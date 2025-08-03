package rules

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseRules(path string) (*RuleSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules RuleSet
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules yaml: %w", err)
	}

	return &rules, nil
}
