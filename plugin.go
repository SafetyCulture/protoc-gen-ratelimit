package genratelimit

import (
	"os"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// PluginOptions encapsulates options for the plugin. The type of renderer, template file, and the name of the output
// file are included.
type PluginOptions struct {
	ConfigFile string
}

// Config is the configuration of the plugin
type Config struct {
	Domain        string   `yaml:"domain"`
	Descriptors   []string `yaml:"descriptors"`
	DefaultLimits []Limit  `yaml:"default_limits"`
	Delimiter     string   `yaml:"delimiter"`
}

var delimiter = "|"

// SupportedFeatures describes a flag setting for supported features.
var SupportedFeatures = uint64(plugin_go.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// Plugin describes a protoc code generate plugin. It's an implementation of Plugin from github.com/pseudomuto/protokit
type Plugin struct{}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *Plugin) Generate(r *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	options, err := ParseOptions(r)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(options.ConfigFile)
	if err != nil {
		return nil, err
	}

	var configYaml Config
	err = yaml.NewDecoder(f).Decode(&configYaml)
	if err != nil {
		return nil, err
	}
	if configYaml.Delimiter != "" {
		delimiter = configYaml.Delimiter
	}

	result := protokit.ParseCodeGenRequest(r)

	resp := new(plugin_go.CodeGeneratorResponse)
	template := gendoc.NewTemplate(result)

	luaOutput, err := GenerateLuaBucketer(template)
	if err != nil {
		return nil, err
	}

	yamlOutput, err := GenerateRateLimitsConfig(template, configYaml)
	if err != nil {
		return nil, err
	}

	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String("ratelimit_bucketer.lua"),
		Content: proto.String(string(luaOutput)),
	})

	resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String("config.yaml"),
		Content: proto.String(string(yamlOutput)),
	})

	resp.SupportedFeatures = proto.Uint64(SupportedFeatures)

	return resp, nil
}

// ParseOptions parses plugin options from a CodeGeneratorRequest. It does this by splitting the `Parameter` field from
// the request object and parsing out the type of renderer to use and the name of the file to be generated.
//
// The parameter (`--doc_opt`) must be of the format <TYPE|TEMPLATE_FILE>,<OUTPUT_FILE>[,default|source_relative]:<EXCLUDE_PATTERN>,<EXCLUDE_PATTERN>*.
// The file will be written to the directory specified with the `--doc_out` argument to protoc.
func ParseOptions(req *plugin_go.CodeGeneratorRequest) (*PluginOptions, error) {
	options := &PluginOptions{
		ConfigFile: "config.yaml",
	}

	params := req.GetParameter()
	if params == "" {
		return options, nil
	}

	options.ConfigFile = params

	return options, nil
}
