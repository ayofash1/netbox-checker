package netbox

import "github.com/ayofash1/netbox-checker/internal/rules"

func MapToRuleDevice(d DeviceRaw) rules.Device {
	tags := make([]string, len(d.Tags))
	for i, tag := range d.Tags {
		tags[i] = tag.Name
	}

	ip := ""
	if d.PrimaryIP != nil {
		ip = d.PrimaryIP.Address
	}

	role := ""
	if d.Role != nil {
		role = d.Role.Name
	}

	site := ""
	if d.Site != nil {
		site = d.Site.Name
	}

	rack := ""
	if d.Rack != nil {
		rack = d.Rack.Name
	}

	return rules.Device{
		Name:       d.Name,
		Tags:       tags,
		PrimaryIP:  ip,
		Role:       role,
		Site:       site,
		Rack:       rack,
		Interfaces: []string{}, // TODO: fetch interfaces later
	}
}
