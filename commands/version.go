package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const (
	ProtocolVersion int32  = 70015
	UserAgent       string = "/onecoin-btc:0.0.1/"
)

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP != nil {
				return ipnet.IP.To4().String()
			}

		}
	}
	return ""
}

func (m *VersionMessage) Timestamp(timeNow int) []byte {
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(timeNow))

	return timestamp
}

func (m *VersionMessage) AddrRecvIpAddress() []byte {
	buf := new(bytes.Buffer)
	ip := []byte(LocalIP())

	//ZEROs first
	//TODO: Try to shift bits instead
	for i := 0; i <= 16-len(ip)-1; i++ {
		err := binary.Write(buf, binary.BigEndian, []byte{0x00})
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}

	for _, v := range ip {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}

	return buf.Bytes()
}

func (m *VersionMessage) AddrRecvPort() []byte {
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(18334)) //TODO: take port from another place
	return port
}

type Peer struct {
	Connection net.Conn
}

func Connect(host string, port string) (*Peer, error) {
	connectionString := host + ":" + port
	log.Println("connect to", connectionString)
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		return nil, err
	}

	peer := &Peer{
		Connection: conn,
	}

	return peer, nil
}

//
const maxPayloadSize uint32 = 0x02000000 //32MiB

type Message struct {
	StartString []byte
	CommandName []byte
	PayloadSize []byte
	Checksum    []byte
}

func NewMessage(network string, controlMessage string) *Message {
	var startString []byte
	if network == "testnet" {
		startString = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	}

	commandName := make([]byte, 12)
	copy(commandName, controlMessage)

	// Fixed in this case that the payload is empty
	payloadSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadSize, uint32(0))

	// Fixed in this case that the payload is empty
	checksum := []byte{0x5d, 0xf6, 0xe0, 0xe2}

	return &Message{
		StartString: startString,
		CommandName: commandName,
		PayloadSize: payloadSize,
		Checksum:    checksum,
	}
}

type VersionMessage struct {
	Version  []byte
	Services []byte
}

func NewVersionMessage(fullNode bool) *VersionMessage {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(ProtocolVersion))

	services := make([]byte, 8)
	if fullNode {
		binary.LittleEndian.PutUint64(services, 0x01)
	} else {
		binary.LittleEndian.PutUint64(services, 0x00)
	}

	return &VersionMessage{
		Version:  version,
		Services: services,
	}
}

func main() {

}
