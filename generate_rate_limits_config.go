package genratelimit

import (
	_ "embed"
	"fmt"

	"github.com/SafetyCulture/s12-apis-go/common"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"gopkg.in/yaml.v3"
)

type yamlRateLimit struct {
	RequestsPerUnit uint32 `yaml:"requests_per_unit,omitempty"`
	Unit            string `yaml:"unit,omitempty"`
	Unlimited       bool   `yaml:"unlimited,omitempty"`
}

type yamlDescriptor struct {
	Key         string
	Value       string           `yaml:"value,omitempty"`
	RateLimit   *yamlRateLimit   `yaml:"rate_limit,omitempty"`
	Descriptors []yamlDescriptor `yaml:"descriptors,omitempty"`
}

type yamlRoot struct {
	Domain      string
	Descriptors []yamlDescriptor
}

type override struct {
	OrgID     string         `yaml:"org_id"`
	RateLimit *yamlRateLimit `yaml:"rate_limit,omitempty"`
}

type config struct {
	Overrides []override
}

type apiLimit struct {
	Bucket string
	Limit  int32
}

func generateDescriptors(clientClass string, limits []apiLimit, overrideDescriptors []yamlDescriptor) yamlDescriptor {
	descriptors := make([]yamlDescriptor, 0)

	for _, limit := range limits {
		rateLimit := &yamlRateLimit{
			uint32(limit.Limit),
			"minute",
			false,
		}
		if limit.Limit == -1 {
			rateLimit = &yamlRateLimit{
				Unlimited: false,
			}
		}

		descriptors = append(descriptors, yamlDescriptor{
			Key:       "bucket",
			Value:     limit.Bucket,
			RateLimit: rateLimit,
		})
	}

	classDescriptor := newDescriptorTuple(clientClass, "", "", descriptors)
	// Overrides should be at the top of the descriptor list
	classDescriptor.Descriptors = append(overrideDescriptors, classDescriptor.Descriptors...)

	return classDescriptor
}

func GenerateRateLimitsConfig(template *gendoc.Template, cfg config) ([]byte, error) {
	var defaultLimit *common.Limits

	apiLimits := map[string][]apiLimit{}
	appendApiLimit := func(bucket string, limits *common.Limits) {
		def := limits.Default
		apiLimits["default"] = append(apiLimits["default"], apiLimit{bucket, def})

		if limits.Api != 0 {
			apiLimits["sc_api"] = append(apiLimits["sc_api"], apiLimit{bucket, limits.Api})
		}
		if limits.Integration != 0 {
			apiLimits["sc_integration"] = append(apiLimits["sc_integration"], apiLimit{bucket, limits.Integration})
		}

		if limits.Mobile != 0 {
			apiLimits["sc_device"] = append(apiLimits["sc_device"], apiLimit{bucket, limits.Mobile})
		}
		if limits.Web != 0 {
			apiLimits["sc_web"] = append(apiLimits["sc_web"], apiLimit{bucket, limits.Web})
		}

		if limits.Unauthenticated != 0 {
			apiLimits["unauthenticated"] = append(apiLimits["unauthenticated"], apiLimit{bucket, limits.Unauthenticated})
		}
	}

	for _, file := range template.Files {
		if file.Package == "s12.common" {
			for _, enum := range file.Enums {
				if enum.Name == "RateLimitBucket" {
					for _, value := range enum.Values {
						if opts, ok := value.Option("s12.common.limits").(*common.Limits); ok {
							if value.Number == "0" {
								defaultLimit = opts
							}
						}
						appendApiLimit(fmt.Sprintf("s12.common.ratelimit:%s", value.Number), defaultLimit)
					}
				}
			}
		}
	}

	for _, file := range template.Files {
		for _, service := range file.Services {
			servicePath := getServicePath(file, service)
			limit := defaultLimit

			if opts, ok := service.Option("s12.common.api_limits").(*common.ApiRateLimits); ok {
				if opts.Limits != nil && opts.Bucket != common.RateLimitBucket_RATE_LIMIT_BUCKET_UNSPECIFIED {
					return nil, fmt.Errorf("%s %s cannot use bucket and limits together", file.Name, service.FullName)
				}
				if opts.Limits != nil {
					limit = opts.Limits
				}
			}
			appendApiLimit(servicePath, limit)

			for _, method := range service.Methods {
				if opts, ok := method.Option("s12.common.ratelimit").(*common.ApiRateLimits); ok {
					if opts.Limits != nil && opts.Bucket != common.RateLimitBucket_RATE_LIMIT_BUCKET_UNSPECIFIED {
						return nil, fmt.Errorf("%s %s %s cannot use bucket and limits together", file.Name, service.FullName, method.Name)
					}
					if opts.Limits != nil {
						appendApiLimit(getDefaultMethodPath(file, service, method), opts.Limits)
					}
				}
			}
		}
	}

	overrideDescriptors := []yamlDescriptor{}
	for _, override := range cfg.Overrides {
		overrideDescriptors = append(
			overrideDescriptors,
			newDescriptorTuple("", override.OrgID, "", []yamlDescriptor{
				{
					Key:       "bucket",
					RateLimit: override.RateLimit,
				},
			}).Descriptors...,
		)
	}

	root := yamlRoot{
		Domain: "rate_per_user_bucket",
		Descriptors: []yamlDescriptor{
			generateDescriptors("sc_api", apiLimits["sc_api"], overrideDescriptors),
			generateDescriptors("sc_integration", apiLimits["sc_integration"], overrideDescriptors),
			generateDescriptors("sc_web", apiLimits["sc_web"], overrideDescriptors),
			generateDescriptors("sc_mobile", apiLimits["sc_mobile"], overrideDescriptors),
			generateDescriptors("unauthenticated", apiLimits["unauthenticated"], overrideDescriptors),
			generateDescriptors("", apiLimits["default"], overrideDescriptors),
		},
	}

	return yaml.Marshal(root)
}
