package genratelimit

import (
	"bytes"
	// For embedded templates
	_ "embed"
	"fmt"
	"regexp"
	tmpl "text/template"

	ratelimit "github.com/SafetyCulture/protoc-gen-ratelimit/s12/protobuf/ratelimit"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	httpext "github.com/pseudomuto/protoc-gen-doc/extensions/google_api_http"
)

//go:embed templates/bucketer.lua.tmpl
var bucketerTemplate string

// Used to identify parameters in a path e.g. `/users/{used_id}`
var paramMatch = regexp.MustCompile(`({\w+})`)

// A path with an associated bucket
type pattern struct {
	Path   string
	Bucket string
}

// GenerateLuaBucketer generates the Lua bucketer
func GenerateLuaBucketer(template *gendoc.Template) ([]byte, error) {
	// All of the paths grouped by method, and if they include path parameters
	// Paths with parameters require pattern matching.
	httpMethods := map[string]map[bool][]pattern{}
	addPattern := func(method, path, bucket string, hasParams bool) {
		if httpMethods[method] == nil {
			httpMethods[method] = map[bool][]pattern{}
		}
		if httpMethods[method][hasParams] == nil {
			httpMethods[method][hasParams] = []pattern{}
		}

		httpMethods[method][hasParams] = append(httpMethods[method][hasParams], pattern{path, bucket})
	}

	for _, file := range template.Files {
		for _, service := range file.Services {
			for _, method := range service.Methods {
				bucket := service.FullName
				defaultPath := getDefaultMethodPath(service, method)

				if opts, ok := method.Option("s12.protobuf.ratelimit.limit").(*ratelimit.MethodOptionsRateLimits); ok {
					if opts.Limits != nil {
						bucket = defaultPath
					}

					if opts.Bucket != "" {
						// We use the number as the strings value is likely out of date
						bucket = fmt.Sprintf("custom_bucket:%s", opts.Bucket)
					}
				}

				// Default method/path combination
				addPattern("POST", defaultPath, bucket, false)

				if opts, ok := method.Option("google.api.http").(httpext.HTTPExtension); ok {
					for _, rule := range opts.Rules {
						pathWithoutParams := paramMatch.ReplaceAllString(rule.Pattern, ".+")
						hasParams := pathWithoutParams != rule.Pattern

						addPattern(rule.Method, pathWithoutParams, bucket, hasParams)
					}
				}
			}
		}
	}

	tp, err := tmpl.New("Lua Template").Parse(bucketerTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tp.Execute(&buf, httpMethods)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
