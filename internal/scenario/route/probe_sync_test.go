package route

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

// The Mesen probe executes the route; MinesRoute() is the tracked model of
// it. Two sources of truth drift silently, so this test parses the probe's
// ROUTE table and asserts the leg sequence matches. If the probe is
// intentionally changed, update MinesRoute() in the same commit.
const probePath = "../../../mesen/probes/EXP-0036.lua"

var (
	legLineRe = regexp.MustCompile(`\{\s*n\s*=\s*(\d+),\s*kind\s*=\s*"(\w+)"(.*)$`)
	fieldRe   = regexp.MustCompile(`(\w+)\s*=\s*(?:"([^"]*)"|(0x[0-9A-Fa-f]+|\d+))`)
)

type probeLeg struct {
	num    int
	kind   string
	dir    string
	target int64
	cmp    string
	until  string
	hasTgt bool
}

func parseProbeRoute(t *testing.T) []probeLeg {
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
			}
		}
		legs = append(legs, leg)
	}
	return legs
}

func TestProbeRouteMatchesModel(t *testing.T) {
	probe := parseProbeRoute(t)
	model := MinesRoute()

	if len(probe) != len(model.Legs) {
		t.Fatalf("probe has %d legs, model has %d", len(probe), len(model.Legs))
	}

	untilName := map[Until]string{
		UntilBattleStart: "battle_start",
		UntilBattleEnd:   "battle_end",
		UntilElapsed:     "elapsed",
		UntilMapChange:   "map_change",
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
	}
}
