package roadmap

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// RoadmapSchemaURI is the only schema accepted by the policy module contract.
const RoadmapSchemaURI = "https://github.com/AndySabina/RoadmapControl/schemas/roadmap/v1/schema.json"

// PolicyRoadmap is the read-only policy module contract.
type PolicyRoadmap struct {
	schema                 string
	modulePath             string
	additionalTrackerTypes []string
}

// ParsePolicyRoadmap validates and reads the policy module declared by manifest.
func ParsePolicyRoadmap(manifest Manifest, policyYAML []byte) (PolicyRoadmap, error) {
	var policy PolicyRoadmap
	if manifest.Schema != RoadmapSchemaURI {
		return policy, fmt.Errorf("policy roadmap requires schema %q", RoadmapSchemaURI)
	}
	if len(manifest.Modules) != 1 || manifest.Modules[0] != "policy.yaml" {
		return policy, fmt.Errorf("policy roadmap requires exactly policy.yaml")
	}

	n, err := document(policyYAML)
	if err != nil {
		return policy, err
	}
	if n.Kind != yaml.MappingNode {
		return policy, fmt.Errorf("policy must be a mapping")
	}

	seen := map[string]bool{}
	for i := 0; i < len(n.Content); i += 2 {
		key, value := n.Content[i], n.Content[i+1]
		if key.Value != "kind" && key.Value != "additional_tracker_types" {
			return policy, fmt.Errorf("unknown policy field %q", key.Value)
		}
		if seen[key.Value] {
			return policy, fmt.Errorf("duplicate policy field %q", key.Value)
		}
		seen[key.Value] = true

		switch key.Value {
		case "kind":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" || value.Value != "policy" {
				return policy, fmt.Errorf("kind must be the string policy")
			}
		case "additional_tracker_types":
			if value.Kind != yaml.SequenceNode {
				return policy, fmt.Errorf("additional_tracker_types must be a sequence")
			}
			policy.additionalTrackerTypes = make([]string, len(value.Content))
			for i, item := range value.Content {
				if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
					return policy, fmt.Errorf("additional_tracker_types must contain strings")
				}
				policy.additionalTrackerTypes[i] = item.Value
			}
		}
	}
	if !seen["kind"] || !seen["additional_tracker_types"] {
		return PolicyRoadmap{}, fmt.Errorf("policy requires kind and additional_tracker_types")
	}

	policy.schema = manifest.Schema
	policy.modulePath = manifest.Modules[0]
	return policy, nil
}

// Schema returns the validated manifest schema URI.
func (p PolicyRoadmap) Schema() string { return p.schema }

// ModulePath returns the validated policy module path.
func (p PolicyRoadmap) ModulePath() string { return p.modulePath }

// AdditionalTrackerTypes returns a copy of the configured additional tracker types.
func (p PolicyRoadmap) AdditionalTrackerTypes() []string {
	types := make([]string, len(p.additionalTrackerTypes))
	copy(types, p.additionalTrackerTypes)
	return types
}
