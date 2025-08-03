package netbox

type DevicePage struct {
	Count   int         `json:"count"`
	Next    *string     `json:"next"`
	Results []DeviceRaw `json:"results"`
}

type DeviceRaw struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	PrimaryIP *IPAddress `json:"primary_ip4"`
	Role      *Nested    `json:"device_role"`
	Site      *Nested    `json:"site"`
	Rack      *Nested    `json:"rack"`
	Tags      []Tag      `json:"tags"`
}

type IPAddress struct {
	Address string `json:"address"`
}

type Nested struct {
	Name string `json:"name"`
}

type Tag struct {
	Name string `json:"name"`
}

type InterfacePage struct {
	Count   int            `json:"count"`
	Results []InterfaceRaw `json:"results"`
}

type InterfaceRaw struct {
	Name string `json:"name"`
}
