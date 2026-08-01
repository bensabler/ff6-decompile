package route

// MinesRoute is the SCN-0001 segment-4 route: milestone 04
// (`04-free-movement`, Narshe exterior) to milestone 05
// (`05-mines-entry`, mines interior).
//
// EXP-0035 walked and documented eleven legs interactively. Three of those
// combine a one-tile sidestep with a climb, and EXP-0036 found one omitted
// climb plus a walk-on to the documented tile, so this encoding splits them
// into seventeen legs, each with one input and one completion condition.
// The zigzag matters — a single held direction stalls at Y=$21 and again
// at Y=$16, which is why the earlier segments' hold-one-direction cadence
// could not walk this stretch.
//
// Expected coordinates come from the EXP-0035 recon and are re-asserted by
// EXP-0036's scheduled runs. StartFrame is milestone 04's capture frame.
func MinesRoute() Route {
	return Route{
		Name:          "SCN-0001 segment 4: milestone 04 to mines interior",
		StartFrame:    46375,
		NeutralFrames: 8,
		Legs: []Leg{
			{Num: 1, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x1B, Timeout: 900,
				ExpectX: 0x1B, ExpectY: 0x2A,
				Note: "leave the milestone-04 tile eastward"},
			{Num: 2, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x27, Timeout: 900,
				ExpectX: 0x1B, ExpectY: 0x27,
				Note: "first climb"},
			{Num: 3, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x1E, Timeout: 900,
				ExpectX: 0x1E, ExpectY: 0x27,
				Note: "east to the trigger column"},
			// EXP-0036 needed two runs to place the trigger. Standing on
			// (1E,27) and tapping A did nothing, and neither did holding
			// right there: the guard dialogue appears one tile further
			// north, at (1E,25). EXP-0035's condensed leg table had dropped
			// this intermediate climb; its own recon log has it.
			{Num: 4, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x25, Timeout: 900,
				ExpectX: 0x1E, ExpectY: 0x25,
				Note: "climb to the guard's tile"},
			{Num: 5, Input: Right, Pulse: true, Until: UntilBattleStart, Timeout: 1800,
				ExpectX: 0x1E, ExpectY: 0x25,
				Note: "push east into the guard while clearing dialogue; starts the fifth battle"},
			{Num: 6, Pulse: true, Until: UntilBattleEnd, Timeout: 7200,
				Note: "fight the fifth scripted battle to victory"},
			{Num: 7, Pulse: true, Until: UntilElapsed, Duration: 900, Timeout: 1200,
				ExpectX: 0x1E, ExpectY: 0x25,
				Note: "dismiss the reward windows; control returns here"},
			{Num: 8, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x21, Timeout: 900,
				ExpectX: 0x1E, ExpectY: 0x21,
				Note: "climb to the first zigzag ledge"},
			{Num: 9, Input: Left, Until: UntilPosition, Axis: AxisX,
				Compare: AtMost, Target: 0x1D, Timeout: 900,
				ExpectX: 0x1D, ExpectY: 0x21,
				Note: "zigzag sidestep west"},
			{Num: 10, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x1E, Timeout: 900,
				ExpectX: 0x1D, ExpectY: 0x1E,
				Note: "climb"},
			{Num: 11, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x1E, Timeout: 900,
				ExpectX: 0x1E, ExpectY: 0x1E,
				Note: "zigzag sidestep east"},
			{Num: 12, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x18, Timeout: 900,
				ExpectX: 0x1E, ExpectY: 0x18,
				Note: "climb"},
			{Num: 13, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x1F, Timeout: 900,
				ExpectX: 0x1F, ExpectY: 0x18,
				Note: "zigzag sidestep east"},
			{Num: 14, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x16, Timeout: 900,
				ExpectX: 0x1F, ExpectY: 0x16,
				Note: "final climb; trips the mine-shaft dialogue"},
			{Num: 15, Pulse: true, Until: UntilElapsed, Duration: 600, Timeout: 900,
				ExpectX: 0x1F, ExpectY: 0x16,
				Note: "clear the shaft dialogue"},
			// The transition is detected by the position jump. EXP-0036 first
			// tried +$1EA5 reaching $0D and the leg completed in 0 frames with
			// the party still on the exterior; CONTRA-0002 later showed why —
			// that byte is an event-flag byte, not a map indicator, so it was
			// never going to mark the transition. X only reaches $26 on the far
			// side of it.
			{Num: 16, Input: Up, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x26, Timeout: 1800,
				ExpectX: 0x26, ExpectY: 0x21,
				Note: "step into the shaft: map transition into the mines"},
			// The transition lands at Y=$21/$20; EXP-0035's recon walked on
			// to (26,1C), which is the documented milestone-05 coordinate.
			{Num: 17, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x1C, Timeout: 900,
				ExpectX: 0x26, ExpectY: 0x1C,
				Note: "walk to the documented mines-entry tile"},
		},
	}
}

