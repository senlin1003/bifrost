package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fasthttp/router"
	ws "github.com/fasthttp/websocket"
	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/transports/bifrost-http/integrations"
	bfws "github.com/maximhq/bifrost/transports/bifrost-http/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/maximhq/bifrost/framework/configstore"
	configstoreTables "github.com/maximhq/bifrost/framework/configstore/tables"
	"github.com/maximhq/bifrost/framework/grant"
	"github.com/maximhq/bifrost/framework/kvstore"
	"github.com/maximhq/bifrost/framework/logstore"
	"github.com/maximhq/bifrost/transports/bifrost-http/lib"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type codexTestHooks struct {
	pre     int
	post    int
	cleanup int
	deny    int
	t       *testing.T
}

type codexCompressionWireConn struct {
	net.Conn
	captured []byte
}

func (c *codexCompressionWireConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	if n > 0 && len(c.captured) < 65536 {
		keep := n
		if keep > 65536-len(c.captured) {
			keep = 65536 - len(c.captured)
		}
		c.captured = append(c.captured, p[:keep]...)
	}
	return n, err
}

// 验收: CXD-R1~ FWD-67-A2~
func TestChatGPTWSCompressionAndLongMessages(t *testing.T) {
	for _, offer := range []bool{true, false} {
		t.Run(fmt.Sprint(offer), func(t *testing.T) {
			long := strings.Repeat("fake long delta ", 8192)
			delta := `{"type":"response.output_text.delta","delta":"` + long + `"}`
			done := `{"type":"response.completed","response":{"id":"resp-compression","usage":{"input_tokens":2,"output_tokens":3,"total_tokens":5}}}`
			handshake := make(chan string, 1)
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handshake <- r.Header.Get("Sec-WebSocket-Extensions")
				conn, err := (&ws.Upgrader{EnableCompression: true}).Upgrade(w, r, http.Header{"Set-Cookie": {"fake=one; Path=/", "fake=two; Path=/backend-api"}})
				if err != nil {
					t.Error(err)
					return
				}
				defer conn.Close()
				conn.EnableWriteCompression(true)
				for {
					_, raw, err := conn.ReadMessage()
					if err != nil {
						return
					}
					if !strings.Contains(string(raw), `"model":"gpt-6.1-sol"`) {
						t.Error("native payload changed")
					}
					_ = conn.WriteMessage(ws.TextMessage, []byte(delta))
					_ = conn.WriteMessage(ws.TextMessage, []byte(done))
				}
			}))
			defer upstream.Close()
			t.Setenv(integrations.ChatGPTUpstreamEnv, upstream.URL)
			h := &WSResponsesHandler{sessions: bfws.NewSessionManager(1), handlerStore: testWSHandlerStore{}}
			hooks := &codexTestHooks{t: t}
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			server := &fasthttp.Server{Handler: func(ctx *fasthttp.RequestCtx) { h.handleChatGPTUpgradeWithHooks(ctx, hooks) }}
			go server.Serve(listener)
			defer server.Shutdown()
			dialer := ws.Dialer{EnableCompression: offer, HandshakeTimeout: 3 * time.Second}
			var wire *codexCompressionWireConn
			dialer.NetDial = func(network, address string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, address, 3*time.Second)
				if err != nil {
					return nil, err
				}
				wire = &codexCompressionWireConn{Conn: conn}
				return wire, nil
			}
			conn, resp, err := dialer.Dial("ws://"+listener.Addr().String()+"/chatgpt_passthrough/backend-api/codex/responses", http.Header{"Authorization": {"Bearer fake-company"}, "X-Bf-Vk": {"sk-bf-test"}})
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			if ext := resp.Header.Get("Sec-WebSocket-Extensions"); strings.Contains(ext, "permessage-deflate") != offer {
				t.Errorf("downstream compression offer=%v, negotiated=%q", offer, ext)
			}
			if ext := <-handshake; !strings.Contains(ext, "permessage-deflate") {
				t.Errorf("upstream compression absent: %q", ext)
			}
			if len(resp.Header.Values("Set-Cookie")) != 2 {
				t.Fatal("cookies lost")
			}
			conn.EnableWriteCompression(offer)
			for i := 0; i < 2; i++ {
				_ = conn.WriteMessage(ws.TextMessage, []byte(`{"type":"response.create","model":"gpt-6.1-sol","store":false,"reasoning":{"effort":"high"},"input":[]}`))
				_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
				for _, expected := range []string{delta, done} {
					_, raw, err := conn.ReadMessage()
					if err != nil || string(raw) != expected {
						t.Fatal("compressed native message changed or interrupted", err)
					}
				}
			}
			first := bytes.Index(wire.captured, []byte("\r\n\r\n")) + 4
			if first < 4 || first >= len(wire.captured) {
				t.Fatal("wire frame missing")
			}
			if compressed := wire.captured[first]&0x40 != 0; compressed != offer {
				t.Errorf("actual downstream frame RSV1=%v offer=%v", compressed, offer)
			}
			t.Logf("offer=%v downstream=%s RSV1=%v recorded_bytes=%d plain_delta_bytes=%d", offer, resp.Header.Get("Sec-WebSocket-Extensions"), wire.captured[first]&0x40 != 0, len(wire.captured), len(long))
			if offer && len(wire.captured) >= len(long) {
				t.Fatal("long messages were not compressed on wire")
			}
		})
	}
}

