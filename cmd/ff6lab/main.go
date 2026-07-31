package main

import (
	"fmt"
	"os"

	"github.com/bensabler/ff6-decompile/internal/audit"
	"github.com/bensabler/ff6-decompile/internal/project"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ff6lab:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Print(project.Help())
		return nil
	}
	switch args[0] {
	case "version":
		fmt.Println(project.Version)
		return nil
	case "audit":
		findings, err := audit.Run(".")
		if err != nil {
			return err
		}
		for _, f := range findings {
			fmt.Println(f)
		}
		if len(findings) > 0 {
			return fmt.Errorf("audit: %d finding(s)", len(findings))
		}
		fmt.Println("audit: clean")
		return nil
	case "indexes":
		if len(args) > 1 && args[1] == "generate" {
			out, err := audit.GenerateExperimentIndex(".")
			if err != nil {
				return err
			}
			if err := os.WriteFile("indexes/EXPERIMENTS.md", []byte(out), 0o644); err != nil {
				return err
			}
			fmt.Println("indexes/EXPERIMENTS.md regenerated")
			return nil
		}
		return fmt.Errorf("usage: ff6lab indexes generate")
	case "help", "-h", "--help":
		fmt.Print(project.Help())
		return nil
	default:
		return fmt.Errorf("unknown command %q; run ff6lab help", args[0])
	}
}
