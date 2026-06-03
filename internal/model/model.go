package model

type Address struct {
	Domain string `yaml:"domain"`
	IP     string `yaml:"ip"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
