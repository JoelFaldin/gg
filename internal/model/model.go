package model

type Address struct {
	IPv4 string `yaml:"ipv4,omitempty"`
	IPv6 string `yaml:"ipv6,omitempty"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
