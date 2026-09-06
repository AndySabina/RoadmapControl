package roadmap

import (
	"strings"
	"testing"
)

func TestPolicyRoadmapCanonicalJSON(t *testing.T) {
	policy, err := ParsePolicyRoadmap(
		Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}},
		[]byte("kind: policy\nadditional_tracker_types: [initiative]\n"),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := policy.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"modules":[{"kind":"policy","path":"policy.yaml","policy":{"additional_tracker_types":["initiative"]}}],"schema":"https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"}`
	if string(got) != want {
		t.Fatalf("CanonicalJSON() = %q, want %q", got, want)
	}
}

func TestPolicyRoadmapCanonicalJSONPreservesEmptyTrackerTypes(t *testing.T) {
	policy, err := ParsePolicyRoadmap(
		Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}},
		[]byte("kind: policy\nadditional_tracker_types: []\n"),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := policy.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"modules":[{"kind":"policy","path":"policy.yaml","policy":{"additional_tracker_types":[]}}],"schema":"https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"}`
	if string(got) != want {
		t.Fatalf("CanonicalJSON() = %q, want %q", got, want)
	}
}

func TestParsePolicyRoadmapRejectsSurrogateEscapes(t *testing.T) {
	manifest := Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}}
	for _, tt := range []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{"high surrogate", "kind: policy\nadditional_tracker_types: [\"\\uD800\"]\n", true},
		{"high surrogate followed by ASCII", "kind: policy\nadditional_tracker_types: [\"\\uD800\\u0041\"]\n", true},
		{"low followed by high surrogate", "kind: policy\nadditional_tracker_types: [\"\\uDC00\\uD800\"]\n", true},
		{"literal replacement character", "kind: policy\nadditional_tracker_types: [\"�\"]\n", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			policy, err := ParsePolicyRoadmap(manifest, []byte(tt.yaml))
			if tt.wantErr {
				if err == nil {
					t.Fatal("ParsePolicyRoadmap() unexpectedly accepted surrogate escape")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePolicyRoadmap() error = %v", err)
			}
			got, err := policy.CanonicalJSON()
			if err != nil {
				t.Fatal(err)
			}
			const want = `{"modules":[{"kind":"policy","path":"policy.yaml","policy":{"additional_tracker_types":["�"]}}],"schema":"https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"}`
			if string(got) != want {
				t.Fatalf("CanonicalJSON() = %q, want %q", got, want)
			}
		})
	}
}

func TestPolicyRoadmapCanonicalJSONPreservesValuesAndIsolation(t *testing.T) {
	policy := PolicyRoadmap{
		schema:     RoadmapSchemaURI,
		modulePath: "policy.yaml",
		additionalTrackerTypes: []string{
			"", "second", "first", "CJK 漢字 😀 � <>& \\\" / \\ \n\x01 é é",
		},
	}

	got, err := policy.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"modules":[{"kind":"policy","path":"policy.yaml","policy":{"additional_tracker_types":["","second","first","CJK 漢字 😀 � <>& \\\" / \\ \n\u0001 é é"]}}],"schema":"https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"}`
	if string(got) != want {
		t.Fatalf("CanonicalJSON() = %q, want %q", got, want)
	}
	got[0] = '!'
	again, err := policy.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != want {
		t.Fatalf("CanonicalJSON() after output mutation = %q, want %q", again, want)
	}
}

func TestPolicyRoadmapCanonicalJSONRejectsInvalidState(t *testing.T) {
	valid := PolicyRoadmap{schema: RoadmapSchemaURI, modulePath: "policy.yaml", additionalTrackerTypes: []string{}}
	for _, tt := range []struct {
		name   string
		policy PolicyRoadmap
	}{
		{"zero aggregate", PolicyRoadmap{}},
		{"wrong schema", PolicyRoadmap{schema: "urn:other", modulePath: "policy.yaml", additionalTrackerTypes: []string{}}},
		{"wrong module path", PolicyRoadmap{schema: RoadmapSchemaURI, modulePath: "other.yaml", additionalTrackerTypes: []string{}}},
		{"nil tracker types", PolicyRoadmap{schema: RoadmapSchemaURI, modulePath: "policy.yaml"}},
		{"invalid schema UTF-8", PolicyRoadmap{schema: string([]byte{0xff}), modulePath: "policy.yaml", additionalTrackerTypes: []string{}}},
		{"invalid path UTF-8", PolicyRoadmap{schema: RoadmapSchemaURI, modulePath: string([]byte{0xff}), additionalTrackerTypes: []string{}}},
		{"invalid tracker type UTF-8", PolicyRoadmap{schema: RoadmapSchemaURI, modulePath: "policy.yaml", additionalTrackerTypes: []string{string([]byte{0xed, 0xa0, 0x80})}}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.policy.CanonicalJSON()
			if err == nil || !strings.Contains(err.Error(), "invalid") {
				t.Fatalf("CanonicalJSON() error = %v, want invalid state", err)
			}
		})
	}
	if _, err := valid.CanonicalJSON(); err != nil {
		t.Fatalf("valid CanonicalJSON() error = %v", err)
	}
}

func TestValidateCanonicalStringsRejectsEveryStringPosition(t *testing.T) {
	valid := canonicalPolicyRoadmap{
		Modules: []canonicalPolicyModule{{
			Kind:   "policy",
			Path:   "policy.yaml",
			Policy: canonicalPolicy{AdditionalTrackerTypes: []string{"initiative"}},
		}},
		Schema: RoadmapSchemaURI,
	}
	for _, tt := range []struct {
		name  string
		value canonicalPolicyRoadmap
	}{
		{"schema", canonicalPolicyRoadmap{Modules: valid.Modules, Schema: string([]byte{0xff})}},
		{"path", canonicalPolicyRoadmap{Modules: []canonicalPolicyModule{{Kind: "policy", Path: string([]byte{0xff}), Policy: valid.Modules[0].Policy}}, Schema: RoadmapSchemaURI}},
		{"kind", canonicalPolicyRoadmap{Modules: []canonicalPolicyModule{{Kind: string([]byte{0xff}), Path: "policy.yaml", Policy: valid.Modules[0].Policy}}, Schema: RoadmapSchemaURI}},
		{"tracker type", canonicalPolicyRoadmap{Modules: []canonicalPolicyModule{{Kind: "policy", Path: "policy.yaml", Policy: canonicalPolicy{AdditionalTrackerTypes: []string{string([]byte{0xff})}}}}, Schema: RoadmapSchemaURI}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCanonicalStrings(tt.value); err == nil {
				t.Fatal("validateCanonicalStrings() accepted invalid UTF-8")
			}
		})
	}
	if err := validateCanonicalStrings(valid); err != nil {
		t.Fatalf("validateCanonicalStrings(valid) = %v", err)
	}
}
