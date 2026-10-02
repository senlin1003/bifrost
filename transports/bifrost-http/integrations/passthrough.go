package integrations

import (
	"os"
	"strings"

	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"
	"github.com/maximhq/bifrost/transports/bifrost-http/lib"
	"github.com/valyala/fasthttp"
)

// IsPassthroughRequest matches the native routes before the outer middlewares run.
// ChatGPT has a single allowed endpoint, unlike the other catch-all routers.
func IsPassthroughRequest(ctx *fasthttp.RequestCtx) bool {
	path := string(ctx.Path())
	if path == "/chatgpt_passthrough/backend-api/codex/responses" {
		return string(ctx.Method()) == fasthttp.MethodPost
	}
	for _, prefix := range []string{"/anthropic_passthrough/", "/openai_passthrough/", "/azure_passthrough/", "/runware_passthrough/", "/genai_passthrough/"} {
		if strings.HasPrefix(path, prefix) {
			switch string(ctx.Method()) {
			case fasthttp.MethodGet, fasthttp.MethodPost, fasthttp.MethodPut, fasthttp.MethodDelete, fasthttp.MethodPatch, fasthttp.MethodHead:
				return true
			}
			return false
		}
	}
	return false
}

// PassthroughRouter is a catch-all router that forwards all requests directly
// to the provider without matching against known route patterns.
type PassthroughRouter struct {
	*GenericRouter
}

// NewPassthroughRouter creates a passthrough-only router for any prefix/provider combo.
func NewPassthroughRouter(
	client *bifrost.Bifrost,
	handlerStore lib.HandlerStore,
	accessResolver AccessResolver,
	logger schemas.Logger,
	cfg *PassthroughConfig,
) *PassthroughRouter {
	if cfg == nil {
		cfg = &PassthroughConfig{}
	}
	return &PassthroughRouter{
		GenericRouter: NewGenericRouter(client, handlerStore, accessResolver, nil, cfg, logger),
	}
}

// NewAnthropicPassthroughRouter creates a passthrough router for /anthropic_passthrough.
func NewAnthropicPassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider: schemas.Anthropic,
		StripPrefix: []string{
			"/anthropic_passthrough",
		},
	})
}

// NewOpenAIPassthroughRouter creates a passthrough router for /openai_passthrough.
func NewOpenAIPassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider: schemas.OpenAI,
		StripPrefix: []string{
			"/openai_passthrough",
		},
	})
}

// NewChatGPTPassthroughRouter creates a passthrough router for /chatgpt_passthrough.
// Restricted to the Codex responses endpoint only — this is not a general-purpose
// ChatGPT backend proxy.
func NewChatGPTPassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider:    schemas.OpenAI,
		UpstreamURL: chatGPTUpstreamURL(),
		StripPrefix: []string{
			"/chatgpt_passthrough",
		},
		AllowedRoutes: []PassthroughRoute{
			{Method: fasthttp.MethodPost, Path: "/chatgpt_passthrough/backend-api/codex/responses"},
		},
	})
}

// ChatGPTUpstreamEnv overrides the /chatgpt_passthrough upstream (default https://chatgpt.com).
// It exists only for diagnostics: pointing it at a local plaintext recorder lets us capture
// exactly what Bifrost sends upstream and what ChatGPT returns, without trusting a MITM CA.
// Leave it unset in normal operation.
const ChatGPTUpstreamEnv = "BIFROST_CHATGPT_PASSTHROUGH_UPSTREAM"

const defaultChatGPTUpstream = "https://chatgpt.com"

func chatGPTUpstreamURL() string {
	if v := strings.TrimRight(strings.TrimSpace(os.Getenv(ChatGPTUpstreamEnv)), "/"); v != "" {
		return v
	}
	return defaultChatGPTUpstream
}

// NewAzurePassthroughRouter creates a passthrough router for /azure_passthrough.
func NewAzurePassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider: schemas.Azure,
		StripPrefix: []string{
			"/azure_passthrough",
		},
	})
}

// NewRunwarePassthroughRouter creates a passthrough router for /runware_passthrough. Runware exposes
// a single task-based endpoint, so this forwards raw task arrays and unlocks any Runware task type
// (3D, upscaling, background removal, ...) that Bifrost does not model natively.
func NewRunwarePassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider: schemas.Runware,
		StripPrefix: []string{
			"/runware_passthrough",
		},
	})
}

// NewGenAIPassthroughRouter creates a passthrough router for /genai_passthrough.
func NewGenAIPassthroughRouter(client *bifrost.Bifrost, handlerStore lib.HandlerStore, accessResolver AccessResolver, logger schemas.Logger) *PassthroughRouter {
	return NewPassthroughRouter(client, handlerStore, accessResolver, logger, &PassthroughConfig{
		Provider:         schemas.Gemini,
		ProviderDetector: detectProviderFromGenAIRequest,
		StripPrefix: []string{
			"/genai_passthrough/v1beta1",
			"/genai_passthrough/v1beta",
			"/genai_passthrough/v1",
			"/genai_passthrough",
		},
	})
}