// 验收: FWD-67-A2~ Real handshake, without company credentials or a paid upstream.
func TestChatGPTWSHandshakeCookies(t *testing.T) {
	want := []string{"affinity=one; Path=/; HttpOnly", "affinity=two; Path=/backend-api; Expires=Wed, 21 Oct 2037 07:28:00 GMT; Secure"}
	closed := make(chan struct{})
	var dials, messages atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dials.Add(1)
		conn, err := (&ws.Upgrader{}).Upgrade(w, r, http.Header{"Set-Cookie": want})
		if err != nil {
			t.Error(err)
			return
		}
		defer close(closed)
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
			messages.Add(1)
			_ = conn.WriteMessage(ws.TextMessage, []byte(`{"type":"response.completed","response":{"id":"resp-test","usage":{"input_tokens":2,"output_tokens":3,"total_tokens":5}}}`))
		}
	}))
	defer upstream.Close()
	t.Setenv(integrations.ChatGPTUpstreamEnv, upstream.URL)
	h := &WSResponsesHandler{sessions: bfws.NewSessionManager(1), handlerStore: testWSHandlerStore{}}
	hooks := &codexTestHooks{t: t, deny: 3}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &fasthttp.Server{Handler: func(ctx *fasthttp.RequestCtx) { h.handleChatGPTUpgradeWithHooks(ctx, hooks) }}
	go server.Serve(listener)
	defer server.Shutdown()
	header := http.Header{"Authorization": {"Bearer company.test.token"}, "X-Bf-Vk": {"sk-bf-test"}}
	conn, resp, err := ws.DefaultDialer.Dial("ws://"+listener.Addr().String()+"/chatgpt_passthrough/backend-api/codex/responses", header)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if got := resp.Header.Values("Set-Cookie"); !reflect.DeepEqual(got, want) {
		t.Errorf("cookies = %q, want %q", got, want)
	}
	for i := 0; i < 3; i++ {
		_ = conn.WriteMessage(ws.TextMessage, []byte(`{"type":"response.create","model":"gpt-5","store":false,"input":[]}`))
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		if i < 2 && !strings.Contains(string(raw), "response.completed") {
			t.Fatalf("unexpected response: %s", raw)
		}
		if i == 2 && !strings.Contains(string(raw), "quota exhausted") {
			t.Fatal("quota bypass")
		}
	}
	_ = conn.Close()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Error("upstream not closed after client disconnect")
	}
	if dials.Load() != 1 || messages.Load() != 2 {
		t.Fatalf("dials=%d messages=%d", dials.Load(), messages.Load())
	}
}

