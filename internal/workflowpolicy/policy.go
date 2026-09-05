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
	permissionsIndent := -1
	for scanner.Scan() {
		rawLine := scanner.Text()
		indent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))
		line := strings.TrimSpace(strings.SplitN(rawLine, "#", 2)[0])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "permissions:") {
			permissionsIndent = indent
			permission := strings.TrimSpace(strings.TrimPrefix(line, "permissions:"))
			if permission == "write-all" {
				violations = append(violations, "write permission is forbidden")
			}
			if strings.HasPrefix(permission, "{") && strings.HasSuffix(permission, "}") {
				for _, entry := range strings.Split(strings.TrimSuffix(strings.TrimPrefix(permission, "{"), "}"), ",") {
					_, access, isPermission := strings.Cut(entry, ":")
					if isPermission && strings.TrimSpace(access) == "write" {
						violations = append(violations, "write permission is forbidden")
					}
				}
			}
			continue
		}
		if permissionsIndent < 0 || indent <= permissionsIndent {
			permissionsIndent = -1
			continue
		}
		_, access, isPermission := strings.Cut(line, ":")
		if isPermission && strings.TrimSpace(access) == "write" {
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
