package genratelimit

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	tmpl "text/template"

	"github.com/SafetyCulture/s12-apis-go/common"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	httpext "github.com/pseudomuto/protoc-gen-doc/extensions/google_api_http"
)

//go:embed templates/bucketer.lua.tmpl
var bucketerTemplate string

var paramMatch = regexp.MustCompile(`({\w+})`)

type pattern struct {
	Path   string
	Bucket string
}

func GenerateLuaBucketer(template *gendoc.Template) ([]byte, error) {
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
			servicePath := getServicePath(file, service)

			for _, method := range service.Methods {
				bucket := servicePath
				defaultPath := getDefaultMethodPath(file, service, method)

				if opts, ok := method.Option("s12.common.ratelimit").(*common.RateLimits); ok {
					if opts.Limits != nil {
						bucket = defaultPath
					}

					if opts.Bucket != common.RateLimitBucket_RATE_LIMIT_BUCKET_UNSPECIFIED {
						// We use the number as the strings value is likely out of date
						bucket = fmt.Sprintf("s12.common.ratelimit:%d", opts.Bucket.Number())
					}
				}

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
