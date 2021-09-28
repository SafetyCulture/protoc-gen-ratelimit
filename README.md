# protoc-gen-ratelimit

This is a Rate Limit generator for Google Protocol Buffers compiler `protoc`. The plugin generates a Lua filter to bucket requests based on their paths, and a descriptor file for [envoyproxy/ratelimit](https://github.com/envoyproxy/ratelimit).

## Installation

```
go get -u github.com/SafetyCulture/protoc-gen-ratelimit/cmd/protoc-gen-ratelimit
```

## Invoking the Plugin

The plugin is invoked by passing the --ratelimit_out, and --ratelimit_opt options to the protoc compiler. The option has the following format:

```
--doc_opt=ratelimit/config.yaml
```
