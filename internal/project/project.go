package project

const Version = "0.1.0-dev"

func Help() string {
	return `ff6lab — FF6 reconstruction laboratory

Usage:
  ff6lab help
  ff6lab version

Planned command groups:
  project, rom, evidence, experiment, asset, graphics, audio, validate, report

ROM-derived files remain local and are never embedded in the executable.
`
}
