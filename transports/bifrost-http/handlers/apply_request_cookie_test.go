package handlers

import (
	"strings"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/valyala/fasthttp"
)

// 验收: FWD-65-A2
// Writing the caller's headers back after each plugin phase must not multiply the Cookie
// header: the provider has to receive the cookies exactly once, as from a direct client.
func TestApplyHTTPRequestToCtxKeepsCookieOnce(t *testing.T) {
	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/backend-api/codex/responses")
	ctx.Request.Header.Set("Cookie", "__cflb=a; __cf_bm=b")

	for phase := 0; phase < 2; phase++ { // pre-auth and pre-hook both write back
		req := &schemas.HTTPRequest{
			Method:  string(ctx.Method()),
			Path:    string(ctx.Path()),
			Headers: map[string]string{"Cookie": string(ctx.Request.Header.Peek("Cookie"))},
		}
		if !applyHTTPRequestToCtx(&ctx, req) {
			t.Fatal("applyHTTPRequestToCtx rejected an unchanged request")
		}
	}

	got := ""
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		if strings.EqualFold(string(k), "cookie") {
			got += string(v)
		}
	})
	if got != "__cflb=a; __cf_bm=b" {
		t.Fatalf("Cookie after two write-backs = %q, want it unchanged", got)
	}
}
