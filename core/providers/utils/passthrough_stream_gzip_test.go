package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/valyala/fasthttp"
)

const passthroughGzipSSE = "event: message_start\ndata: {\"type\":\"message_start\"}\n\n" +
	"event: message_delta\ndata: {\"type\":\"message_delta\",\"usage\":{\"output_tokens\":7}}\n\n"

// runPassthroughStream feeds body (labelled with contentEncoding) through StreamPassthrough
// and returns the forwarded bytes plus how many events Observe saw.
func runPassthroughStream(t *testing.T, body []byte, contentEncoding string) ([]byte, int) {
	t.Helper()
	resp := fasthttp.AcquireResponse()
	if contentEncoding != "" {
		resp.Header.Set("Content-Encoding", contentEncoding)
	}
	resp.SetBodyStream(bytes.NewReader(body), -1)

	ctx := schemas.NewBifrostContext(context.Background(), time.Time{})
	passthrough := func(_ *schemas.BifrostContext, r *schemas.BifrostResponse, e *schemas.BifrostError) (*schemas.BifrostResponse, *schemas.BifrostError) {
		return r, e
	}
	observed := 0
	ch := StreamPassthrough(ctx, passthrough, nil, resp, resp.BodyStream(), PassthroughStreamParams{
		StatusCode: 200,
		StartTime:  time.Now(),
		Observe: func(event []byte) *schemas.BifrostPassthroughUsage {
			observed++
			return nil
		},
	})

	var forwarded []byte
	timeout := time.After(5 * time.Second)
	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				return forwarded, observed
			}
			if chunk != nil && chunk.BifrostPassthroughResponse != nil {
				forwarded = append(forwarded, chunk.BifrostPassthroughResponse.Body...)
			}
		case <-timeout:
			t.Fatal("passthrough stream did not finish")
		}
	}
}

// Upstream gzip must be decoded: Content-Encoding is not forwarded, so encoded bytes would
// reach the client labelled as plain SSE (FWD-56: Claude Code saw a double-gzipped stream
// and retried every turn without streaming).
func TestStreamPassthroughDecodesGzip(t *testing.T) {
	var gz bytes.Buffer
	w := gzip.NewWriter(&gz)
	w.Write([]byte(passthroughGzipSSE))
	w.Close()

	forwarded, observed := runPassthroughStream(t, gz.Bytes(), "gzip")
	if string(forwarded) != passthroughGzipSSE {
		t.Fatalf("forwarded %q, want decoded SSE", forwarded)
	}
	if observed != 2 {
		t.Fatalf("Observe saw %d events, want 2", observed)
	}
}

func TestStreamPassthroughKeepsPlainBody(t *testing.T) {
	forwarded, observed := runPassthroughStream(t, []byte(passthroughGzipSSE), "")
	if string(forwarded) != passthroughGzipSSE {
		t.Fatalf("forwarded %q, want body unchanged", forwarded)
	}
	if observed != 2 {
		t.Fatalf("Observe saw %d events, want 2", observed)
	}
}

// fasthttp reports text/plain for a response without Content-Type; passthrough must not
// forward that invented type (FWD-56: Codex SSE reached the client as text/plain).
func TestExtractPassthroughHeadersKeepsMissingContentType(t *testing.T) {
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	resp.Header.Set("X-Codex-Plan-Type", "team")
	headers := ExtractPassthroughProviderResponseHeaders(resp)
	for k := range headers {
		if strings.EqualFold(k, "content-type") {
			t.Fatalf("invented content-type %q", headers[k])
		}
	}
	if headers["X-Codex-Plan-Type"] != "team" {
		t.Fatalf("headers = %v", headers)
	}

	typed := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(typed)
	typed.Header.SetContentType("text/event-stream")
	if got := ExtractPassthroughProviderResponseHeaders(typed)["Content-Type"]; got != "text/event-stream" {
		t.Fatalf("content-type = %q, want text/event-stream", got)
	}
}
