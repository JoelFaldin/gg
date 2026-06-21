package model

type Address struct {
	TTL  int      `yaml:"ttl,omitempty"`
	IPv4 []string `yaml:"ipv4,omitempty"`
	IPv6 []string `yaml:"ipv6,omitempty"`
}

type Data struct {
	Storage map[string]Address `yaml:"storage"`
}

type Question struct {
	QName       string
	QType       uint16
	QClass      uint16
	QuestionEnd int
}

type Answer struct {
	Name string
	Type uint16
	TTL  uint32
	IP   []string
}