// 验收: FWD-67-A2~ Failed upstream/upgrade releases the reserved slot and cookies.
func TestChatGPTWSHandshakeFailureCleanup(t *testing.T) {
	for _, rejected := range []bool{false, true} {
		t.Run(fmt.Sprint(rejected), func(t *testing.T) {
			closed := make(chan struct{})
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !rejected {
					w.Header().Set("Set-Cookie", "failure=private")
					w.WriteHeader(401)
					return
				}
				conn, err := (&ws.Upgrader{}).Upgrade(w, r, http.Header{"Set-Cookie": {"success=private"}})
				if err != nil {
					t.Error(err)
					return
				}
				defer conn.Close()
				defer close(closed)
				_, _, _ = conn.ReadMessage()
			}))
			defer upstream.Close()
			t.Setenv(integrations.ChatGPTUpstreamEnv, upstream.URL)
			h := &WSResponsesHandler{sessions: bfws.NewSessionManager(1)}
			if rejected {
				h.upgrader.CheckOrigin = func(*fasthttp.RequestCtx) bool { return false }
			}
			var ctx fasthttp.RequestCtx
			ctx.Request.SetRequestURI("/chatgpt_passthrough/backend-api/codex/responses")
			ctx.Request.Header.SetMethod("GET")
			for key, value := range map[string]string{"Connection": "Upgrade", "Upgrade": "websocket", "Sec-WebSocket-Version": "13", "Sec-WebSocket-Key": "dGhlIHNhbXBsZSBub25jZQ==", "Authorization": "Bearer company.test.token", "X-Bf-Vk": "sk-bf-test"} {
				ctx.Request.Header.Set(key, value)
			}
			h.handleChatGPTUpgrade(&ctx)
			want := 502
			if rejected {
				want = 403
			}
			if ctx.Response.StatusCode() != want || len(ctx.Response.Header.Peek("Set-Cookie")) != 0 {
				t.Fatalf("status=%d cookie leaked=%v", ctx.Response.StatusCode(), len(ctx.Response.Header.Peek("Set-Cookie")) != 0)
			}
			r, err := h.sessions.Reserve()
			if err != nil {
				t.Fatal("reservation leaked")
			}
			r.Release()
			if rejected {
				select {
				case <-closed:
				case <-time.After(2 * time.Second):
					t.Fatal("upstream leaked")
				}
			}
		})
	}
}

// 验收: FWD-67-A2~ No invented cookies; concurrent handshakes obey admission.
func TestChatGPTWSHandshakeAdmission(t *testing.T) {
	var dials atomic.Int32
	closed := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dials.Add(1)
		conn, err := (&ws.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		defer close(closed)
		_, _, _ = conn.ReadMessage()
	}))
	defer upstream.Close()
	t.Setenv(integrations.ChatGPTUpstreamEnv, upstream.URL)
	h := &WSResponsesHandler{sessions: bfws.NewSessionManager(1)}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := &fasthttp.Server{Handler: h.handleChatGPTUpgrade}
	go server.Serve(listener)
	defer server.Shutdown()
	endpoint := "ws://" + listener.Addr().String() + "/chatgpt_passthrough/backend-api/codex/responses"
	header := http.Header{"Authorization": {"Bearer company.test.token"}, "X-Bf-Vk": {"sk-bf-test"}}
	conn, resp, err := ws.DefaultDialer.Dial(endpoint, header)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if len(resp.Header.Values("Set-Cookie")) != 0 {
		t.Fatal("invented cookie")
	}
	second, resp, err := ws.DefaultDialer.Dial(endpoint, header)
	if second != nil {
		second.Close()
	}
	if err == nil || resp == nil || resp.StatusCode != 429 {
		t.Fatal("connection limit bypass")
	}
	if dials.Load() != 1 {
		t.Fatal("rejected connection dialed upstream")
	}
	conn.Close()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("idle upstream leaked")
	}
}

