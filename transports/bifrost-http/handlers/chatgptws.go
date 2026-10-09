package handlers

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	ws "github.com/fasthttp/websocket"
	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/maximhq/bifrost/transports/bifrost-http/integrations"
	bfws "github.com/maximhq/bifrost/transports/bifrost-http/websocket"
	"github.com/valyala/fasthttp"
)

// Only the swap proxy supplies the company bearer. This endpoint never selects
// another provider key or falls back to HTTP after a partially delivered turn.
type chatGPTWSHooks interface {
	RunPreRequestHooks(*schemas.BifrostContext, *schemas.BifrostRequest)
	RunStreamPreHooks(*schemas.BifrostContext, *schemas.BifrostRequest) (*bifrost.WSStreamHooks, *schemas.BifrostError)
}

func chatGPTWSTarget(rawQuery string) (string, error) {
	base := strings.TrimRight(strings.TrimSpace(os.Getenv(integrations.ChatGPTUpstreamEnv)), "/")
	if base == "" {
		base = "https://chatgpt.com"
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("invalid ChatGPT upstream")
	}
	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		ip := net.ParseIP(u.Hostname())
		if u.Hostname() != "localhost" && (ip == nil || !ip.IsLoopback()) {
			return "", fmt.Errorf("plaintext ChatGPT upstream must be loopback")
		}
		u.Scheme = "ws"
	default:
		return "", fmt.Errorf("invalid ChatGPT upstream scheme")
	}
	u.Path += "/backend-api/codex/responses"
	u.RawQuery = rawQuery
	return u.String(), nil
}

func chatGPTWSHeaders(auth *authHeaders) http.Header {
	headers := http.Header{}
	for name, values := range auth.headers {
		if strings.HasPrefix(name, "x-bf-") || strings.HasPrefix(name, "sec-websocket-") {
			continue
		}
		switch name {
		case "host", "connection", "upgrade", "content-length", "transfer-encoding", "proxy-authorization", "x-request-id":
			continue
		}
		for _, value := range values {
			headers.Add(name, value)
		}
	}
	return headers
}

func (h *WSResponsesHandler) handleChatGPTUpgrade(ctx *fasthttp.RequestCtx) {
	h.handleChatGPTUpgradeWithHooks(ctx, h.client)
}

func (h *WSResponsesHandler) handleChatGPTUpgradeWithHooks(ctx *fasthttp.RequestCtx, runner chatGPTWSHooks) {
	if !ws.FastHTTPIsWebSocketUpgrade(ctx) {
		ctx.Error("websocket upgrade required", fasthttp.StatusBadRequest)
		return
	}
	auth := captureAuthHeaders(ctx)
	if auth.virtualKey == "" || !strings.HasPrefix(auth.authorization, "Bearer ") || strings.HasPrefix(auth.authorization, "Bearer sk-bf-") {
		ctx.Error("company credential and virtual key required", fasthttp.StatusUnauthorized)
		return
	}
	target, err := chatGPTWSTarget(string(ctx.URI().QueryString()))
	if err != nil {
		ctx.Error("invalid ChatGPT upstream", fasthttp.StatusBadGateway)
		return
	}
	reservation, err := h.sessions.Reserve()
	if err != nil {
		ctx.Error("websocket connection limit reached", fasthttp.StatusTooManyRequests)
		return
	}
	var proxy *schemas.ProxyConfig
	if h.config != nil {
		if cfg, cfgErr := h.config.GetProviderConfigRaw(schemas.OpenAI); cfgErr == nil && cfg != nil {
			proxy = cfg.ProxyConfig
		}
	}
	upstream, response, err := bfws.DialUpstreamWithResponse(target, chatGPTWSHeaders(auth), schemas.OpenAI, "", proxy)
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
	if err != nil {
		reservation.Release()
		ctx.Error("ChatGPT websocket connection failed", fasthttp.StatusBadGateway)
		return
	}
	for _, cookie := range response.Header.Values("Set-Cookie") {
		ctx.Response.Header.Add("Set-Cookie", cookie)
	}
	// Native Codex offers permessage-deflate. Each hop negotiates its own
	// compression context; never copy the upstream extension response verbatim.
	upgrader := h.upgrader
	upgrader.EnableCompression = true
	err = upgrader.Upgrade(ctx, func(conn *ws.Conn) {
		defer conn.Close()
		defer upstream.Close()
		session, err := reservation.Complete(conn)
		if err != nil {
			writeWSError(conn, 429, "websocket_connection_limit_reached", err.Error())
			return
		}
		defer h.sessions.Remove(conn)
		session.SetUpstream(upstream)
		h.chatGPTEventLoop(session, auth, target, runner)
	})
	if err != nil {
		reservation.Release()
		_ = upstream.Close()
		ctx.Response.Header.Del("Set-Cookie")
	}
}

