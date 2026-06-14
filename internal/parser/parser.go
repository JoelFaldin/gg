package parser

import (
	"encoding/binary"
	"fmt"
	"gg/internal/helpers"
	"gg/internal/logger"
	"gg/internal/model"
	"gg/internal/records"
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

func ParseMessage(data []byte, store *store.Store, connection *net.UDPConn, addr *net.UDPAddr) ([]byte, error) {
	customLogger := logger.CustomLogger()

	q := parseBody(data[12:])
	customLogger.Info(fmt.Sprintf("Type: %d", q.QType))

	ip, exists := store.SearchDomain(q.QName)

	// If domain does not exists on YAML:
	if !exists {
		return helpers.ForwardAndCache(data, q, store)
	}

	// Prepare for response:
	var answerBytes []byte
	var err error

	switch q.QType {
	case 1:
		if len(ip.IPv4) == 0 {
			return helpers.ForwardAndCache(data, q, store)
		}
		answerBytes, err = records.RecordA{}.Serialize(ip.IPv4)
	case 28:
		if len(ip.IPv6) == 0 {
			return helpers.ForwardAndCache(data, q, store)
		}
		answerBytes, err = records.RecordAAAA{}.Serialize(ip.IPv6)
	default:
		return helpers.ForwardAndCache(data, q, store)
	}

	if err != nil {
		customLogger.Error(fmt.Sprintf("Error while serializing: %v", err))
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

	customLogger.Debug("Record found on yaml!")

	return response, nil
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

func parseBody(body []byte) model.Question {
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

	return model.Question{
		QName:       domain,
		QType:       binary.BigEndian.Uint16(qtype),
		QClass:      binary.BigEndian.Uint16(qclass),
		QuestionEnd: pointer,
	}
}
