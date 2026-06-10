package model

type Address struct {
	IPV4 string `yaml:"ipv4,omitempty"`
	IPV6 string `yaml:"ipv6,omitempty"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