func (h *WSResponsesHandler) chatGPTEventLoop(session *bfws.Session, auth *authHeaders, target string, runner chatGPTWSHooks) {
	// Keep reading while a turn is running: a client disconnect must interrupt an
	// upstream read immediately. At most one subsequent message is queued.
	messages := make(chan []byte, 1)
	done := make(chan struct{})
	defer close(done)
	defer session.Close()
	go func() {
		defer close(messages)
		for {
			kind, raw, err := session.ClientConn().ReadMessage()
			if err != nil {
				session.Close()
				return
			}
			if kind != ws.TextMessage {
				session.Close()
				return
			}
			select {
			case messages <- raw:
			case <-done:
				return
			}
		}
	}()
	for raw := range messages {
		if !h.chatGPTTurn(session, auth, target, raw, runner) {
			return
		}
	}
}

func (h *WSResponsesHandler) chatGPTTurn(session *bfws.Session, auth *authHeaders, target string, raw []byte, runner chatGPTWSHooks) bool {
	mapped, mappingErr := rewriteCodexWS(raw, auth)
	if mappingErr != nil {
		writeWSError(session, 400, "invalid_request_error", "invalid native identity context")
		return false
	}
	raw = mapped
	var event schemas.WebSocketResponsesEvent
	if sonic.Unmarshal(raw, &event) != nil || event.Type != schemas.WSEventResponseCreate || event.Model == "" || strings.Contains(event.Model, "/") {
		writeWSError(session, 400, "invalid_request_error", "native response.create and model required")
		return true
	}
	req, err := h.convertEventToRequestWithEmptyInput(&event, true)
	if err != nil {
		writeWSError(session, 400, "invalid_request_error", "invalid response.create")
		return true
	}
	ctx, cancel := createBifrostContextFromAuth(h.handlerStore, auth)
	defer cancel()
	ctx.SetValue(schemas.BifrostContextKeySkipKeySelection, true)
	ctx.SetValue(schemas.BifrostContextKeyParentRequestID, session.ID())
	request := &schemas.BifrostRequest{RequestType: schemas.WebSocketResponsesRequest, ResponsesRequest: req}
	runner.RunPreRequestHooks(ctx, request)
	// Company account and payload are fixed by the swap proxy. A routing plugin
	// must not silently send the unchanged native payload to a different model.
	if req.Provider != schemas.OpenAI || req.Model != event.Model {
		writeWSError(session, 400, "invalid_request_error", "native ChatGPT routing cannot change provider or model")
		return true
	}
	hooks, preErr := runner.RunStreamPreHooks(ctx, request)
	if preErr != nil {
		writeWSBifrostError(session, preErr)
		return true
	}
	defer hooks.Cleanup()
	if hooks.ShortCircuitResponse != nil {
		writeWSShortCircuitResponse(session, hooks.ShortCircuitResponse)
		return true
	}
	finished := false
	tracer, _ := ctx.Value(schemas.BifrostContextKeyTracer).(schemas.Tracer)
	var span schemas.SpanHandle
	if tracer != nil {
		spanID, handle := tracer.StartSpanID(ctx, schemas.OTelOperationName(schemas.WebSocketResponsesRequest)+" "+req.Model, schemas.SpanKindLLMCall)
		span = handle
		tracer.PopulateLLMRequestAttributes(span, request)
		if spanID != "" {
			ctx.SetValue(schemas.BifrostContextKeySpanID, spanID)
		}
	}
	defer func() {
		if !finished {
			failure := newBifrostError(502, "upstream_connection_error", "ChatGPT websocket turn interrupted")
			if tracer != nil {
				tracer.PopulateLLMResponseAttributes(ctx, span, nil, failure)
				tracer.EndSpan(span, schemas.SpanStatusError, "request interrupted")
			}
			ctx.SetValue(schemas.BifrostContextKeyStreamEndIndicator, true)
			hooks.PostHookRunner(ctx, nil, failure)
		}
	}()
	provider, model, _ := request.GetRequestFields()
	if provider != schemas.OpenAI || model != event.Model {
		writeWSError(session, 400, "invalid_request_error", "native ChatGPT routing cannot change provider or model")
		return false
	}
	upstream := session.Upstream()
	if upstream == nil {
		var proxy *schemas.ProxyConfig
		if h.config != nil {
			if cfg, cfgErr := h.config.GetProviderConfigRaw(schemas.OpenAI); cfgErr == nil && cfg != nil {
				proxy = cfg.ProxyConfig
			}
		}
		upstream, err = bfws.DialUpstream(target, chatGPTWSHeaders(auth), schemas.OpenAI, "", proxy)
		if err != nil {
			writeWSError(session, 502, "upstream_connection_error", "ChatGPT websocket connection failed")
			return false
		}
		session.SetUpstream(upstream)
	}
	if err := upstream.SetWriteDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return false
	}
	if err := upstream.WriteMessage(ws.TextMessage, raw); err != nil {
		return false
	}
	for {
		if err := upstream.SetReadDeadline(time.Now().Add(resolveWSStreamIdleTimeout(h.config, schemas.OpenAI))); err != nil {
			return false
		}
		kind, data, err := upstream.ReadMessage()
		if err != nil {
			writeWSError(session, 502, "upstream_connection_error", "ChatGPT websocket stream interrupted")
			return false
		}
		chunk := parseUpstreamWSEvent(data, schemas.OpenAI, req.Model)
		terminal := chunk != nil && isTerminalStreamType(chunk.Type)
		if terminal {
			ctx.SetValue(schemas.BifrostContextKeyStreamEndIndicator, true)
			finished = true
			if tracer != nil {
				failure := buildWSTurnSpanError(chunk)
				tracer.PopulateLLMResponseAttributes(ctx, span, buildWSTurnSpanResponse(chunk, req), failure)
				if failure != nil {
					tracer.EndSpan(span, schemas.SpanStatusError, "upstream response failed")
				} else {
					tracer.EndSpan(span, schemas.SpanStatusOk, "")
				}
			}
		}
		if chunk != nil {
			response := &schemas.BifrostResponse{ResponsesStreamResponse: chunk}
			if tracer, ok := ctx.Value(schemas.BifrostContextKeyTracer).(schemas.Tracer); ok && tracer != nil {
				if traceID, _ := ctx.Value(schemas.BifrostContextKeyTraceID).(string); traceID != "" {
					tracer.AddStreamingChunk(traceID, response)
				}
			}
			_, postErr := hooks.PostHookRunner(ctx, response, buildWSTurnSpanError(chunk))
			if postErr != nil && !terminal {
				writeWSBifrostError(session, postErr)
				return false
			}
		}
		if session.WriteMessage(kind, data) != nil {
			return false
		}
		if terminal {
			session.MarkResponsesTurnCompleted()
			h.trackResponseID(session, data)
			_ = upstream.SetReadDeadline(time.Time{})
			return true
		}
	}
}
