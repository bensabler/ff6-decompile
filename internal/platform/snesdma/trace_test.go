package snesdma

import (
	"strings"
	"testing"
)

// A line in exactly the shape mesen/probes/dma-trace.lua writes: an 8 KB
// A->B transfer from $7EC180 to $2118, mode 1.
const sampleLine = "dma frame=51578 pc=$C0A1B2 reg=$420B mask=$01 ch=0 regs=011880C17E00200000000000"

func TestParseTraceLine(t *testing.T) {
	e, ok, err := ParseTraceLine(sampleLine)
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	if e.Frame != 51578 {
		t.Errorf("Frame = %d, want 51578", e.Frame)
	}
	if e.PC != 0xC0A1B2 {
		t.Errorf("PC = $%06X, want $C0A1B2", e.PC)
	}
	if e.Register != "$420B" || e.HDMA() {
		t.Errorf("Register = %q, HDMA = %v", e.Register, e.HDMA())
	}
	if e.Channel.Index != 0 || e.Channel.Source() != 0x7EC180 {
		t.Errorf("channel %d source $%06X", e.Channel.Index, e.Channel.Source())
	}
	if !e.Channel.TargetsVRAM() {
		t.Error("dest $18 should classify as VRAM")
	}
	if got := e.Channel.ByteCount(); got != 0x2000 {
		t.Errorf("ByteCount = %d, want %d", got, 0x2000)
	}
	start, end, spanOK := e.Channel.SourceSpan()
	if !spanOK || start != 0x7EC180 || end != 0x7EE17F {
		t.Errorf("span $%06X-$%06X ok=%v, want $7EC180-$7EE17F", start, end, spanOK)
	}
}

func TestParseTraceLineIgnoresNonRecords(t *testing.T) {
	for _, line := range []string{
		"",
		"   ",
		"dma-trace: watching $420B/$420C, appending to mesen/out/dma-trace.log",
		"# a comment",
	} {
		_, ok, err := ParseTraceLine(line)
		if err != nil {
			t.Errorf("line %q errored: %v", line, err)
		}
		if ok {
			t.Errorf("line %q was parsed as a record", line)
		}
	}
}

func TestParseTraceLineRejectsMalformed(t *testing.T) {
	tests := []struct {
		name string
		line string
	}{
		{"no ch field", "dma frame=1 pc=$C00000 reg=$420B mask=$01 regs=011880C17E00200000000000"},
		{"no regs field", "dma frame=1 pc=$C00000 reg=$420B mask=$01 ch=0"},
		{"short regs", "dma frame=1 pc=$C00000 reg=$420B mask=$01 ch=0 regs=0118"},
		{"odd hex", "dma frame=1 pc=$C00000 reg=$420B mask=$01 ch=0 regs=011880C17E0020000000000"},
		{"channel out of range", "dma frame=1 pc=$C00000 reg=$420B mask=$FF ch=9 regs=011880C17E00200000000000"},
		{"unknown field", "dma frame=1 wat=3 pc=$C00000 reg=$420B mask=$01 ch=0 regs=011880C17E00200000000000"},
		{"field with no equals", "dma frame=1 lonely pc=$C00000"},
		{"bad frame", "dma frame=xx pc=$C00000 reg=$420B mask=$01 ch=0 regs=011880C17E00200000000000"},
		// The self-consistency check: the line claims channel 3 but the
		// mask only enables channel 0.
		{"mask disagrees with ch", "dma frame=1 pc=$C00000 reg=$420B mask=$01 ch=3 regs=011880C17E00200000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := ParseTraceLine(tt.line); err == nil {
				t.Error("want an error, got none")
			}
		})
	}
}

func TestParseTrace(t *testing.T) {
	log := strings.Join([]string{
		"dma-trace: watching $420B/$420C, appending to mesen/out/dma-trace.log",
		sampleLine,
		"dma frame=51579 pc=$C0A1C0 reg=$420B mask=$02 ch=1 regs=022200400000020000000000",
		"",
		"dma frame=51580 pc=$C0A200 reg=$420C mask=$80 ch=7 regs=410400500000000000001C00",
	}, "\n")

	events, err := ParseTrace(strings.NewReader(log))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("got %d events, want 3", len(events))
	}
	if !events[1].Channel.TargetsCGRAM() {
		t.Error("second event writes $2122 and should classify as CGRAM")
	}
	if !events[2].HDMA() {
		t.Error("third event came from $420C and should report HDMA")
	}
	if got := FilterTarget(events, DestVRAMLow); len(got) != 1 || got[0].Frame != 51578 {
		t.Errorf("FilterTarget(VRAM) = %v", got)
	}
}

// A partial trace analysed as if complete is how a missing transfer becomes a
// wrong conclusion, so parsing stops rather than skipping.
func TestParseTraceStopsAtTheFirstBadRecord(t *testing.T) {
	log := strings.Join([]string{
		sampleLine,
		"dma frame=2 pc=$C00000 reg=$420B mask=$01 ch=0 regs=zz",
		sampleLine,
	}, "\n")

	events, err := ParseTrace(strings.NewReader(log))
	if err == nil {
		t.Fatal("want an error for the malformed record")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error should name the line, got %v", err)
	}
	if len(events) != 1 {
		t.Errorf("want the events parsed before the failure, got %d", len(events))
	}
}

func FuzzParseTraceLine(f *testing.F) {
	f.Add(sampleLine)
	f.Add("dma")
	f.Add("")
	f.Fuzz(func(t *testing.T, line string) {
		e, ok, err := ParseTraceLine(line)
		if err != nil || !ok {
			return
		}
		// Anything reported as a valid record must be self-consistent.
		if e.Channel.Index < 0 || e.Channel.Index >= ChannelCount {
			t.Fatalf("channel index %d out of range", e.Channel.Index)
		}
		if e.Mask&(1<<uint(e.Channel.Index)) == 0 {
			t.Fatalf("channel %d not enabled in mask $%02X", e.Channel.Index, e.Mask)
		}
		if n := e.Channel.ByteCount(); n < 1 || n > 65536 {
			t.Fatalf("ByteCount %d out of range", n)
		}
	})
}
