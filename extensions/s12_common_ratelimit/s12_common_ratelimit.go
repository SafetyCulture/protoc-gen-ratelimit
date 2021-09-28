package extensions

import (
	"github.com/SafetyCulture/s12-apis-go/common"
	"github.com/pseudomuto/protoc-gen-doc/extensions"
)

// This is needed so protoc-gen-doc will transform our extension
func init() {
	extensions.SetTransformer("s12.common.ratelimit", func(payload interface{}) interface{} {
		ratelimit, ok := payload.(*common.RateLimits)
		if !ok {
			return nil
		}

		return ratelimit
	})
}

// This is needed so protoc-gen-doc will transform our extension
func init() {
	extensions.SetTransformer("s12.common.limits", func(payload interface{}) interface{} {
		ratelimit, ok := payload.(*common.Limits)
		if !ok {
			return nil
		}

		return ratelimit
	})
}

// This is needed so protoc-gen-doc will transform our extension
func init() {
	extensions.SetTransformer("s12.common.api_limits", func(payload interface{}) interface{} {
		ratelimit, ok := payload.(*common.ApiRateLimits)
		if !ok {
			return nil
		}

		return ratelimit
	})
}
