package workflowpolicy

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

var actionReference = regexp.MustCompile(`(?m)^\s*-?\s*uses:\s*[^@\s]+@([^\s#]+)`)
var fullCommitSHA = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

// Validate returns policy violations for an unprivileged pull-request workflow.
func Validate(workflow []byte) []string {
	var violations []string
	text := string(workflow)
	if strings.Contains(text, "pull_request_target") {
		violations = append(violations, "pull_request_target is forbidden")
	}
	if strings.Contains(text, "secrets:") || strings.Contains(text, "secrets.") {
		violations = append(violations, "secrets are forbidden")
	}
	if !strings.Contains(text, "pull_request:") {
		violations = append(violations, "pull_request trigger is required")
	}

	scanner := bufio.NewScanner(bytes.NewReader(workflow))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasSuffix(line, ": write") {
			violations = append(violations, "write permission is forbidden")
		}
	}
	for _, match := range actionReference.FindAllStringSubmatch(text, -1) {
		if !fullCommitSHA.MatchString(match[1]) {
			violations = append(violations, "action reference must use a full commit SHA")
		}
	}
	return violations
}
