package genratelimit

import (
	"bytes"
	// For embedded templates
	_ "embed"
	"fmt"
	"regexp"
	tmpl "text/template"

	ratelimit "github.com/SafetyCulture/protoc-gen-ratelimit/s12/protobuf/ratelimit"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
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

func httpRuleToPattern(rule *annotations.HttpRule, bucket string) (string, string, string, bool) {
	var method string
	var path string

	switch rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		method = "GET"
		path = rule.GetGet()
	case *annotations.HttpRule_Put:
		method = "PUT"
		path = rule.GetPut()
	case *annotations.HttpRule_Post:
		method = "POST"
		path = rule.GetPost()
	case *annotations.HttpRule_Delete:
		method = "DELETE"
		path = rule.GetDelete()
	case *annotations.HttpRule_Patch:
		method = "PATCH"
		path = rule.GetPatch()
	}

	pathPattern := paramMatch.ReplaceAllString(path, ".+")
	return method, pathPattern, bucket, pathPattern != path
}

// GenerateLuaBucketer generates the Lua bucketer
func GenerateLuaBucketer(plugin *protogen.Plugin) ([]byte, error) {
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

	for _, file := range plugin.Request.SourceFileDescriptors {
		for _, service := range file.GetService() {
			for _, method := range service.GetMethod() {
				bucket := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())
				defaultPath := fmt.Sprintf("/%s.%s/%s", file.GetPackage(), service.GetName(), method.GetName())

				if methodOpts, ok := proto.GetExtension(method.GetOptions(), ratelimit.E_Limit).(*ratelimit.MethodOptionsRateLimits); ok && methodOpts != nil {
					if methodOpts.Limits != nil {
						bucket = defaultPath
					}

					if methodOpts.Bucket != "" {
						// We use the number as the strings value is likely out of date
						bucket = fmt.Sprintf("custom_bucket:%s", methodOpts.Bucket)
					}
				}

				// Default method/path combination
				addPattern("POST", defaultPath, bucket, false)

				// Extract HTTP rules
				if httpOpts, ok := proto.GetExtension(method.GetOptions(), annotations.E_Http).(*annotations.HttpRule); ok {
					addPattern(httpRuleToPattern(httpOpts, bucket))

					for _, rule := range httpOpts.AdditionalBindings {
						addPattern(httpRuleToPattern(rule, bucket))
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
