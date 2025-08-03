package rules

type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	ID           string   `yaml:"id"`
	Description  string   `yaml:"description"`
	Type         string   `yaml:"type"`
	Field        string   `yaml:"field,omitempty"`
	RequiredTags []string `yaml:"required_tags,omitempty"`
	Regex        string   `yaml:"regex,omitempty"`
	CIDR         string   `yaml:"cidr,omitempty"`
	Allowed      []string `yaml:"allowed,omitempty"`
	Fields       []string `yaml:"fields,omitempty"`
}
