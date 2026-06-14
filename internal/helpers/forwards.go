package helpers

import (
	"encoding/binary"
	"fmt"
	"gg/internal/config"
	"gg/internal/model"
	"gg/internal/store"
	"net"
	"net/netip"
	"time"
)

func extractAnswer(res []byte, q model.Question) (model.Answer, error) {
	ancount := binary.BigEndian.Uint16(res[6:8])

	answer := model.Answer{
		Name: q.QName,
		Type: q.QType,
		TTL:  60,
		IP:   []string{},
	}

	if ancount == 0 {
		return answer, nil
	}

	answerStart := 12 + q.QuestionEnd + 4
	offset := answerStart

	// Skip name section:
	for i := range int(ancount) {
		nextOffset, err := skipDomainName(res, offset)
		if err != nil {
			return answer, fmt.Errorf("Failed skipping name in answer %d: %w", i, err)
		}

		offset = nextOffset

		// Check for out of bonds errors:
		if offset+10 > len(res) {
			return answer, fmt.Errorf("Truncated data while reading record headers")
		}

		// Get type:
		ansType := binary.BigEndian.Uint16(res[offset : offset+2])
		offset += 2 // Move pointer past type
		offset += 2 // Skip class

		answer.Type = ansType

		// TTL:
		ttl := binary.BigEndian.Uint32(res[offset : offset+4])
		offset += 4 // Move pointer past ttl

		if i == 0 {
			answer.TTL = ttl
		}

		// RDLength:
		rdLength := int(binary.BigEndian.Uint16(res[offset : offset+2]))
		offset += 2 // Move pointer past rdLength

		rData := res[offset : offset+rdLength]

		var ipStr string
		switch ansType {
		case 1:
			if rdLength == 4 {
				addr := netip.AddrFrom4([4]byte(rData))
				ipStr = addr.String()
			}
		case 28:
			if rdLength == 16 {
				addr := netip.AddrFrom16([16]byte(rData))
				ipStr = addr.String()
			}
		default:
			return answer, fmt.Errorf("Skipping unsupported record type: %d", ansType)
		}

		if ipStr != "" {
			answer.IP = append(answer.IP, ipStr)
		}

		offset += rdLength
	}

	return answer, nil
}

func skipDomainName(data []byte, initialOffset int) (int, error) {
	offset := initialOffset
	if offset < len(data) {
		// Extra check to prevent out of bounds errors:
		if offset > len(data) {
			return 0, fmt.Errorf("Unexpected end of data while skipping name")
		}

		b := data[offset]

		if (b & 0xC0) == 0xC0 {
			return offset + 2, nil
		}
		if b == 0 {
			return offset + 1, nil
		}

		labelLength := int(b)

		offset += 1 + labelLength
	}

	return 0, fmt.Errorf("Domain name did not terminate properly")
}

func ForwardAndCache(data []byte, question model.Question, store *store.Store) ([]byte, error) {
	res, err := ForwardQuery(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to forward query: %w", err)
	}

	answer, err := extractAnswer(res, question)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract answer from remote: %w", err)
	}

	if len(answer.IP) == 0 && (question.QType == 1 || question.QType == 28) {
		go store.WriteToYaml(question.QName, answer.IP, question.QType)
	}

	return res, nil
}

func ForwardQuery(data []byte) ([]byte, error) {
	address := config.GetAddress()

	upstream_addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("error resolving upstream: %w", err)
	}

	local_conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("error creating local socket for forward: %w", err)
	}
	defer local_conn.Close()

	_, err = local_conn.WriteTo(data, upstream_addr)
	if err != nil {
		return nil, fmt.Errorf("error writting to upstream: %w", err)
	}

	buf := make([]byte, 4096)

	local_conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := local_conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading upstream response: %w", err)
	}

	return buf[:n], nil
}
