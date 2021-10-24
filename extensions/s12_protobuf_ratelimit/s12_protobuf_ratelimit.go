package extensions

import (
	ratelimit "github.com/SafetyCulture/protoc-gen-ratelimit/proto"
	"github.com/pseudomuto/protoc-gen-doc/extensions"
)

// This is needed so protoc-gen-doc will transform our extension
func init() {
	extensions.SetTransformer("s12.protobuf.ratelimit.limit", func(payload interface{}) interface{} {
		ratelimit, ok := payload.(*ratelimit.MethodOptionsRateLimits)
		if !ok {
			return nil
		}

		return ratelimit
	})
}

// This is needed so protoc-gen-doc will transform our extension
func init() {
	extensions.SetTransformer("s12.protobuf.ratelimit.api_limit", func(payload interface{}) interface{} {
		ratelimit, ok := payload.(*ratelimit.ServiceOptionsRateLimits)
		if !ok {
			return nil
		}

		return ratelimit
	})
}
