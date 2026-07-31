package project

const Version = "0.1.0-dev"

func Help() string {
	return `ff6lab — FF6 reconstruction laboratory

Usage:
  ff6lab help
  ff6lab version
  ff6lab audit                          repository integrity checks
  ff6lab indexes generate               regenerate indexes/EXPERIMENTS.md
  ff6lab attackdata scan <hexdump>      decode a $C46AC0 table dump
                                        (bridge "read cpu" hex format)

Planned command groups:
  rom, evidence, asset, graphics, audio, validate, report

ROM-derived files remain local and are never embedded in the executable.
`
}
