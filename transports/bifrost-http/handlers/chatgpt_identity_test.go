package handlers

import (
	"bytes"
	"github.com/tidwall/gjson"
	"testing"
)

// 验收: FWD-67-A2~
func TestCodexWSIdentityPreservesBusinessBytes(t *testing.T) {
	raw := []byte(`{ "type":"response.create", "prompt_cache_key":"original", "input":[{"text":"original"}], "unknown":900719925474099312345, "client_metadata":{"session_id":"original","thread_id":"other","turn_id":"original","x-codex-window-id":"original:2","x-codex-installation-id":"local","x-codex-turn-metadata":"{\"session_id\":\"original\",\"installation_id\":\"local\",\"turn_id\":\"original\"}"} }`)
	auth := &authHeaders{headers: map[string][]string{codexIdentityHeader: {`{"original":"original","mapped":"mapped","installation":"company-device"}`}}}
	out, err := rewriteCodexWS(raw, auth)
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{"prompt_cache_key": "mapped", "client_metadata.session_id": "mapped", "client_metadata.thread_id": "other", "client_metadata.turn_id": "original", "client_metadata.x-codex-window-id": "mapped:2", "client_metadata.x-codex-installation-id": "company-device", "input.0.text": "original"} {
		if got := gjson.GetBytes(out, path).String(); got != want {
			t.Fatalf("%s: %s", path, got)
		}
	}
	if !bytes.Contains(out, []byte(`"unknown":900719925474099312345`)) || !bytes.HasPrefix(out, []byte(`{ "type"`)) {
		t.Fatal("unrelated bytes changed")
	}
	nested := gjson.GetBytes(out, "client_metadata.x-codex-turn-metadata").String()
	if gjson.Get(nested, "session_id").String() != "mapped" || gjson.Get(nested, "turn_id").String() != "original" {
		t.Fatal("nested metadata mismatch")
	}
	again, err := rewriteCodexWS(raw, auth)
	if err != nil || !bytes.Equal(out, again) {
		t.Fatal("mapping unstable")
	}
	other, err := rewriteCodexWS(raw, &authHeaders{headers: map[string][]string{codexIdentityHeader: {`{"original":"original","mapped":"other-account-session","installation":"other-account-device"}`}}})
	if err != nil || gjson.GetBytes(other, "prompt_cache_key").String() == gjson.GetBytes(out, "prompt_cache_key").String() {
		t.Fatal("accounts share identity")
	}
	if chatGPTWSHeaders(auth).Get(codexIdentityHeader) != "" {
		t.Fatal("internal context leaked")
	}
	untouched, err := rewriteCodexWS(raw, &authHeaders{headers: map[string][]string{}})
	if err != nil || !bytes.Equal(untouched, raw) {
		t.Fatal("legacy message changed")
	}
}

// 验收: FWD-67-A2~
func TestCodexWSIdentityRejectsInvalidContext(t *testing.T) {
	for _, value := range []string{`{`, `{}`, `{"original":"a","mapped":"","installation":"b"}`} {
		if _, err := rewriteCodexWS([]byte(`{}`), &authHeaders{headers: map[string][]string{codexIdentityHeader: {value}}}); err == nil {
			t.Fatal("invalid context accepted")
		}
	}
}
