package roadmap

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/gowebpki/jcs"
)

// CanonicalJSON returns the RFC 8785 canonical JSON representation of p.
func (p PolicyRoadmap) CanonicalJSON() ([]byte, error) {
	if err := p.validateCanonicalState(); err != nil {
		return nil, err
	}

	aggregate := canonicalPolicyRoadmap{
		Modules: []canonicalPolicyModule{{
			Kind:   "policy",
			Path:   p.modulePath,
			Policy: canonicalPolicy{AdditionalTrackerTypes: append(make([]string, 0, len(p.additionalTrackerTypes)), p.additionalTrackerTypes...)},
		}},
		Schema: p.schema,
	}
	if err := validateCanonicalStrings(aggregate); err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(aggregate)
	if err != nil {
		return nil, fmt.Errorf("marshal canonical policy roadmap: %w", err)
	}
	canonical, err := jcs.Transform(encoded)
	if err != nil {
		return nil, fmt.Errorf("canonicalize policy roadmap: %w", err)
	}
	return append([]byte(nil), canonical...), nil
}

type canonicalPolicyRoadmap struct {
	Modules []canonicalPolicyModule `json:"modules"`
	Schema  string                  `json:"schema"`
}

type canonicalPolicyModule struct {
	Kind   string          `json:"kind"`
	Path   string          `json:"path"`
	Policy canonicalPolicy `json:"policy"`
}

type canonicalPolicy struct {
	AdditionalTrackerTypes []string `json:"additional_tracker_types"`
}

func (p PolicyRoadmap) validateCanonicalState() error {
	if p.schema != RoadmapSchemaURI || p.modulePath != "policy.yaml" || p.additionalTrackerTypes == nil {
		return fmt.Errorf("invalid policy roadmap state")
	}
	return nil
}

func validateCanonicalStrings(aggregate canonicalPolicyRoadmap) error {
	if !utf8.ValidString(aggregate.Schema) {
		return fmt.Errorf("policy roadmap contains invalid UTF-8")
	}
	for _, module := range aggregate.Modules {
		if !utf8.ValidString(module.Kind) || !utf8.ValidString(module.Path) {
			return fmt.Errorf("policy roadmap contains invalid UTF-8")
		}
		for _, trackerType := range module.Policy.AdditionalTrackerTypes {
			if !utf8.ValidString(trackerType) {
				return fmt.Errorf("policy roadmap contains invalid UTF-8")
			}
		}
	}
	return nil
}
