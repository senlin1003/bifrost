package utils

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/valyala/fasthttp"
)

// flushWriter is a streaming compressor that can push out everything written so far.
type flushWriter interface {
	io.WriteCloser
	Flush() error
}

var streamCodings = map[string]func(w io.Writer) flushWriter{
	"gzip":    func(w io.Writer) flushWriter { return gzip.NewWriter(w) },
	"deflate": func(w io.Writer) flushWriter { return zlib.NewWriter(w) },
	"br":      func(w io.Writer) flushWriter { return brotli.NewWriter(w) },
	"zstd": func(w io.Writer) flushWriter {
		enc, _ := zstd.NewWriter(w)
		return enc
	},
}

// 验收: FWD-64-A2
// Every coding the streaming Accept-Encoding now offers must come out decoded, with usage
// still observed -- otherwise the client gets compressed bytes labelled as plain SSE.
func TestStreamPassthroughDecodesEveryOfferedCoding(t *testing.T) {
	for name, newWriter := range streamCodings {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			w := newWriter(&buf)
			w.Write([]byte(passthroughGzipSSE))
			w.Close()

			forwarded, observed := runPassthroughStream(t, buf.Bytes(), name)
			if string(forwarded) != passthroughGzipSSE {
				t.Fatalf("forwarded %q, want decoded SSE", forwarded)
			}
			if observed != 2 {
				t.Fatalf("Observe saw %d events, want 2", observed)
			}
		})
	}
}

// 验收: FWD-64-A2
// Decoding must not hold events back until the body ends: the first event, once the
// upstream flushes it, has to reach the client while the rest is still being written.
func TestStreamPassthroughForwardsEachCodingIncrementally(t *testing.T) {
	first := "event: message_start\ndata: {\"type\":\"message_start\"}\n\n"
	rest := "event: message_delta\ndata: {\"type\":\"message_delta\",\"usage\":{\"output_tokens\":7}}\n\n"
	for name, newWriter := range streamCodings {
		t.Run(name, func(t *testing.T) {
			pr, pw := io.Pipe()
			w := newWriter(pw)
			resp := fasthttp.AcquireResponse()
			resp.Header.Set("Content-Encoding", name)
			resp.SetBodyStream(pr, -1)

			// The upstream sends its first event before the decoder starts reading (the gzip
			// reader parses the header up front), then holds the connection open.
			go func() {
				w.Write([]byte(first))
				w.Flush()
			}()

			ctx := schemas.NewBifrostContext(context.Background(), time.Time{})
			passthrough := func(_ *schemas.BifrostContext, r *schemas.BifrostResponse, e *schemas.BifrostError) (*schemas.BifrostResponse, *schemas.BifrostError) {
				return r, e
			}
			ch := StreamPassthrough(ctx, passthrough, nil, resp, resp.BodyStream(), PassthroughStreamParams{
				StatusCode: 200,
				StartTime:  time.Now(),
			})

			var forwarded []byte
			deadline := time.After(5 * time.Second)
			for !strings.Contains(string(forwarded), "message_start") {
				select {
				case chunk, ok := <-ch:
					if !ok {
						t.Fatal("stream ended before the first event")
					}
					if chunk != nil && chunk.BifrostPassthroughResponse != nil {
						forwarded = append(forwarded, chunk.BifrostPassthroughResponse.Body...)
					}
				case <-deadline:
					t.Fatalf("first event not forwarded while the stream is still open; got %q", forwarded)
				}
			}

			go func() {
				w.Write([]byte(rest))
				w.Close()
				pw.Close()
			}()
			for chunk := range ch {
				if chunk != nil && chunk.BifrostPassthroughResponse != nil {
					forwarded = append(forwarded, chunk.BifrostPassthroughResponse.Body...)
				}
			}
			if string(forwarded) != first+rest {
				t.Fatalf("forwarded %q, want %q", forwarded, first+rest)
			}
		})
	}
}
