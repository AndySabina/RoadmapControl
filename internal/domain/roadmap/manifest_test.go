package roadmap

import (
	"strings"
	"testing"
)

func TestParseManifestRejectsUnsafeInput(t *testing.T) {
	cases := []struct{ name, input string }{
		{"duplicate key", "schema: https://schemas.example/v1\nschema: https://schemas.example/v1\nmodules: [a.yaml]\n"},
		{"ambiguous schema scalar", "schema: 12\nmodules: [a.yaml]\n"},
		{"ambiguous path scalar", "schema: https://schemas.example/v1\nmodules: [12]\n"},
		{"unknown field", "schema: https://schemas.example/v1\nmodules: [a.yaml]\nextra: no\n"},
		{"trailing document", "schema: https://schemas.example/v1\nmodules: [a.yaml]\n---\nextra: no\n"},
		{"malformed trailing document", "schema: https://schemas.example/v1\nmodules: [a.yaml]\n---\n[\n"},
		{"alias expansion", "schema: https://schemas.example/v1\nmodules: [&a a.yaml, *a]\n"},
		{"unsafe path", "schema: https://schemas.example/v1\nmodules: [../a.yaml]\n"},
		{"NUL module path", "schema: urn:roadmapcontrol:schema:v1\nmodules: [\"a\\0.yaml\"]\n"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseManifest([]byte(tt.input)); err == nil {
				t.Fatal("accepted unsafe manifest")
			}
		})
	}
	if _, err := ParseManifest([]byte{'s', 0xff}); err == nil {
		t.Fatal("accepted non-UTF-8")
	}
}

func TestParseManifestValidatesAbsoluteSchemaURI(t *testing.T) {
	for _, tt := range []struct {
		name   string
		schema string
		valid  bool
	}{
		{"HTTPS URI", "https://schemas.example/v1", true},
		{"URN", "urn:roadmapcontrol:schema:v1", true},
		{"hostless file URI", "file:///schemas/roadmap/v1/schema.json", true},
		{"relative path", "schemas/roadmap/v1/schema.json", false},
		{"invalid URI", "://schemas.example/v1", false},
		{"malformed opaque URI escape", "urn:bad%zz", false},
		{"malformed opaque URI character", "urn:bad value", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseManifest([]byte("schema: " + tt.schema + "\nmodules: [a.yaml]\n"))
			if (err == nil) != tt.valid {
				t.Fatalf("ParseManifest(%q) error = %v, valid = %v", tt.schema, err, tt.valid)
			}
		})
	}
}

func TestValidateYAMLRejectsDuplicateAndAliasCycle(t *testing.T) {
	for _, input := range []string{
		"one: 1\none: 2\n",
		"items: &a [*a]\n",
		"0xB: first\n11: second\n",
		"true: first\nTRUE: second\n",
	} {
		if err := ValidateYAML([]byte(input)); err == nil {
			t.Fatal("accepted unsafe module YAML")
		}
	}
}

func TestParseManifestNormalizesOnlySafeYAMLPaths(t *testing.T) {
	got, err := ParseManifest([]byte("schema: https://schemas.example/v1\nmodules:\n  - nested/one.YAML\n  - two.yml\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Schema != "https://schemas.example/v1" || strings.Join(got.Modules, ",") != "nested/one.YAML,two.yml" {
		t.Fatalf("got %#v", got)
	}
}
