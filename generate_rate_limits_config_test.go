package genratelimit_test

import (
	"os"
	"path/filepath"
	"testing"

	genratelimit "github.com/SafetyCulture/protoc-gen-ratelimit"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:generate buf build -o fixtures/image.bin

func TestGenerateRateLimitsConfig(t *testing.T) {
	f, err := os.ReadFile(filepath.Join("fixtures", "image.bin"))
	require.NoError(t, err)

	fileDesc := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(f, fileDesc)
	require.NoError(t, err)

	opts := &protogen.Options{}
	plugin, err := opts.New(&pluginpb.CodeGeneratorRequest{
		SourceFileDescriptors: fileDesc.GetFile(),
	})
	require.NoError(t, err)

	content, err := genratelimit.GenerateRateLimitsConfig(plugin, genratelimit.Config{
		Domain:      "my_domain",
		Descriptors: []string{"api_class", "bucket"},
		DefaultLimits: []genratelimit.Limit{
			{
				Key: "",
				Value: &genratelimit.YamlRateLimit{
					Unit:            "minute",
					RequestsPerUnit: 1,
				},
			},
		},
	})
	require.NoError(t, err)

	f, err = os.ReadFile("./fixtures/_generated/config.yaml")
	require.NoError(t, err)

	require.Equal(t, string(f), string(content))
}
