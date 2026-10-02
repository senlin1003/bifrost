package integrations

import "testing"

// 验收: FWD-65-A2
// The caller's cookie reaches the provider, as on a direct connection; gateway credentials and
// hop-by-hop headers still do not.
func TestForwardPassthroughRequestHeaderKeepsCookie(t *testing.T) {
	for key, want := range map[string]bool{
		"cookie":              true,
		"user-agent":          true,
		"accept-encoding":     true,
		"originator":          true,
		"x-api-key":           false,
		"api-key":             false,
		"x-goog-api-key":      false,
		"proxy-authorization": false,
		"host":                false,
		"connection":          false,
		"transfer-encoding":   false,
		"set-cookie":          false,
		"x-request-id":        false,
		"x-bf-vk":             false,
	} {
		if got := forwardPassthroughRequestHeader(key); got != want {
			t.Errorf("forwardPassthroughRequestHeader(%q) = %v, want %v", key, got, want)
		}
	}
}
