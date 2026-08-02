//go:build gui

package main

import (
	"fmt"
	"os"

	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/engine/ebitenhost"
)

// runWindowed opens a window. Built only with the `gui` tag, because the
// host needs cgo and system OpenGL/X11 libraries (ADR-0001).
func runWindowed(m *engine.Machine) error {
	fmt.Fprintln(os.Stdout, "controls:")
	fmt.Fprintln(os.Stdout, ebitenhost.Controls)
	return ebitenhost.Run(m, "FF6 Reconstruction — DEMO-0001")
}

func windowedInput() engine.Source { return ebitenhost.NewKeyboardSource() }
