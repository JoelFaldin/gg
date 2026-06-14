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
func (r RecordA) Serialize(ipStrs []string) ([]byte, error) {
	buf := make([]byte, 0, 16*len(ipStrs))

	for _, ipStr := range ipStrs {
		// Name:
		buf = append(buf, 0xc0, 0x0c)

		// Type:
		buf = append(buf, 0x00, 0x01) // 1, IPv4

		// Class:
		buf = append(buf, 0x00, 0x01) // 1, IN

		// TTL:
		ttl := make([]byte, 4)
		binary.BigEndian.PutUint32(ttl, 300)
		buf = append(buf, ttl...)

		// RDLength:
		buf = append(buf, 0x00, 0x04) // 4, IPv4, 4 bytes

		// RData:
		ip := net.ParseIP(ipStr).To4()
		if ip == nil {
			return nil, fmt.Errorf("invalid ipv4")
		}
		buf = append(buf, ip...)
	}

	return buf, nil
}

type RecordAAAA struct{}

func (r RecordAAAA) Type() uint16 { return 28 }
func (r RecordAAAA) Serialize(ipStrs []string) ([]byte, error) {
	buf := make([]byte, 0, 16*len(ipStrs))

	for _, ipStr := range ipStrs {
		// Name:
		buf = append(buf, 0xc0, 0x0c)

		// Type:
		buf = append(buf, 0x00, 0x1c)

		// Class:
		buf = append(buf, 0x00, 0x01)

		// TTL:
		ttl := make([]byte, 4)
		binary.BigEndian.PutUint32(ttl, 300)
		buf = append(buf, ttl...)

		// RDLength:
		buf = append(buf, 0x00, 0x10)

		// RData:
		ip := net.ParseIP(ipStr).To16()
		if ip == nil {
			return nil, fmt.Errorf("invalid ipv6")
		}
		buf = append(buf, ip...)
	}

	return buf, nil
}
