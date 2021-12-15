// protoc-gen-ratelimit is used to generate supporting files for https://github.com/envoyproxy/ratelimit.
//
// It is a protoc plugin, and can be invoked by passing `--ratelimit_out` and `--ratelimit_opt` arguments to protoc.
//
// Example: generate ratelimit configuration files for Envoy:
//
//     protoc --ratelimit_out=. --doc_opt=ratelimit/config.yaml protos/*.proto
//
// For more details, check out the README at https://github.com/SafetyCulture/protoc-gen-ratelimit
package main

import (
	"github.com/pseudomuto/protokit"

	"log"
	"os"

	genratelimit "github.com/SafetyCulture/protoc-gen-ratelimit"

	_ "github.com/SafetyCulture/protoc-gen-ratelimit/extensions/s12_protobuf_ratelimit" // imported for side effects
	_ "github.com/pseudomuto/protoc-gen-doc/extensions/google_api_http"                 // imported for side effects
)

func main() {
	if flags := ParseFlags(os.Stdout, os.Args); HandleFlags(flags) {
		os.Exit(flags.Code())
	}

	if err := protokit.RunPlugin(new(genratelimit.Plugin)); err != nil {
		log.Fatal(err)
	}
}

// HandleFlags checks if there's a match and returns true if it was "handled"
func HandleFlags(f *Flags) bool {
	if !f.HasMatch() {
		return false
	}

	if f.ShowHelp() {
		f.PrintHelp()
	}

	if f.ShowVersion() {
		f.PrintVersion()
	}

	return true
}