func (r *codexTestHooks) RunPreRequestHooks(ctx *schemas.BifrostContext, req *schemas.BifrostRequest) {
}
func (r *codexTestHooks) RunStreamPreHooks(ctx *schemas.BifrostContext, req *schemas.BifrostRequest) (*bifrost.WSStreamHooks, *schemas.BifrostError) {
	r.pre++
	if ctx.Value(schemas.BifrostContextKeyVirtualKey) != "sk-bf-test" || ctx.Value(schemas.BifrostContextKeySkipKeySelection) != true {
		r.t.Error("missing per-turn governance identity")
	}
	if req.ResponsesRequest.Params.Store == nil || *req.ResponsesRequest.Params.Store {
		r.t.Error("store=false overwritten")
	}
	if r.pre == r.deny {
		return nil, newBifrostError(429, "rate_limit_exceeded", "test quota exhausted")
	}
	return &bifrost.WSStreamHooks{
		Cleanup: func() { r.cleanup++ },
		PostHookRunner: func(ctx *schemas.BifrostContext, resp *schemas.BifrostResponse, err *schemas.BifrostError) (*schemas.BifrostResponse, *schemas.BifrostError) {
			if ctx.Value(schemas.BifrostContextKeyStreamEndIndicator) == true {
				r.post++
				if resp != nil && (resp.ResponsesStreamResponse.Response == nil || resp.ResponsesStreamResponse.Response.Usage == nil) {
					r.t.Error("terminal usage missing from post hook")
				}
			}
			return resp, err
		},
	}, nil
}

// 验收: FWD-67-A1~ FWD-67-A2~
func TestChatGPTWSMultiTurnRawAndQuota(t *testing.T) {
	for _, denied := range []bool{false, true} {
		t.Run(fmt.Sprint(denied), func(t *testing.T) {
			var dials, requests atomic.Int32
			upgrader := ws.Upgrader{}
			payload := `{ "type":"response.create", "model":"gpt-5", "store":false, "input":[], "unknown":{"preserve":1}, "prompt_cache_key":"native-session", "client_metadata":{"session_id":"native-session","turn_id":"independent-turn"} }`
			expectedPayload := strings.ReplaceAll(payload, `"native-session"`, `"company-session"`)
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				dials.Add(1)
				if req.URL.RawQuery != "test=1" || req.Header.Get("Authorization") != "Bearer company.test.token" || req.Header.Get("Cookie") != "session=test" || req.Header.Get("Session-Id") != "company-session" || req.Header.Get("X-Bf-Vk") != "" || req.Header.Get(codexIdentityHeader) != "" {
					t.Error("upstream handshake fidelity or internal key leak")
				}
				conn, err := upgrader.Upgrade(w, req, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer conn.Close()
				for {
					_, raw, err := conn.ReadMessage()
					if err != nil {
						return
					}
					requests.Add(1)
					if string(raw) != expectedPayload {
						t.Error("native message changed")
					}
					_ = conn.WriteMessage(ws.TextMessage, []byte(`{"type":"response.completed","sequence_number":1,"response":{"id":"resp-test","usage":{"input_tokens":2,"output_tokens":3,"total_tokens":5}}}`))
				}
			}))
			defer upstream.Close()
			t.Setenv(integrations.ChatGPTUpstreamEnv, upstream.URL)
			target, err := chatGPTWSTarget("test=1")
			if err != nil {
				t.Fatal(err)
			}
			hooks := &codexTestHooks{t: t}
			if denied {
				hooks.deny = 2
			}
			done := make(chan struct{})
			handler := &WSResponsesHandler{handlerStore: testWSHandlerStore{}}
			auth := &authHeaders{virtualKey: "sk-bf-test", authorization: "Bearer company.test.token", headers: map[string][]string{
				"authorization": {"Bearer company.test.token"}, "x-bf-vk": {"sk-bf-test"}, "cookie": {"session=test"}, "session-id": {"company-session"},
				codexIdentityHeader: {`{"original":"native-session","mapped":"company-session","installation":"company-device"}`},
			}}
			gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				conn, err := upgrader.Upgrade(w, req, nil)
				if err != nil {
					t.Error(err)
					return
				}
				defer close(done)
				handler.chatGPTEventLoop(bfws.NewSession(conn), auth, target, hooks)
			}))
			defer gateway.Close()
			conn, _, err := ws.DefaultDialer.Dial("ws"+strings.TrimPrefix(gateway.URL, "http"), nil)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			for i := 0; i < 2; i++ {
				_ = conn.WriteMessage(ws.TextMessage, []byte(payload))
				_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				_, raw, err := conn.ReadMessage()
				if err != nil {
					t.Fatal(err)
				}
				if denied && i == 1 {
					if !strings.Contains(string(raw), "quota exhausted") {
						t.Fatalf("quota bypass: %s", raw)
					}
				} else if !strings.Contains(string(raw), "response.completed") {
					t.Fatalf("unexpected event: %s", raw)
				}
			}
			_ = conn.Close()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("disconnect did not release session")
			}
			want := 2
			if denied {
				want = 1
			}
			if dials.Load() != 1 || int(requests.Load()) != want || hooks.pre != 2 || hooks.post != want || hooks.cleanup != want {
				t.Fatalf("dials=%d requests=%d pre=%d post=%d cleanup=%d", dials.Load(), requests.Load(), hooks.pre, hooks.post, hooks.cleanup)
			}
		})
	}
}

