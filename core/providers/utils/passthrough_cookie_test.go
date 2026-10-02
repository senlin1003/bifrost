package utils

import (
	"bufio"
	"strings"
	"testing"

	"github.com/valyala/fasthttp"
)

// What chatgpt.com/backend-api/codex/responses sets on a direct connection (names from the
// 2026-10-02 capture; values made up).
const codexSetCookieReply = "HTTP/1.1 200 OK\r\n" +
	"Set-Cookie: __cf_bm=abc; path=/; expires=Fri, 02-Oct-26 05:00:00 GMT; domain=chatgpt.com; HttpOnly; Secure; SameSite=None\r\n" +
	"Set-Cookie: __cflb=def; SameSite=None; Secure; path=/; expires=Fri, 02-Oct-26 06:00:00 GMT; HttpOnly\r\n" +
	"Set-Cookie: __oailb=ghi; path=/; HttpOnly; Secure\r\n" +
	"Content-Length: 0\r\n\r\n"

// 验收: FWD-65-A1
// Every Set-Cookie the provider sends is kept -- they cannot be comma-joined into one value.
func TestExtractPassthroughKeepsEverySetCookie(t *testing.T) {
	var resp fasthttp.Response
	if err := resp.Read(bufio.NewReader(strings.NewReader(codexSetCookieReply))); err != nil {
		t.Fatal(err)
	}
	stored := ExtractPassthroughProviderResponseHeaders(&resp)[PassthroughSetCookieHeader]
	lines := strings.Split(stored, "\n")
	if len(lines) != 3 {
		t.Fatalf("stored %d Set-Cookie values, want 3: %q", len(lines), stored)
	}
	for _, name := range []string{"__cf_bm=abc", "__cflb=def", "__oailb=ghi"} {
		if !strings.Contains(stored, name) {
			t.Fatalf("missing %s in %q", name, stored)
		}
	}
}

// 验收: FWD-65-A1
// Written back to the caller, each cookie is its own header line with its attributes intact.
func TestWritePassthroughSetCookiesWritesEachLine(t *testing.T) {
	var resp fasthttp.Response
	if err := resp.Read(bufio.NewReader(strings.NewReader(codexSetCookieReply))); err != nil {
		t.Fatal(err)
	}
	stored := ExtractPassthroughProviderResponseHeaders(&resp)[PassthroughSetCookieHeader]

	var out fasthttp.ResponseHeader
	WritePassthroughSetCookies(&out, stored)
	raw := out.String()
	if got := strings.Count(strings.ToLower(raw), "set-cookie:"); got != 3 {
		t.Fatalf("wrote %d Set-Cookie lines, want 3:\n%s", got, raw)
	}
	lower := strings.ToLower(raw)
	for _, want := range []string{"__cf_bm=abc", "domain=chatgpt.com", "httponly", "secure", "samesite=none", "__cflb=def", "__oailb=ghi"} {
		if !strings.Contains(lower, strings.ToLower(want)) {
			t.Fatalf("missing %q in written headers:\n%s", want, raw)
		}
	}
}

// 验收: FWD-65-A3
// A reply without cookies (Anthropic never sets any) gains no Set-Cookie.
func TestExtractPassthroughWithoutCookiesAddsNone(t *testing.T) {
	var resp fasthttp.Response
	resp.Header.Set("Content-Type", "text/event-stream")
	if _, ok := ExtractPassthroughProviderResponseHeaders(&resp)[PassthroughSetCookieHeader]; ok {
		t.Fatal("Set-Cookie appeared on a reply that had none")
	}
}
