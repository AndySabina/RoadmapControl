package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const helpText = `roadmapctl — RoadmapControl bootstrap

Usage:
  roadmapctl --help

RoadmapControl is in bootstrap: this executable has no product commands.
`

func run(args []string, stdout io.Writer) error {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		_, err := fmt.Fprint(stdout, helpText)
		return err
	}
	return errors.New("no product commands are available; use --help")
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