func TestChatGPTWSTargetRejectsPlaintextRemote(t *testing.T) {
	for _, base := range []string{"http://example.com", "https://user:password@example.com", "https://chatgpt.com?redirect=1"} {
		t.Setenv(integrations.ChatGPTUpstreamEnv, base)
		if _, err := chatGPTWSTarget(""); err == nil {
			t.Fatalf("unsafe upstream accepted: %s", base)
		}
	}
}

// 验收: FWD-67-A1~
func TestChatGPTWSRouteRequiresUpgradeAndCredentials(t *testing.T) {
	h := &WSResponsesHandler{}
	r := router.New()
	h.RegisterRoutes(r)
	for _, upgrade := range []bool{false, true} {
		var ctx fasthttp.RequestCtx
		ctx.Request.SetRequestURI("/chatgpt_passthrough/backend-api/codex/responses")
		ctx.Request.Header.SetMethod("GET")
		want := 400
		if upgrade {
			ctx.Request.Header.Set("Connection", "Upgrade")
			ctx.Request.Header.Set("Upgrade", "websocket")
			want = 401
		}
		r.Handler(&ctx)
		if ctx.Response.StatusCode() != want {
			t.Fatalf("upgrade=%v status=%d want=%d", upgrade, ctx.Response.StatusCode(), want)
		}
	}
}

// 验收: FWD-67-A1~
func TestChatGPTWSDisconnectDuringTurnFinalizesHooks(t *testing.T) {
	upgrader := ws.Upgrader{}
	started := make(chan struct{})
	upstreamClosed := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
		defer close(upstreamClosed)
		if _, _, err := conn.ReadMessage(); err != nil {
			t.Error(err)
			return
		}
		close(started)
		_, _, _ = conn.ReadMessage()
	}))
	defer upstream.Close()
	hooks := &codexTestHooks{t: t}
	h := &WSResponsesHandler{handlerStore: testWSHandlerStore{}}
	done := make(chan struct{})
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			t.Error(err)
			return
		}
		defer close(done)
		h.chatGPTEventLoop(bfws.NewSession(conn), &authHeaders{virtualKey: "sk-bf-test", headers: map[string][]string{}}, "ws"+strings.TrimPrefix(upstream.URL, "http"), hooks)
	}))
	defer gateway.Close()
	conn, _, err := ws.DefaultDialer.Dial("ws"+strings.TrimPrefix(gateway.URL, "http"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_ = conn.WriteMessage(ws.TextMessage, []byte(`{"type":"response.create","model":"gpt-5","store":false,"input":[]}`))
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("turn did not reach upstream")
	}
	_ = conn.Close()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("turn not interrupted by client disconnect")
	}
	select {
	case <-upstreamClosed:
	case <-time.After(5 * time.Second):
		t.Fatal("upstream connection leaked")
	}
	if hooks.pre != 1 || hooks.post != 1 || hooks.cleanup != 1 {
		t.Fatalf("unbalanced hooks: %+v", hooks)
	}
}

