// protoc-gen-ratelimit is used to generate supporting files for https://github.com/envoyproxy/ratelimit.
//
// It is a protoc plugin, and can be invoked by passing `--ratelimit_out` and `--ratelimit_opt` arguments to protoc.
//
// Example: generate ratelimit configuration files for Envoy:
//
//	protoc --ratelimit_out=. --ratelimit_opt=config.yaml protos/*.proto
//
// For more details, check out the README at https://github.com/SafetyCulture/protoc-gen-ratelimit
package main

import (
	"os"

	genratelimit "github.com/SafetyCulture/protoc-gen-ratelimit"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	if flags := ParseFlags(os.Stdout, os.Args); HandleFlags(flags) {
		os.Exit(flags.Code())
	}

	// Run the plugin - it doesn't return anything
	protogen.Options{}.Run(genratelimit.Run)
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
