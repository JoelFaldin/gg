package model

type Address struct {
	IP string `yaml:"ip"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
