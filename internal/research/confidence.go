package research

import "fmt"

type Confidence string

const (
	Confirmed           Confidence = "confirmed"
	StrongHypothesis    Confidence = "strong-hypothesis"
	TentativeHypothesis Confidence = "tentative-hypothesis"
	Unknown             Confidence = "unknown"
)

func (c Confidence) Validate() error {
	switch c {
	case Confirmed, StrongHypothesis, TentativeHypothesis, Unknown:
		return nil
	default:
		return fmt.Errorf("invalid confidence %q", c)
	}
}
