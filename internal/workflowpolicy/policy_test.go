package workflowpolicy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPullRequestWorkflowIsUnprivilegedAndPinned(t *testing.T) {
	workflow := filepath.Join("..", "..", ".github", "workflows", "verify-pr.yml")
	contents, err := os.ReadFile(workflow)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}

	if violations := Validate(contents); len(violations) != 0 {
		t.Fatalf("workflow policy violations: %v", violations)
	}
}

func TestUnsafeWorkflowFixtureIsRejected(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("testdata", "unsafe.yml"))
	if err != nil {
		t.Fatalf("read unsafe fixture: %v", err)
	}

	violations := Validate(contents)
	if len(violations) != 5 {
		t.Fatalf("got %d violations (%v), want 5", len(violations), violations)
	}
}
