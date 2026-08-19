package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/portfolio/pf-developer-portal/internal/specbreak"
)

func main() {
	failOn := flag.String("fail-on", "ERR", "fail when findings at this level exist (only ERR is implemented)")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: oasdiff-gate [-fail-on ERR] base.yaml revision.yaml")
		os.Exit(2)
	}
	if *failOn != "ERR" {
		fmt.Fprintln(os.Stderr, "only -fail-on ERR is supported")
		os.Exit(2)
	}
	base, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	rev, err := os.ReadFile(flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	rep, err := specbreak.CompareYAML(base, rev)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	for _, e := range rep.Errors {
		fmt.Printf("%s %s %s: %s\n", e.ID, e.Method, e.Path, e.Message)
	}
	if rep.HasERR() {
		os.Exit(1)
	}
}
