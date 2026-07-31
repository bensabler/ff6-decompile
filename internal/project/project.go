package project

const Version = "0.1.0-dev"

func Help() string {
	return `ff6lab — FF6 reconstruction laboratory

Usage:
  ff6lab help
  ff6lab version
  ff6lab audit                          repository integrity checks
                                        (includes census validation)
  ff6lab indexes generate               regenerate all generated indexes
                                        and dashboards
  ff6lab census validate                content census + ROM ledger checks
  ff6lab census sync                    regenerate census-derived files,
                                        then validate
  ff6lab coverage summary               per-domain breadth counts
  ff6lab coverage gaps                  unlocated census entries
  ff6lab coverage domain <name>         one domain's entries + next actions
  ff6lab rom gaps                       unknown ROM spans, largest first
  ff6lab attackdata scan <hexdump>      decode a $C46AC0 table dump
                                        (bridge "read cpu" hex format)

Planned command groups:
  rom, evidence, asset, graphics, audio, validate, report

ROM-derived files remain local and are never embedded in the executable.
`
}
