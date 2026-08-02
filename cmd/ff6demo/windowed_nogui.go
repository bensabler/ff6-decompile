//go:build !gui

package main

import (
	"errors"

	"github.com/bensabler/ff6-decompile/internal/engine"
)

// runWindowed is unavailable in the default build.
//
// The windowed host needs cgo and system OpenGL/X11 libraries on Linux and
// macOS, so it sits behind the `gui` build tag. That keeps `go build ./...`,
// `go vet ./...` and `go test ./...` free of any system-library requirement,
// which is what lets CI run the full demo test surface on a bare container.
func runWindowed(*engine.Machine) error {
	return errors.New(`this binary was built without the windowed host.

Build with the gui tag to open a window:

  go run -tags gui ./cmd/ff6demo

Or run headless, which needs no display and is the authoritative mode for
automated comparison:

  go run ./cmd/ff6demo -headless -frames 120 -capture-last`)
}

// windowedInput is never reached without the gui tag; runWindowed errors
// first. It exists so both builds share one main.go.
func windowedInput() engine.Source { return engine.NoInput }
