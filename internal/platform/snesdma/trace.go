package snesdma

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Trace parsing for mesen/probes/dma-trace.lua.
//
// The probe logs raw register bytes rather than a decoded summary, so a log
// can be re-decoded when this package improves. A log of someone's
// interpretation cannot be, which is the whole argument for the split.

// Event is one channel enabled by one $420B/$420C write.
type Event struct {
	// Frame is the emulator frame count when the transfer was kicked.
	Frame uint64
	// PC is the 24-bit program counter that wrote the enable register.
	PC uint32
	// Register is "$420B" for DMA or "$420C" for HDMA.
	Register string
	// Mask is the value written, so a caller can see the other channels
	// enabled by the same write.
	Mask uint8
	// Channel is the decoded configuration for this channel.
	Channel Channel
}

// HDMA reports whether the event came from an HDMA enable.
func (e Event) HDMA() bool { return e.Register == "$420C" }

// ParseTraceLine parses one line of dma-trace.log.
//
// Lines that are not trace records — the probe's own startup message, blank
// lines, anything else that shares the file — return ok=false with no error.
// A line that looks like a record but is malformed returns an error, because
// silently skipping a corrupt record would understate what a trace saw.
func ParseTraceLine(line string) (Event, bool, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 || fields[0] != "dma" {
		return Event{}, false, nil
	}

	var e Event
	var chIndex = -1
	var regBytes []byte

	for _, f := range fields[1:] {
		k, v, found := strings.Cut(f, "=")
		if !found {
			return Event{}, false, fmt.Errorf("snesdma: field %q has no '='", f)
		}
		var err error
		switch k {
		case "frame":
			e.Frame, err = strconv.ParseUint(v, 10, 64)
		case "pc":
			var n uint64
			n, err = strconv.ParseUint(strings.TrimPrefix(v, "$"), 16, 32)
			e.PC = uint32(n)
		case "reg":
			e.Register = v
		case "mask":
			var n uint64
			n, err = strconv.ParseUint(strings.TrimPrefix(v, "$"), 16, 8)
			e.Mask = uint8(n)
		case "ch":
			chIndex, err = strconv.Atoi(v)
		case "regs":
			regBytes, err = hex.DecodeString(v)
		default:
			return Event{}, false, fmt.Errorf("snesdma: unknown field %q", k)
		}
		if err != nil {
			return Event{}, false, fmt.Errorf("snesdma: field %s=%q: %w", k, v, err)
		}
	}

	if chIndex < 0 {
		return Event{}, false, fmt.Errorf("snesdma: line has no ch= field")
	}
	if regBytes == nil {
		return Event{}, false, fmt.Errorf("snesdma: channel %d has no regs= field", chIndex)
	}
	ch, err := Decode(chIndex, regBytes)
	if err != nil {
		return Event{}, false, err
	}
	e.Channel = ch

	// The mask must actually enable the channel the line claims. A line
	// that disagrees with itself means the probe or the log is wrong, and
	// that is worth failing on rather than analysing.
	if e.Mask&(1<<uint(chIndex)) == 0 {
		return Event{}, false, fmt.Errorf("snesdma: channel %d is not set in mask $%02X", chIndex, e.Mask)
	}
	return e, true, nil
}

// ParseTrace reads a whole dma-trace.log.
//
// It stops at the first malformed record. A partial trace analysed as if
// complete is how a missing transfer becomes a wrong conclusion.
func ParseTrace(r io.Reader) ([]Event, error) {
	var out []Event
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	line := 0
	for sc.Scan() {
		line++
		e, ok, err := ParseTraceLine(sc.Text())
		if err != nil {
			return out, fmt.Errorf("line %d: %w", line, err)
		}
		if ok {
			out = append(out, e)
		}
	}
	if err := sc.Err(); err != nil {
		return out, fmt.Errorf("snesdma: reading trace: %w", err)
	}
	return out, nil
}

// FilterTarget returns the events whose channel writes the given destination
// register, for example DestVRAMLow.
func FilterTarget(events []Event, dest uint8) []Event {
	var out []Event
	for _, e := range events {
		if e.Channel.Dest == dest {
			out = append(out, e)
		}
	}
	return out
}