// MinesEncounterRoute is the SCN-0001 segment-5 route: milestone 05
// onward into the mines toward the first random encounter (milestone
// 06). It is MinesRoute plus nine legs over the corridor EXP-0038's
// recon mapped: straight north from (26,1C) to (26,0B), east to
// (28,0B), north to (28,09) — stopping short of (2A,09), where recon
// tripped a scripted event (dialogue, party splits) that is out of
// this unit's scope. Legs 21-26 patrol back down the corridor and up
// again to extend the step budget without entering the event zone.
//
// The random encounter is expected to interrupt one of legs 18-26; the
// controller suspends walk legs while a battle owns the screen and
// restarts their timeout window afterwards, so an interrupted leg
// resumes and completes. Completing all legs without a sixth battle is
// the bounded negative result, not a route failure.
func MinesEncounterRoute() Route {
	base := MinesRoute()
	return Route{
		Name:          "SCN-0001 segment 5: mines corridor to the first random encounter",
		StartFrame:    base.StartFrame,
		NeutralFrames: base.NeutralFrames,
		Legs: append(base.Legs, []Leg{
			{Num: 18, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x0B, Timeout: 1200,
				ExpectX: 0x26, ExpectY: 0x0B,
				Note: "north corridor from the mines-entry tile; any post-transition fade is absorbed by the hold"},
			{Num: 19, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x28, Timeout: 900,
				ExpectX: 0x28, ExpectY: 0x0B,
				Note: "east along the trestle (up is blocked at 26,0B)"},
			{Num: 20, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x09, Timeout: 900,
				ExpectX: 0x28, ExpectY: 0x09,
				Note: "climb; stops one turn short of the (2A,09) event trigger"},
			{Num: 21, Input: Down, Until: UntilPosition, Axis: AxisY,
				Compare: AtLeast, Target: 0x0B, Timeout: 900,
				ExpectX: 0x28, ExpectY: 0x0B,
				Note: "patrol: back down the climb"},
			{Num: 22, Input: Left, Until: UntilPosition, Axis: AxisX,
				Compare: AtMost, Target: 0x26, Timeout: 900,
				ExpectX: 0x26, ExpectY: 0x0B,
				Note: "patrol: west to the corridor"},
			{Num: 23, Input: Down, Until: UntilPosition, Axis: AxisY,
				Compare: AtLeast, Target: 0x14, Timeout: 900,
				ExpectX: 0x26, ExpectY: 0x14,
				Note: "patrol: south down the corridor"},
			{Num: 24, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x0B, Timeout: 1200,
				ExpectX: 0x26, ExpectY: 0x0B,
				Note: "patrol: north again"},
			{Num: 25, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x28, Timeout: 900,
				ExpectX: 0x28, ExpectY: 0x0B,
				Note: "patrol: east again"},
			{Num: 26, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x09, Timeout: 900,
				ExpectX: 0x28, ExpectY: 0x09,
				Note: "patrol: final climb; route end = bounded no-encounter budget exhausted"},
		}...),
	}
}
