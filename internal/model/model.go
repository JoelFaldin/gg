package model

type Address struct {
	Value string `yaml:"value"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}
