package records

import (
	"encoding/binary"
	"fmt"
	"net"
)

type DNSRecord interface {
	Type() uint16
	Serialize(value string) ([]byte, error)
}

type RecordA struct{}

func (r RecordA) Type() uint16 { return 1 }
func (r RecordA) Serialize(ipStr string) ([]byte, error) {
	buf := make([]byte, 0, 16)

	// Name:
	buf = append(buf, 0xc0, 0x0c)

	// Type:
	buf = append(buf, 0x00, 0x01) // 1, IPv4

	// Class:
	buf = append(buf, 0x00, 0x01) // 1, IN

	// TTL:
	ttl := make([]byte, 0, 4)
	binary.BigEndian.PutUint16(ttl, 300)
	buf = append(buf, ttl...)

	// RDLength:
	buf = append(buf, 0x00, 0x04) // 4, IPv4, 4 bytes

	// RData:
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return nil, fmt.Errorf("invalid ipv4\n")
	}
	buf = append(buf, ip...)

	return buf, nil
}