type testWSHandlerStore struct {
	matcher *lib.HeaderMatcher
}

func (s testWSHandlerStore) GetHeaderMatcher() *lib.HeaderMatcher {
	return s.matcher
}

func (s testWSHandlerStore) GetStreamChunkInterceptor() lib.StreamChunkInterceptor {
	return nil
}

func (s testWSHandlerStore) GetAsyncJobExecutor() *logstore.AsyncJobExecutor {
	return nil
}

func (s testWSHandlerStore) GetAsyncJobResultTTL() int {
	return 0
}

func (s testWSHandlerStore) GetKVStore() *kvstore.Store {
	return nil
}

func (s testWSHandlerStore) GetMCPHeaderCombinedAllowlist() schemas.WhiteList {
	return nil
}

func (s testWSHandlerStore) ShouldAllowPerRequestStorageOverride() bool { return false }
func (s testWSHandlerStore) ShouldAllowPerRequestRawOverride() bool     { return false }
func (s testWSHandlerStore) ShouldAllowDirectKeys() bool                { return false }
func (s testWSHandlerStore) GetMCPExternalServerURL() string            { return "" }
func (s testWSHandlerStore) GetMCPExternalClientURL() string            { return "" }

type timeoutNetError struct{}

func (timeoutNetError) Error() string   { return "i/o timeout" }
func (timeoutNetError) Timeout() bool   { return true }
func (timeoutNetError) Temporary() bool { return false }

func TestResolveWSStreamIdleTimeoutUsesProviderOverride(t *testing.T) {
	cfg := &lib.Config{
		Providers: map[schemas.ModelProvider]configstore.ProviderConfig{
			schemas.OpenAI: {
				NetworkConfig: &schemas.NetworkConfig{StreamIdleTimeoutInSeconds: 7},
			},
		},
	}

	timeout := resolveWSStreamIdleTimeout(cfg, schemas.OpenAI)
	assert.Equal(t, 7*time.Second, timeout)
}

func TestResolveWSStreamIdleTimeoutFallsBackToDefault(t *testing.T) {
	timeout := resolveWSStreamIdleTimeout(&lib.Config{}, schemas.OpenAI)
	assert.Equal(t, time.Duration(schemas.DefaultStreamIdleTimeoutInSeconds)*time.Second, timeout)
}

func TestIsWSReadTimeout(t *testing.T) {
	assert.True(t, isWSReadTimeout(timeoutNetError{}))
	assert.False(t, isWSReadTimeout(net.UnknownNetworkError("unknown")))
	assert.False(t, isWSReadTimeout(errors.New("boom")))
	assert.False(t, isWSReadTimeout(nil))
}

func TestNewBifrostError(t *testing.T) {
	bifrostErr := newBifrostError(504, "upstream_timeout", "upstream websocket stream timed out")
	if bifrostErr == nil {
		t.Fatal("expected bifrost error, got nil")
	}
	if bifrostErr.StatusCode == nil || *bifrostErr.StatusCode != 504 {
		t.Fatalf("status code = %#v, want 504", bifrostErr.StatusCode)
	}
	if bifrostErr.Error == nil {
		t.Fatal("expected error field, got nil")
	}
	if bifrostErr.Error.Type == nil || *bifrostErr.Error.Type != "upstream_timeout" {
		t.Fatalf("error type = %#v, want upstream_timeout", bifrostErr.Error.Type)
	}
	if bifrostErr.Error.Message != "upstream websocket stream timed out" {
		t.Fatalf("error message = %q, want upstream websocket stream timed out", bifrostErr.Error.Message)
	}
}

func TestCreateBifrostContextFromAuth_BaggageSessionIDSetsGrouping(t *testing.T) {
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{}, &authHeaders{
		baggage: "foo=bar, session-id=rt-ws-123, baz=qux",
	})
	defer cancel()

	if got, _ := ctx.Value(schemas.BifrostContextKeyParentRequestID).(string); got != "rt-ws-123" {
		t.Fatalf("parent request id = %q, want %q", got, "rt-ws-123")
	}
}

