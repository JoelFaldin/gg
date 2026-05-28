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

func ParseMessage(data []byte) (string, error) {
	header := Header{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		QDCount: binary.BigEndian.Uint16(data[4:6]),
		ANCount: binary.BigEndian.Uint16(data[6:8]),
		NSCount: binary.BigEndian.Uint16(data[8:10]),
		ARCount: binary.BigEndian.Uint16(data[10:12]),
	}

	fmt.Println("id", header.ID)
	fmt.Println("flags", header.Flags)
	fmt.Println("qdCount", header.QDCount)
	fmt.Println("anCount", header.ANCount)
	fmt.Println("nsCount", header.NSCount)
	fmt.Println("arCount", header.ARCount)

	return "", nil
}
