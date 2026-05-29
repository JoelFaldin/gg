package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gg/internal/store"
	"net"
	"strings"
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
	QName  string
	QType  uint16
	QClass uint16
}

func ParseMessage(data []byte, store *store.Store) (string, error) {
	parseHeader(data[:12])
	q := parseBody(data[12:])

	ip, err := store.BuscarIP(q.QName)
	if err != nil {
		return "", err
	}

	fmt.Printf("client searches %s -> IP in yaml: %s\n", q.QName, ip)

	return "", nil
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
	// Classic question section:
	// fmt.Println("----------")
	// fmt.Println("full body:", body)
	// fmt.Println("----------")

	var str []string
	pointer := 0
	for {
		length := int(body[pointer])

		if length == 0 {
			break
		}

		start := pointer + 1
		end := start + length
		piece := string(body[start:end])

		str = append(str, piece)
		pointer = end
	}

	domain := strings.Join(str, ".")

	qtype := []byte{body[pointer+1], body[pointer+2]}
	qclass := []byte{body[pointer+3], body[pointer+4]}

	// fmt.Println(domain)
	// fmt.Println(qtype)
	// fmt.Println(qclass)

	return Question{
		QName:  domain,
		QType:  binary.BigEndian.Uint16(qtype),
		QClass: binary.BigEndian.Uint16(qclass),
	}
}

func buildResponse(rawRequest []byte, questionEnd int, ipStr string) []byte {
	buf := new(bytes.Buffer)

	id := rawRequest[0:2]
	// Client id:
	buf.Write(id)

	// Its a message, dns server is working with no problems:
	binary.Write(buf, binary.BigEndian, uint16(0x8180))

	binary.Write(buf, binary.BigEndian, uint16(1)) // 1 question
	binary.Write(buf, binary.BigEndian, uint16(1)) // 1 answer
	binary.Write(buf, binary.BigEndian, uint16(0)) // NSCount
	binary.Write(buf, binary.BigEndian, uint16(0)) // ARCount

	// Copy Question section:
	buf.Write(rawRequest[12:questionEnd])

	// Build answer section:
	binary.Write(buf, binary.BigEndian, uint16(0xC00C))
	binary.Write(buf, binary.BigEndian, uint16(1))   // Type A
	binary.Write(buf, binary.BigEndian, uint16(1))   // Class IN
	binary.Write(buf, binary.BigEndian, uint16(300)) // TTL: 5 mins

	ip := net.ParseIP(ipStr).To4()

	binary.Write(buf, binary.BigEndian, uint16(len(ip)))
	buf.Write(ip)

	return buf.Bytes()
}
