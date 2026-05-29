package parser

import (
	"encoding/binary"
	"fmt"
	"gg/internal/model"
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
	QLen    uint16
	QName   string
	QExtLen uint16
	QExt    string
	QType   uint16
	QClass  uint16
}

func ParseMessage(data []byte, store *model.Store) (string, error) {
	parseHeader(data[:12])
	parseBody(data[12:])

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

func parseBody(body []byte) {
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

	fmt.Println(domain)
	fmt.Println(qtype)
	fmt.Println(qclass)
}
