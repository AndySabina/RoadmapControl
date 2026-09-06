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

func TestParsePolicyRoadmapRejectsWrongManifestSchema(t *testing.T) {
	manifest := Manifest{
		Schema:  "urn:roadmapcontrol:schema:v1",
		Modules: []string{"policy.yaml"},
	}

	if _, err := ParsePolicyRoadmap(manifest, []byte("kind: policy\nadditional_tracker_types: []\n")); err == nil {
		t.Fatal("accepted a manifest without the RoadmapControl schema URI")
	}
}
