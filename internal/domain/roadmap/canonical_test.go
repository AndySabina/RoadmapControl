package roadmap

import "testing"

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
