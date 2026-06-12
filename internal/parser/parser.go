package parser

import (
	"encoding/binary"
	"fmt"
	"gg/internal/config"
	"gg/internal/records"
	"gg/internal/store"
	"net"
	"strings"
	"time"
)

type Header struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type Question struct {
	QName       string
	QType       uint16
	QClass      uint16
	QuestionEnd int
}

func ParseMessage(data []byte, store *store.Store, connection *net.UDPConn, addr *net.UDPAddr) ([]byte, error) {
	q := parseBody(data[12:])
	fmt.Println("type:", q.QType)

	ip, exists := store.SearchDomain(q.QName)

	// If domain does not exists on YAML:
	if !exists {
		return forwardAndCache(data, q, store)
	}

	// Prepare for response:
	var answerBytes []byte
	var err error

	switch q.QType {
	case 1:
		if ip.IPV4 == "" {
			return forwardAndCache(data, q, store)
		}
		answerBytes, err = records.RecordA{}.Serialize(ip.IPV4)
	case 28:
		if ip.IPV6 == "" {
			return forwardAndCache(data, q, store)
		}
		answerBytes, err = records.RecordAAAA{}.Serialize(ip.IPV6)
	default:
		return forwardAndCache(data, q, store)
	}

	if err != nil {
		fmt.Println("error while serializing", err)
	}

	// Prepare response header:
	headerBytes := make([]byte, 12)
	copy(headerBytes, data[0:12])

	headerBytes[2] |= 0x80 // Set QRbit to 1, "response"
	headerBytes[3] |= 0x80 // Set RAbit to 1

	// Counters:
	headerBytes[6] = 0x00
	headerBytes[7] = 0x01 // Set ANCOUNT to 1

	headerBytes[8] = 0x00
	headerBytes[9] = 0x00 // Set NSCOUNT to 0

	headerBytes[10] = 0x00
	headerBytes[11] = 0x00 // Set ARCOUNT to 0

	// Extract Question bytes:
	questionEnd := 12 + q.QuestionEnd + 4
	questionBytes := data[12:questionEnd]

	response := []byte{}
	response = append(response, headerBytes...)
	response = append(response, questionBytes...)
	response = append(response, answerBytes...)

	return response, nil

	// fmt.Printf("client searches %s -> IP in yaml: %s\n", q.QName, ip)
}

func parseHeader(header []byte) Header {
	return Header{
		ID:      binary.BigEndian.Uint16(header[0:2]),
		Flags:   binary.BigEndian.Uint16(header[2:4]),
		QDCount: binary.BigEndian.Uint16(header[4:6]),
		ANCount: binary.BigEndian.Uint16(header[6:8]),
		NSCount: binary.BigEndian.Uint16(header[8:10]),
		ARCount: binary.BigEndian.Uint16(header[10:12]),
	}
}

func parseBody(body []byte) Question {
	var str []string
	pointer := 0
	for {
		length := int(body[pointer])

		if length == 0 {
			pointer++
			break
		}

		start := pointer + 1
		end := start + length
		piece := string(body[start:end])

		str = append(str, piece)
		pointer = end
	}

	domain := strings.Join(str, ".")

	qtype := []byte{body[pointer], body[pointer+1]}
	qclass := []byte{body[pointer+2], body[pointer+3]}

	return Question{
		QName:       domain,
		QType:       binary.BigEndian.Uint16(qtype),
		QClass:      binary.BigEndian.Uint16(qclass),
		QuestionEnd: pointer,
	}
}

func extractIp(res []byte, q Question) (string, error) {
	answerStart := 12 + q.QuestionEnd + 4

	if len(res) <= answerStart+12 {
		return "", fmt.Errorf("no answer section\n")
	}

	rdLenPos := answerStart + 10
	rdLen := binary.BigEndian.Uint16(res[rdLenPos : rdLenPos+2])

	rDataPos := rdLenPos + 2
	if len(res) < rDataPos+int(rdLen) {
		return "", fmt.Errorf("malformatted or corrupted package\n")
	}

	ipBytes := res[rDataPos : rDataPos+int(rdLen)]

	ipNet := net.IP(ipBytes)
	return ipNet.String(), nil
}

func forwardAndCache(data []byte, question Question, store *store.Store) ([]byte, error) {
	res, err := forwardQuery(data)
	if err != nil {
		return nil, err
	}

	extractedIp, err := extractIp(res, question)
	if err == nil && extractedIp != "" {
		go store.WriteToYaml(question.QName, extractedIp, question.QType)
	}

	return res, nil
}

func forwardQuery(data []byte) ([]byte, error) {
	address := config.GetAddress()

	upstream_addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("error resolving upstream: %w\n", err)
	}

	local_conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: nil, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("error creating local socket for forward: %w\n", err)
	}
	defer local_conn.Close()

	_, err = local_conn.WriteTo(data, upstream_addr)
	if err != nil {
		return nil, fmt.Errorf("error writting to upstream: %w\n", err)
	}

	buf := make([]byte, 4096)

	local_conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, _, err := local_conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading upstream response: %w\n", err)
	}

	return buf[:n], nil
}