// A connection settles its identity at upgrade the way an HTTP request does at conversion, so
// governance sees a settled request whether or not a key was presented.
func TestCreateBifrostContextFromAuth_SettlesIdentity(t *testing.T) {
	cases := map[string]struct {
		auth     *authHeaders
		wantKind string
		wantKey  string
	}{
		"no credential":       {auth: &authHeaders{}},
		"virtual key header":  {auth: &authHeaders{virtualKey: "sk-bf-ws"}, wantKind: string(grant.CredentialVirtualKey), wantKey: "sk-bf-ws"},
		"bearer virtual key":  {auth: &authHeaders{authorization: "Bearer sk-bf-bearer"}, wantKind: string(grant.CredentialVirtualKey), wantKey: "sk-bf-bearer"},
		"provider key bearer": {auth: &authHeaders{authorization: "Bearer sk-openai"}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{}, tc.auth)
			defer cancel()

			g := ctx.Grant()
			if g == nil {
				t.Fatal("expected a grant on the connection context")
			}
			identity := g.Identity()
			if identity == nil {
				t.Fatal("expected a settled identity on the grant")
			}
			if got := identity.Credential(); got.Kind != tc.wantKind || got.Value != tc.wantKey {
				t.Fatalf("credential = %+v, want kind %q value %q", got, tc.wantKind, tc.wantKey)
			}
		})
	}
}

func TestCreateBifrostContextFromAuth_EmptyBaggageSessionIDIgnored(t *testing.T) {
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{}, &authHeaders{
		baggage: "session-id=   ",
	})
	defer cancel()

	if got := ctx.Value(schemas.BifrostContextKeyParentRequestID); got != nil {
		t.Fatalf("parent request id should be unset, got %#v", got)
	}
}

func TestCreateBifrostContextFromAuth_ForwardsPrefixedHeaders(t *testing.T) {
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{}, &authHeaders{
		headers: map[string][]string{
			"x-bf-eh-originator": {"my-test-client"},
			"x-bf-eh-x-trace-id": {"abc-123"},
			"x-bf-eh-cookie":     {"blocked"},
		},
	})
	defer cancel()

	extraHeaders, ok := ctx.Value(schemas.BifrostContextKeyExtraHeaders).(map[string][]string)
	if !ok {
		t.Fatal("expected websocket extra headers in context")
	}
	assert.Equal(t, []string{"my-test-client"}, extraHeaders["originator"])
	assert.Equal(t, []string{"abc-123"}, extraHeaders["x-trace-id"])
	assert.NotContains(t, extraHeaders, "cookie")
}

func TestCreateBifrostContextFromAuth_AppliesHeaderFilterAndDirectAllowlist(t *testing.T) {
	matcher := lib.NewHeaderMatcher(&configstoreTables.GlobalHeaderFilterConfig{
		Allowlist: []string{"originator", "anthropic-*"},
		Denylist:  []string{"anthropic-secret"},
	})
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{matcher: matcher}, &authHeaders{
		headers: map[string][]string{
			"x-bf-eh-originator":       {"allowed-prefix"},
			"x-bf-eh-x-trace-id":       {"blocked-by-allowlist"},
			"anthropic-beta":           {"allowed-direct"},
			"anthropic-secret":         {"blocked-by-denylist"},
			"x-bf-eh-anthropic-secret": {"blocked-prefix-denylist"},
		},
	})
	defer cancel()

	extraHeaders, ok := ctx.Value(schemas.BifrostContextKeyExtraHeaders).(map[string][]string)
	if !ok {
		t.Fatal("expected websocket extra headers in context")
	}
	assert.Equal(t, []string{"allowed-prefix"}, extraHeaders["originator"])
	assert.Equal(t, []string{"allowed-direct"}, extraHeaders["anthropic-beta"])
	assert.NotContains(t, extraHeaders, "x-trace-id")
	assert.NotContains(t, extraHeaders, "anthropic-secret")
}

