package rules

import (
	"testing"
)

func TestEvaluateTagRule(t *testing.T) {
	device := Device{
		Name: "router1",
		Tags: []string{"compliant", "core"},
	}

	rule := Rule{
		ID:           "tag-check",
		Type:         "tag",
		RequiredTags: []string{"compliant"},
	}

	result := evaluateTagRule(device, rule)
	if !result.Compliant {
		t.Errorf("Expected device to be compliant; got: %v", result)
	}
}

func TestEvaluateInterfaceRule(t *testing.T) {
	device := Device{
		Name:       "switch1",
		Interfaces: []string{"eth0", "ens3"},
	}

	rule := Rule{
		ID:    "iface-check",
		Type:  "interface",
		Regex: "^(eth\\d+|ens\\d+)$",
	}

	result := evaluateInterfaceRule(device, rule)
	if !result.Compliant {
		t.Errorf("Expected all interfaces to match; got: %v", result)
	}
}

func TestEvaluateIPRangeRule(t *testing.T) {
	device := Device{
		Name:      "host1",
		PrimaryIP: "10.10.5.4/24",
	}

	rule := Rule{
		ID:   "ip-range",
		Type: "ip_range",
		CIDR: "10.0.0.0/8",
	}

	result := evaluateIPRangeRule(device, rule)
	if !result.Compliant {
		t.Errorf("Expected IP to be within range; got: %v", result)
	}
}

func TestEvaluateAllowedValuesRule(t *testing.T) {
	device := Device{
		Name: "web1",
		Role: "web",
	}

	rule := Rule{
		ID:      "role-allowed",
		Type:    "allowed_values",
		Field:   "role",
		Allowed: []string{"web", "db", "cache"},
	}

	result := evaluateAllowedValuesRule(device, rule)
	if !result.Compliant {
		t.Errorf("Expected role to be allowed; got: %v", result)
	}
}

func TestEvaluateRequiredFieldsRule(t *testing.T) {
	device := Device{
		Name: "node1",
		Site: "dc1",
		Rack: "rack1",
	}

	rule := Rule{
		ID:     "required-fields",
		Type:   "required_fields",
		Fields: []string{"site", "rack"},
	}

	result := evaluateRequiredFieldsRule(device, rule)
	if !result.Compliant {
		t.Errorf("Expected required fields to be present; got: %v", result)
	}
}
