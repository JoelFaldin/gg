package model

type Address struct {
	IPV4 string `yaml:"ipv4"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
