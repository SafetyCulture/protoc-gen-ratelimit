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

func TestGenerateLuaBucketer(t *testing.T) {
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
	generated, err := genratelimit.GenerateLuaBucketer(plugin)
	require.NoError(t, err)

	f, err = os.ReadFile("./fixtures/_generated/bucketer.lua")
	require.NoError(t, err)

	require.Equal(t, string(f), string(generated))
}