func TestCaptureAuthHeaders_PreservesDuplicateHeaderValues(t *testing.T) {
	var req fasthttp.Request
	req.Header.Set("Host", "example.test")
	req.Header.Add("x-bf-eh-x-trace-id", "trace-a")
	req.Header.Add("x-bf-eh-x-trace-id", "trace-b")

	ctx := &fasthttp.RequestCtx{}
	ctx.Init(&req, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}, nil)

	auth := captureAuthHeaders(ctx)
	assert.Equal(t, []string{"trace-a", "trace-b"}, auth.headers["x-bf-eh-x-trace-id"])
}

func TestCreateBifrostContextFromAuth_PreservesMultipleForwardedHeaderValues(t *testing.T) {
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{}, &authHeaders{
		headers: map[string][]string{
			"x-bf-eh-x-trace-id": {"trace-a", "trace-b"},
		},
	})
	defer cancel()

	extraHeaders, ok := ctx.Value(schemas.BifrostContextKeyExtraHeaders).(map[string][]string)
	if !ok {
		t.Fatal("expected websocket extra headers in context")
	}
	assert.Equal(t, []string{"trace-a", "trace-b"}, extraHeaders["x-trace-id"])
}

func TestCreateBifrostContextFromAuth_BlocksWebSocketHandshakeForwardedHeaders(t *testing.T) {
	matcher := lib.NewHeaderMatcher(&configstoreTables.GlobalHeaderFilterConfig{
		Allowlist: []string{"*"},
	})
	ctx, cancel := createBifrostContextFromAuth(testWSHandlerStore{matcher: matcher}, &authHeaders{
		headers: map[string][]string{
			"x-bf-eh-upgrade":                {"websocket"},
			"x-bf-eh-sec-websocket-protocol": {"realtime"},
			"sec-websocket-extensions":       {"permessage-deflate"},
			"x-bf-eh-originator":             {"safe"},
		},
	})
	defer cancel()

	extraHeaders, ok := ctx.Value(schemas.BifrostContextKeyExtraHeaders).(map[string][]string)
	if !ok {
		t.Fatal("expected websocket extra headers in context")
	}
	assert.Equal(t, []string{"safe"}, extraHeaders["originator"])
	assert.NotContains(t, extraHeaders, "upgrade")
	assert.NotContains(t, extraHeaders, "sec-websocket-protocol")
	assert.NotContains(t, extraHeaders, "sec-websocket-extensions")
}

func TestMergeWebSocketHeaders_ForwardedHeadersOverrideProviderHeadersAndPreserveValues(t *testing.T) {
	ctx := schemas.NewBifrostContext(nil, time.Time{})
	ctx.SetValue(schemas.BifrostContextKeyExtraHeaders, map[string][]string{
		"originator":    {"my-test-client"},
		"authorization": {"Bearer malicious"},
		"x-static":      {"client-value-1", "client-value-2"},
	})

	merged := mergeWebSocketHeaders(ctx, map[string]string{
		"Authorization": "Bearer provider-key",
		"x-static":      "provider-value",
	})

	assert.Equal(t, []string{"my-test-client"}, merged.Values("originator"))
	assert.Equal(t, []string{"Bearer provider-key"}, merged.Values("Authorization"))
	assert.Equal(t, []string{"client-value-1", "client-value-2"}, merged.Values("x-static"))
	assert.NotContains(t, merged.Values("Authorization"), "Bearer malicious")
}

func TestHasWebSocketForwardedHeaders(t *testing.T) {
	ctx := schemas.NewBifrostContext(nil, time.Time{})
	assert.False(t, hasWebSocketForwardedHeaders(ctx))

	ctx.SetValue(schemas.BifrostContextKeyExtraHeaders, map[string][]string{
		"authorization": {"Bearer malicious"},
	})
	assert.False(t, hasWebSocketForwardedHeaders(ctx))

	ctx.SetValue(schemas.BifrostContextKeyExtraHeaders, map[string][]string{
		"x-trace-id": {"abc-123"},
	})
	assert.True(t, hasWebSocketForwardedHeaders(ctx))
}
