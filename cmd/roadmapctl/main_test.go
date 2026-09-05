package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelpDescribesBootstrapBoundary(t *testing.T) {
	var stdout bytes.Buffer

	if err := run([]string{"--help"}, &stdout); err != nil {
		t.Fatalf("run --help: %v", err)
	}

	help := stdout.String()
	for _, want := range []string{"roadmapctl", "bootstrap", "no product commands"} {
		if !strings.Contains(help, want) {
			t.Errorf("help %q does not contain %q", help, want)
		}
	}
}

func TestRunRejectsProductCommands(t *testing.T) {
	var stdout bytes.Buffer

	err := run([]string{"validate"}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "no product commands") {
		t.Fatalf("run validate error = %v, want unavailable product command", err)
	}
	if stdout.Len() != 0 {
		t.Fatalf("run validate output = %q, want none", stdout.String())
	}
}
