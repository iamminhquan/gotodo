// Package version holds the build-time version string.
package version

// Version is injected at build time via ldflags:
//
//	go build -ldflags "-X github.com/iamminhquan/gotodo/internal/version.Version=1.2.3" ./cmd/gotodo
var Version = "v1.0.0"
