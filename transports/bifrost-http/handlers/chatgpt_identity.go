package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
)

const codexIdentityHeader = "x-bf-codex-identity"

// Context comes from the swap proxy after account selection, not native clients.
type codexIdentity struct {
	Original     string `json:"original"`
	Mapped       string `json:"mapped"`
	Installation string `json:"installation"`
}

func rewriteCodexWS(raw []byte, auth *authHeaders) ([]byte, error) {
	values := auth.headers[codexIdentityHeader]
	if len(values) == 0 {
		return raw, nil
	}
	var identity codexIdentity
	if len(values) != 1 || len(values[0]) > 2048 || json.Unmarshal([]byte(values[0]), &identity) != nil || identity.Installation == "" || (identity.Original == "") != (identity.Mapped == "") {
		return nil, fmt.Errorf("invalid native identity context")
	}
	for _, value := range []string{identity.Original, identity.Mapped, identity.Installation} {
		if len(value) > 128 || strings.ContainsAny(value, "\r\n") {
			return nil, fmt.Errorf("invalid native identity context")
		}
	}
	out := raw
	set := func(path string, value string) error {
		var err error
		out, err = sjson.SetBytes(out, path, value)
		return err
	}
	replace := func(path string) error {
		v := gjson.GetBytes(out, path)
		if v.Type == gjson.String && identity.Original != "" && v.String() == identity.Original {
			return set(path, identity.Mapped)
		}
		return nil
	}
	for _, path := range []string{"prompt_cache_key", "client_metadata.session_id", "client_metadata.thread_id"} {
		if err := replace(path); err != nil {
			return nil, err
		}
	}
	window := func(value string) string {
		prefix := identity.Original + ":"
		if identity.Original != "" && strings.HasPrefix(value, prefix) {
			suffix := strings.TrimPrefix(value, prefix)
			if suffix != "" && strings.Trim(suffix, "0123456789") == "" {
				return identity.Mapped + ":" + suffix
			}
		}
		return value
	}
	v := gjson.GetBytes(out, "client_metadata.x-codex-window-id")
	if v.Type == gjson.String && window(v.String()) != v.String() {
		if err := set("client_metadata.x-codex-window-id", window(v.String())); err != nil {
			return nil, err
		}
	}
	if gjson.GetBytes(out, "client_metadata.x-codex-installation-id").Type == gjson.String {
		if err := set("client_metadata.x-codex-installation-id", identity.Installation); err != nil {
			return nil, err
		}
	}
	nested := gjson.GetBytes(out, "client_metadata.x-codex-turn-metadata")
	if nested.Type == gjson.String && gjson.Valid(nested.String()) {
		data := []byte(nested.String())
		before := string(data)
		for _, key := range []string{"session_id", "thread_id"} {
			v := gjson.GetBytes(data, key)
			if identity.Original != "" && v.Type == gjson.String && v.String() == identity.Original {
				var err error
				data, err = sjson.SetBytes(data, key, identity.Mapped)
				if err != nil {
					return nil, err
				}
			}
		}
		for _, key := range []string{"window_id", "installation_id"} {
			v := gjson.GetBytes(data, key)
			if v.Type != gjson.String {
				continue
			}
			value := window(v.String())
			if key == "installation_id" {
				value = identity.Installation
			}
			if value != v.String() {
				var err error
				data, err = sjson.SetBytes(data, key, value)
				if err != nil {
					return nil, err
				}
			}
		}
		if string(data) != before {
			if err := set("client_metadata.x-codex-turn-metadata", string(data)); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}
