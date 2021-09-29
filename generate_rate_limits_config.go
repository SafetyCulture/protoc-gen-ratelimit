package genratelimit

import (
	_ "embed"
	"fmt"

	ratelimit "github.com/SafetyCulture/protoc-gen-ratelimit/proto"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"gopkg.in/yaml.v3"
)

// The rate limit being applied to a descriptor
type YamlRateLimit struct {
	RequestsPerUnit uint32 `yaml:"requests_per_unit,omitempty"`
	Unit            string `yaml:"unit,omitempty"`
	Unlimited       bool   `yaml:"unlimited,omitempty"`
}

// The description of a rate limiting tuple
type YamlDescriptor struct {
	Key         string
	Value       string            `yaml:"value,omitempty"`
	RateLimit   *YamlRateLimit    `yaml:"rate_limit,omitempty"`
	Descriptors []*YamlDescriptor `yaml:"descriptors,omitempty"`
}

// The root of the YAML document
type YamlRoot struct {
	Domain      string
	Descriptors []*YamlDescriptor
}

func GenerateRateLimitsConfig(template *gendoc.Template, cfg Config) ([]byte, error) {
	descriptors := cfg.Descriptors
	descriptorCount := len(descriptors)

	limits := map[string]*Limit{}
	for _, def := range cfg.DefaultLimits {
		key, err := formatKey(def.Key, "", descriptorCount)
		if err != nil {
			return nil, err
		}
		limits[key] = &Limit{
			Key:   key,
			Value: def.Value,
		}
	}

	for _, file := range template.Files {
		for _, service := range file.Services {
			servicePath := getServicePath(file, service)

			if opts, ok := service.Option("s12.protobuf.ratelimit.api_limit").(*ratelimit.ServiceOptionsRateLimits); ok {
				if opts.Limits != nil && opts.Bucket != "" {
					return nil, fmt.Errorf("%s %s cannot use bucket and limits together", file.Name, service.FullName)
				}
				if opts.Limits != nil {
					for key, value := range opts.Limits {
						limitKey, err := formatKey(key, servicePath, descriptorCount)
						if err != nil {
							return nil, err
						}

						limits[limitKey] = &Limit{
							limitKey,
							&YamlRateLimit{
								uint32(value.RequestsPerUnit),
								value.Unit,
								value.Unlimited,
							},
						}
					}
				}
			}

			for _, method := range service.Methods {
				if opts, ok := method.Option("s12.protobuf.ratelimit.limit").(*ratelimit.MethodOptionsRateLimits); ok {
					if opts.Limits != nil && opts.Bucket != "" {
						return nil, fmt.Errorf("%s %s %s cannot use bucket and limits together", file.Name, service.FullName, method.Name)
					}
					if opts.Limits != nil {
						for key, value := range opts.Limits {
							limitKey, err := formatKey(key, getDefaultMethodPath(file, service, method), descriptorCount)
							if err != nil {
								return nil, err
							}

							limits[limitKey] = &Limit{
								limitKey,
								&YamlRateLimit{
									uint32(value.RequestsPerUnit),
									value.Unit,
									value.Unlimited,
								},
							}
						}
					}
				}
			}
		}
	}

	// Sort the limits so that the output is deterministic
	// and empty/default values are last and not immediately matched
	sortedLimits := sortLimits(limits)

	root := YamlRoot{
		Domain:      "rate_per_user_bucket",
		Descriptors: sortedLimits.Descriptors(cfg.Descriptors),
	}

	return yaml.Marshal(root)
}
