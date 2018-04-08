package commands

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"net"
	"time"
)

const (
	ProtocolVersion int32  = 70015
	UserAgent       string = "/onecoin-btc:0.0.1/"
)

func LocalIP() (string, []byte) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP != nil {
				return ipnet.IP.To4().String(), ipnet.IP.To4()
			}

		}
	}
	return "", nil
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

	var message ControlMessage
	if controlMessage == "version" {
		message = NewVersionMessage(true)
	}

	payloadSize := make([]byte, 4)
	checksum := make([]byte, 4)
	// Fixed in this case that the payload is empty
	if msgLength := message.Length(); msgLength < 1 {
		binary.LittleEndian.PutUint32(payloadSize, uint32(0))
		copy(checksum, []byte{0x5d, 0xf6, 0xe0, 0xe2})
	} else {
		binary.LittleEndian.PutUint32(payloadSize, uint32(msgLength))
		//checksum
		payloadSha256 := sha256.Sum256(message.Payload())
		bck := new(bytes.Buffer)
		for _, b := range payloadSha256 {
			bck.WriteByte(b)
		}
		shaOfSha := sha256.Sum256(bck.Bytes())
		copy(checksum, shaOfSha[0:4])
	}

	return &Message{
		StartString: startString,
		CommandName: commandName,
		PayloadSize: payloadSize,
		Checksum:    checksum,
	}
}

type ControlMessage interface {
	Payload() []byte
	Length() int
}

type VersionMessage struct {
	Version    []byte
	Services   []byte
	Timestamp  []byte
	FromIpPort []byte
	ToIpPort   []byte
	Nonce      []byte
	UserAgent  []byte
	LastBlock  []byte
	Relay      []byte
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

	timestamp := make([]byte, 8)
	timeNow := time.Now().Unix()
	binary.LittleEndian.PutUint64(timestamp, uint64(timeNow))

	//from := "127.0.0.1:8333"
	_, ip := LocalIP()
	fromIpPortBuffer := new(bytes.Buffer)
	paddingLeft := make([]byte, 12)
	binary.Write(fromIpPortBuffer, binary.BigEndian, paddingLeft)
	binary.Write(fromIpPortBuffer, binary.BigEndian, ip)
	binary.Write(fromIpPortBuffer, binary.BigEndian, uint16(8333))

	//to := "127.0.0.1:18333"
	toIpPortBuffer := new(bytes.Buffer)
	binary.Write(toIpPortBuffer, binary.BigEndian, paddingLeft)
	binary.Write(toIpPortBuffer, binary.BigEndian, ip)
	binary.Write(toIpPortBuffer, binary.BigEndian, uint16(18333))

	//nonce
	nonce := make([]byte, 8)

	//useragent
	userAgent := make([]byte, len(UserAgent))
	copy(userAgent, []byte(UserAgent))

	//last block
	lastBlock := make([]byte, 4)
	binary.LittleEndian.PutUint32(lastBlock, uint32(0))

	//relay
	relay := []byte{0x01}

	return &VersionMessage{
		Version:    version,
		Services:   services,
		Timestamp:  timestamp,
		FromIpPort: fromIpPortBuffer.Bytes(),
		ToIpPort:   toIpPortBuffer.Bytes(),
		Nonce:      nonce,
		UserAgent:  userAgent,
		LastBlock:  lastBlock,
		Relay:      relay,
	}
}

func (m *VersionMessage) Payload() []byte {
	payload := new(bytes.Buffer)
	payload.Write(m.Version)
	payload.Write(m.Services)
	payload.Write(m.Timestamp)
	payload.Write(m.FromIpPort)
	payload.Write(m.ToIpPort)
	payload.Write(m.Nonce)
	payload.Write(m.UserAgent)
	payload.Write(m.LastBlock)
	payload.Write(m.Relay)

	return payload.Bytes()
}

func (m *VersionMessage) Length() int {
	return len(m.Payload())
}

func main() {
}
