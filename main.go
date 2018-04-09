package main

import (
	"log"
	"net"

	"github.com/itschainkit/btc-protocol/commands"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:18333")
	if conn != nil {
		log.Println("Connected")
	}

	if err != nil {
		log.Fatal("error on connection to peer", err)
		panic(err)
	}

	version := commands.NewMessage("testnet", "version")
	byts := version.Bytes()
	log.Println("bytes to send", byts)

	_, err = conn.Write(version.Bytes())
	if err != nil {
		log.Println("error on send version message", err)
	}

	for {
	}
}
