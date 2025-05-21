package genratelimit

import (
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v3"
)

const Version = "1.1.0"

// Config is the configuration of the plugin
type Config struct {
	Domain        string   `yaml:"domain"`
	Descriptors   []string `yaml:"descriptors"`
	DefaultLimits []Limit  `yaml:"default_limits"`
	Delimiter     string   `yaml:"delimiter"`
}

var delimiter = "|"

// Run executes the plugin using protogen
func Run(gen *protogen.Plugin) error {
	// Get config file path from parameters
	configFile := "config.yaml"
	if gen.Request.GetParameter() != "" {
		configFile = gen.Request.GetParameter()
	}

	// Read config file
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var configYaml Config
	err = yaml.NewDecoder(f).Decode(&configYaml)
	if err != nil {
		return err
	}
	if configYaml.Delimiter != "" {
		delimiter = configYaml.Delimiter
	}

	// Generate outputs
	luaOutput, err := GenerateLuaBucketer(gen)
	if err != nil {
		return err
	}

	yamlOutput, err := GenerateRateLimitsConfig(gen, configYaml)
	if err != nil {
		return err
	}

	// Add outputs to response
	luaFile := gen.NewGeneratedFile("ratelimit_bucketer.lua", "")
	_, err = luaFile.Write(luaOutput)
	if err != nil {
		return err
	}

	yamlFile := gen.NewGeneratedFile("config.yaml", "")
	_, err = yamlFile.Write(yamlOutput)
	if err != nil {
		return err
	}

	return nil
}
