package route

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

// The Mesen probes execute the routes; the model constructors are the
// tracked encodings of them. Multiple sources of truth drift silently, so
// this test parses each probe's ROUTE table and asserts the leg sequence
// matches its model. If a probe is intentionally changed, update the model
// in the same commit. EXP-0037 carries a copy of the EXP-0036 route
// (instrumentation-only additions); EXP-0038 extends it with the mines
// corridor legs.
var probeModels = []struct {
	path  string
	model func() Route
}{
	{"../../../mesen/probes/EXP-0036.lua", MinesRoute},
	{"../../../mesen/probes/EXP-0037.lua", MinesRoute},
	{"../../../mesen/probes/EXP-0038.lua", MinesEncounterRoute},
}

var (
	legLineRe = regexp.MustCompile(`\{\s*n\s*=\s*(\d+),\s*kind\s*=\s*"(\w+)"(.*)$`)
	fieldRe   = regexp.MustCompile(`(\w+)\s*=\s*(?:"([^"]*)"|(0x[0-9A-Fa-f]+|\d+))`)
)

type probeLeg struct {
	num     int
	kind    string
	dir     string
	target  int64
	cmp     string
	until   string
	hasTgt  bool
	dur     int64
	timeout int64
}

func parseProbeRoute(t *testing.T, probePath string) []probeLeg {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(probePath))
	if err != nil {
		t.Fatalf("reading probe: %v", err)
	}
	var legs []probeLeg
	for _, line := range regexp.MustCompile(`\r?\n`).Split(string(data), -1) {
		m := legLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		num, _ := strconv.Atoi(m[1])
		leg := probeLeg{num: num, kind: m[2]}
		for _, f := range fieldRe.FindAllStringSubmatch(m[3], -1) {
			key, str, num := f[1], f[2], f[3]
			switch key {
			case "dir":
				leg.dir = str
			case "cmp":
				leg.cmp = str
			case "until_":
				leg.until = str
			case "target":
				v, err := strconv.ParseInt(num, 0, 64)
				if err == nil {
					leg.target, leg.hasTgt = v, true
				}
			case "dur":
				if v, err := strconv.ParseInt(num, 0, 64); err == nil {
					leg.dur = v
				}
			case "timeout":
				if v, err := strconv.ParseInt(num, 0, 64); err == nil {
					leg.timeout = v
				}
			}
		}
		legs = append(legs, leg)
	}
	return legs
}

func TestProbeRouteMatchesModel(t *testing.T) {
	for _, pm := range probeModels {
		t.Run(filepath.Base(pm.path), func(t *testing.T) {
			assertProbeMatchesModel(t, pm.path, pm.model())
		})
	}
}

func assertProbeMatchesModel(t *testing.T, probePath string, model Route) {
	probe := parseProbeRoute(t, probePath)

	if len(probe) != len(model.Legs) {
		t.Fatalf("probe has %d legs, model has %d", len(probe), len(model.Legs))
	}

	untilName := map[Until]string{
		UntilBattleStart: "battle_start",
		UntilBattleEnd:   "battle_end",
		UntilElapsed:     "elapsed",
		UntilWatchedByte: "watched_byte",
	}
	cmpName := map[Compare]string{AtLeast: "ge", AtMost: "le"}

	for i, want := range model.Legs {
		got := probe[i]
		if got.num != want.Num {
			t.Errorf("leg index %d: probe n=%d, model Num=%d", i, got.num, want.Num)
			continue
		}
		wantKind := "walk"
		if want.Pulse {
			wantKind = "pulse"
		}
		if got.kind != wantKind {
			t.Errorf("leg %d: probe kind %q, model %q", want.Num, got.kind, wantKind)
		}
		if !want.Pulse && got.dir != string(want.Input) {
			t.Errorf("leg %d: probe dir %q, model %q", want.Num, got.dir, want.Input)
		}
		if want.Until == UntilPosition {
			if !got.hasTgt || uint8(got.target) != want.Target {
				t.Errorf("leg %d: probe target %#x (present=%v), model %#x",
					want.Num, got.target, got.hasTgt, want.Target)
			}
			if got.cmp != cmpName[want.Compare] {
				t.Errorf("leg %d: probe cmp %q, model %q", want.Num, got.cmp, cmpName[want.Compare])
			}
		} else if got.until != untilName[want.Until] {
			t.Errorf("leg %d: probe until_ %q, model %q", want.Num, got.until, untilName[want.Until])
		}
		// Timing drift between the encodings is as real a divergence as a
		// wrong target: the EXP-0036 audit found the guard compared shape
		// but not durations, so a timeout edit in one file would have
		// passed silently.
		if int64(want.Timeout) != got.timeout {
			t.Errorf("leg %d: probe timeout %d, model %d", want.Num, got.timeout, want.Timeout)
		}
		if want.Until == UntilElapsed && int64(want.Duration) != got.dur {
			t.Errorf("leg %d: probe dur %d, model %d", want.Num, got.dur, want.Duration)
		}
	}
}
