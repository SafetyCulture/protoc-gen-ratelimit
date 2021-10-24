package genratelimit_test

import (
	"bytes"
	"os"
	"testing"

	genratelimit "github.com/SafetyCulture/protoc-gen-ratelimit"
	_ "github.com/SafetyCulture/protoc-gen-ratelimit/extensions/s12_protobuf_ratelimit" // imported for side effects

	gendoc "github.com/pseudomuto/protoc-gen-doc"
	_ "github.com/pseudomuto/protoc-gen-doc/extensions/google_api_http" // imported for side effects
	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	"gotest.tools/assert"
)

//go:generate buf build -o fixtures/image.bin

func TestGenerateLuaBucketer(t *testing.T) {
	set, err := utils.LoadDescriptorSet("fixtures", "image.bin")
	assert.NilError(t, err)

	req := utils.CreateGenRequest(set, "fixtures/tasks.proto")
	result := protokit.ParseCodeGenRequest(req)

	template := gendoc.NewTemplate(result)

	content, err := genratelimit.GenerateLuaBucketer(template)
	assert.NilError(t, err)

	f, err := os.ReadFile("./fixtures/_generated/bucketer.lua")
	assert.NilError(t, err)

	var buf bytes.Buffer
	buf.Write(f)

	assert.Equal(t, buf.String(), string(content))
}
