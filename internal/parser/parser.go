package parser

import (
	"encoding/binary"
	"fmt"
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
	QName   uint16
	QExtLen uint16
	QExt    uint16
	TTL     uint16
}

func ParseMessage(data []byte) (string, error) {
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
	fmt.Println("----------")
	fmt.Println("full body:", body)
	fmt.Println("----------")

	domain_len := body[0]
	domain := body[1 : domain_len+1]

	start_len := len(domain) + 1
	ext_len := body[len(domain)+1]

	extension := body[start_len+1 : start_len+int(ext_len)+1]

	fmt.Println("domain:", domain)
	fmt.Println("extension:", extension)
}
