package roadmap

import "testing"

func TestParsePolicyRoadmapAcceptsExactManifestAndPolicy(t *testing.T) {
	manifest := Manifest{
		Schema:  RoadmapSchemaURI,
		Modules: []string{"policy.yaml"},
	}

	policy, err := ParsePolicyRoadmap(manifest, []byte("kind: policy\nadditional_tracker_types: [initiative]\n"))
	if err != nil {
		t.Fatal(err)
	}
	if policy.Schema() != RoadmapSchemaURI {
		t.Fatalf("Schema() = %q", policy.Schema())
	}
	if policy.ModulePath() != "policy.yaml" {
		t.Fatalf("ModulePath() = %q", policy.ModulePath())
	}
	if got := policy.AdditionalTrackerTypes(); len(got) != 1 || got[0] != "initiative" {
		t.Fatalf("AdditionalTrackerTypes() = %#v", got)
	}
}

func TestParsePolicyRoadmapRejectsInvalidContracts(t *testing.T) {
	validManifest := Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}}
	for _, tt := range []struct {
		name     string
		manifest Manifest
		policy   string
	}{
		{"wrong schema", Manifest{Schema: "urn:other", Modules: []string{"policy.yaml"}}, "kind: policy\nadditional_tracker_types: []\n"},
		{"wrong path", Manifest{Schema: RoadmapSchemaURI, Modules: []string{"other.yaml"}}, "kind: policy\nadditional_tracker_types: []\n"},
		{"multiple modules", Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml", "trackers.yaml"}}, "kind: policy\nadditional_tracker_types: []\n"},
		{"top-level sequence", validManifest, "[]\n"},
		{"missing kind", validManifest, "additional_tracker_types: []\n"},
		{"missing additional tracker types", validManifest, "kind: policy\n"},
		{"null kind", validManifest, "kind: null\nadditional_tracker_types: []\n"},
		{"mapping kind", validManifest, "kind: {}\nadditional_tracker_types: []\n"},
		{"sequence kind", validManifest, "kind: [policy]\nadditional_tracker_types: []\n"},
		{"non-string scalar kind", validManifest, "kind: true\nadditional_tracker_types: []\n"},
		{"wrong kind", validManifest, "kind: feature\nadditional_tracker_types: []\n"},
		{"null types", validManifest, "kind: policy\nadditional_tracker_types: null\n"},
		{"scalar types", validManifest, "kind: policy\nadditional_tracker_types: feature\n"},
		{"mapping types", validManifest, "kind: policy\nadditional_tracker_types: {}\n"},
		{"wrong item type", validManifest, "kind: policy\nadditional_tracker_types: [1]\n"},
		{"unknown field", validManifest, "kind: policy\nadditional_tracker_types: []\nextra: no\n"},
		{"duplicate field", validManifest, "kind: policy\nkind: policy\nadditional_tracker_types: []\n"},
		{"alias", validManifest, "kind: policy\nadditional_tracker_types: [&types a, *types]\n"},
		{"multiple documents", validManifest, "kind: policy\nadditional_tracker_types: []\n---\nkind: policy\nadditional_tracker_types: []\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParsePolicyRoadmap(tt.manifest, []byte(tt.policy)); err == nil {
				t.Fatal("accepted invalid policy contract")
			}
		})
	}
}

func TestParsePolicyRoadmapPreservesAllowedValuesAndIsolation(t *testing.T) {
	manifest := Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}}
	policy, err := ParsePolicyRoadmap(manifest, []byte("kind: policy\nadditional_tracker_types: ['', feature, feature]\n"))
	if err != nil {
		t.Fatal(err)
	}
	manifest.Schema = "urn:other"
	manifest.Modules[0] = "other.yaml"

	got := policy.AdditionalTrackerTypes()
	if len(got) != 3 || got[0] != "" || got[1] != "feature" || got[2] != "feature" {
		t.Fatalf("AdditionalTrackerTypes() = %#v", got)
	}
	got[1] = "mutated"
	if policy.Schema() != RoadmapSchemaURI || policy.ModulePath() != "policy.yaml" {
		t.Fatalf("policy changed after manifest mutation: schema=%q path=%q", policy.Schema(), policy.ModulePath())
	}
	if types := policy.AdditionalTrackerTypes(); types[1] != "feature" {
		t.Fatalf("policy changed after output mutation: %#v", types)
	}

	empty, err := ParsePolicyRoadmap(Manifest{Schema: RoadmapSchemaURI, Modules: []string{"policy.yaml"}}, []byte("kind: policy\nadditional_tracker_types: []\n"))
	if err != nil {
		t.Fatal(err)
	}
	if types := empty.AdditionalTrackerTypes(); types == nil || len(types) != 0 {
		t.Fatalf("empty types = %#v", types)
	}
}
