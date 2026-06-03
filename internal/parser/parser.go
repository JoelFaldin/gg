package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func ParseMessage(data []byte, store *store.Store, connection *net.UDPConn, addr *net.UDPAddr) {
	parseHeader(data[:12])
	q := parseBody(data[12:])

	ip, err := store.SearchDomain(q.QName)
	// Domain found:
	if err == nil {
		fmt.Println("--- domain found on yaml!!")
		res := buildResponse(data, q.QuestionEnd, ip)
		connection.WriteToUDP(res, addr)

	} else { // Domain not found:
		fmt.Println("--- domain not found, forwarding...")
		internet_res, addr, err := forwardQuery(data)
		if err != nil {
			fmt.Println("error on forward function", err)
			return
		}

		header := parseHeader(internet_res[0:12])
		pointer := 12

		body := parseBody(internet_res[pointer:])
		pointer += body.QuestionEnd

		// Add 4 bytes (Qtype (2) and Qclass (2))
		pointer += 4

		if header.ANCount > 0 {
			start_answer := pointer

			rdLenth := binary.BigEndian.Uint16(internet_res[start_answer+10 : start_answer+12])
			fin_answer := start_answer + 12 + int(rdLenth)
			// 12 is the number of bytes from the start to the field RDLength!!

			answer := internet_res[start_answer:fin_answer]
			ip_start := len(answer) - int(rdLenth)

			ip_final := answer[ip_start:]

			r := net.IP(ip_final).String()

			// Save to yaml:
			// Extract IP from internet_res:
			go store.WriteToYaml(q.QName, r)

			connection.WriteToUDP(internet_res, addr)
		}
	}

	fmt.Printf("client searches %s -> IP in yaml: %s\n", q.QName, ip)
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

	qtype := []byte{body[pointer+1], body[pointer+2]}
	qclass := []byte{body[pointer+3], body[pointer+4]}

	return Question{
		QName:       domain,
		QType:       binary.BigEndian.Uint16(qtype),
		QClass:      binary.BigEndian.Uint16(qclass),
		QuestionEnd: pointer,
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

func forwardQuery(data []byte) ([]byte, *net.UDPAddr, error) {
	upstream_addr, err := net.ResolveUDPAddr("udp", "1.1.1.1:53")
	if err != nil {
		return nil, nil, err
	}

	local_conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, nil, err
	}
	defer local_conn.Close()

	_, err = local_conn.WriteTo(data, upstream_addr)
	if err != nil {
		return nil, nil, err
	}

	buf := make([]byte, 1024)

	local_conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	n, addr, err := local_conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}

	return buf[:n], addr, nil
}
