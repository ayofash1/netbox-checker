package rules

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Device represents a simplified NetBox device
type Device struct {
	Name       string
	Tags       []string
	PrimaryIP  string
	Role       string
	Site       string
	Rack       string
	Interfaces []string // interface names like eth0, ens33
}

// Result captures the outcome of applying a rule to a device
type Result struct {
	DeviceName string
	RuleID     string
	Compliant  bool
	Message    string
}

// Evaluate runs all rules against all devices
func Evaluate(devices []Device, ruleSet *RuleSet) []Result {
	var results []Result

	for _, device := range devices {
		for _, rule := range ruleSet.Rules {
			res := evaluateRule(device, rule)
			results = append(results, res)
		}
	}
	return results
}

// evaluateRule applies a single rule to a single device
func evaluateRule(device Device, rule Rule) Result {
	switch rule.Type {
	case "tag":
		return evaluateTagRule(device, rule)
	case "interface":
		return evaluateInterfaceRule(device, rule)
	case "ip_range":
		return evaluateIPRangeRule(device, rule)
	case "allowed_values":
		return evaluateAllowedValuesRule(device, rule)
	case "required_fields":
		return evaluateRequiredFieldsRule(device, rule)
	default:
		return Result{
			DeviceName: device.Name,
			RuleID:     rule.ID,
			Compliant:  false,
			Message:    fmt.Sprintf("Unknown rule type: %s", rule.Type),
		}
	}
}

func evaluateTagRule(device Device, rule Rule) Result {
	for _, required := range rule.RequiredTags {
		for _, tag := range device.Tags {
			if tag == required {
				return Result{device.Name, rule.ID, true, "Has required tag"}
			}
		}
	}
	return Result{device.Name, rule.ID, false, "Missing required tags"}
}

func evaluateInterfaceRule(device Device, rule Rule) Result {
	re := regexp.MustCompile(rule.Regex)
	for _, iface := range device.Interfaces {
		if !re.MatchString(iface) {
			return Result{device.Name, rule.ID, false, fmt.Sprintf("Interface '%s' does not match regex", iface)}
		}
	}
	return Result{device.Name, rule.ID, true, "All interfaces match regex"}
}

func evaluateIPRangeRule(device Device, rule Rule) Result {
	_, cidr, err := net.ParseCIDR(rule.CIDR)
	if err != nil {
		return Result{device.Name, rule.ID, false, "Invalid CIDR in rule"}
	}

	ip := net.ParseIP(strings.Split(device.PrimaryIP, "/")[0])
	if ip == nil {
		return Result{device.Name, rule.ID, false, "Invalid IP format"}
	}

	if cidr.Contains(ip) {
		return Result{device.Name, rule.ID, true, "IP is within CIDR"}
	}
	return Result{device.Name, rule.ID, false, "IP not in allowed CIDR"}
}

func evaluateAllowedValuesRule(device Device, rule Rule) Result {
	for _, allowed := range rule.Allowed {
		if deviceFieldValue(device, rule.Field) == allowed {
			return Result{device.Name, rule.ID, true, "Field value is allowed"}
		}
	}
	return Result{device.Name, rule.ID, false, fmt.Sprintf("Field '%s' value not allowed", rule.Field)}
}

func evaluateRequiredFieldsRule(device Device, rule Rule) Result {
	for _, field := range rule.Fields {
		if strings.TrimSpace(deviceFieldValue(device, field)) == "" {
			return Result{device.Name, rule.ID, false, fmt.Sprintf("Field '%s' is empty", field)}
		}
	}
	return Result{device.Name, rule.ID, true, "All required fields present"}
}

func deviceFieldValue(d Device, field string) string {
	switch field {
	case "role":
		return d.Role
	case "site":
		return d.Site
	case "rack":
		return d.Rack
	case "primary_ip4", "primary_ip":
		return d.PrimaryIP
	default:
		return ""
	}
}
