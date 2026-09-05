// Package roadmap defines deterministic, filesystem-independent roadmap values.
package roadmap

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

const maxYAMLBytes = 1 << 20
const maxYAMLDepth = 64

type Manifest struct {
	Schema  string   `yaml:"schema"`
	Modules []string `yaml:"modules"`
}

// ParseManifest accepts the one typed manifest document used to declare modules.
func ParseManifest(data []byte) (Manifest, error) {
	var m Manifest
	n, err := document(data)
	if err != nil {
		return m, err
	}
	if n.Kind != yaml.MappingNode {
		return m, fmt.Errorf("manifest must be a mapping")
	}
	seen := map[string]bool{}
	for i := 0; i < len(n.Content); i += 2 {
		key, value := n.Content[i], n.Content[i+1]
		if key.Value != "schema" && key.Value != "modules" {
			return m, fmt.Errorf("unknown manifest field %q", key.Value)
		}
		if seen[key.Value] {
			return m, fmt.Errorf("duplicate manifest field %q", key.Value)
		}
		seen[key.Value] = true
		if key.Value == "schema" && (value.Kind != yaml.ScalarNode || value.Tag != "!!str") {
			return m, fmt.Errorf("schema must be a string")
		}
		if key.Value == "modules" {
			if value.Kind != yaml.SequenceNode {
				return m, fmt.Errorf("modules must be a sequence")
			}
			for _, p := range value.Content {
				if p.Kind != yaml.ScalarNode || p.Tag != "!!str" {
					return m, fmt.Errorf("module paths must be strings")
				}
			}
		}
	}
	if !seen["schema"] || !seen["modules"] {
		return m, fmt.Errorf("manifest requires schema and modules")
	}
	if err := n.Decode(&m); err != nil {
		return m, fmt.Errorf("decode manifest: %w", err)
	}
	if !validAbsoluteURI(m.Schema) {
		return m, fmt.Errorf("schema must be an absolute URI")
	}
	paths := map[string]bool{}
	for i, p := range m.Modules {
		if !SafeRelativePath(p) {
			return m, fmt.Errorf("unsafe module path %q", p)
		}
		if paths[p] {
			return m, fmt.Errorf("duplicate module path %q", p)
		}
		paths[p] = true
		m.Modules[i] = p
	}
	return m, nil
}

// SafeRelativePath accepts normalized, slash-separated YAML paths below a root.
func SafeRelativePath(p string) bool {
	return p != "" && !strings.ContainsRune(p, '\x00') && p == path.Clean(p) && !strings.Contains(p, "\\") && !strings.HasPrefix(p, "/") && p != "." && !strings.HasPrefix(p, "../") && !strings.EqualFold(path.Ext(p), "") && (strings.EqualFold(path.Ext(p), ".yaml") || strings.EqualFold(path.Ext(p), ".yml"))
}

func validAbsoluteURI(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil || !u.IsAbs() {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			if i+2 >= len(s) || !isHex(s[i+1]) || !isHex(s[i+2]) {
				return false
			}
			i += 2
			continue
		}
		if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'()*+,;=", rune(s[i])) {
			return false
		}
	}
	return true
}

func isHex(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'f' || 'A' <= b && b <= 'F'
}

// ValidateYAML applies document-level safety checks without interpreting a module.
func ValidateYAML(data []byte) error { _, err := document(data); return err }

func document(data []byte) (*yaml.Node, error) {
	if len(data) > maxYAMLBytes {
		return nil, fmt.Errorf("YAML exceeds %d-byte limit", maxYAMLBytes)
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("YAML is not UTF-8")
	}
	d := yaml.NewDecoder(bytes.NewReader(data))
	var n yaml.Node
	if err := d.Decode(&n); err != nil {
		return nil, fmt.Errorf("decode YAML: %w", err)
	}
	if n.Kind == 0 {
		return nil, fmt.Errorf("YAML document is empty")
	}
	var extra yaml.Node
	if err := d.Decode(&extra); err == nil {
		return nil, fmt.Errorf("YAML must contain exactly one document")
	} else if err != io.EOF {
		return nil, fmt.Errorf("trailing YAML: %w", err)
	}
	if len(n.Content) != 1 {
		return nil, fmt.Errorf("invalid YAML document")
	}
	if err := safeNode(n.Content[0], 0); err != nil {
		return nil, err
	}
	return n.Content[0], nil
}

func safeNode(n *yaml.Node, depth int) error {
	if depth > maxYAMLDepth {
		return fmt.Errorf("YAML exceeds depth limit")
	}
	if n.Kind == yaml.AliasNode {
		return fmt.Errorf("YAML aliases are not supported")
	}
	if n.Kind == yaml.MappingNode {
		seen := map[string]bool{}
		for i := 0; i < len(n.Content); i += 2 {
			if n.Content[i].Kind != yaml.ScalarNode || n.Content[i].Tag != "!!str" {
				return fmt.Errorf("YAML mapping keys must be strings")
			}
			if seen[n.Content[i].Value] {
				return fmt.Errorf("duplicate YAML key %q", n.Content[i].Value)
			}
			seen[n.Content[i].Value] = true
		}
	}
	for _, c := range n.Content {
		if err := safeNode(c, depth+1); err != nil {
			return err
		}
	}
	return nil
}
