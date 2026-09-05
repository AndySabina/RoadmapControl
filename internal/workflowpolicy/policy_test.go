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

func TestWriteCapablePermissionsAreRejected(t *testing.T) {
	tests := []struct {
		name     string
		workflow string
	}{
		{
			name: "permission mapping with comment",
			workflow: "on:\n" +
				"  pull_request:\n" +
				"permissions:\n" +
				"  contents: write # comment\n",
		},
		{
			name: "write-all shorthand",
			workflow: "on:\n" +
				"  pull_request:\n" +
				"permissions: write-all\n",
		},
		{
			name: "inline permission mapping",
			workflow: "on:\n" +
				"  pull_request:\n" +
				"permissions: {contents: write}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violations := Validate([]byte(tt.workflow))
			for _, violation := range violations {
				if violation == "write permission is forbidden" {
					return
				}
			}
			t.Fatalf("violations = %v, want write permission violation", violations)
		})
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
