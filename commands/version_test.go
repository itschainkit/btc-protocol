package commands

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {
	ip, _ := LocalIP()
	if ip == "" {
		t.Error("should not be nil", ip)
	}
}

func TestVersionAddrRecvPort(t *testing.T) {
	version := NewVersionMessage(true)
	port := version.AddrRecvPort()
	if hex := fmt.Sprintf("%x", port); hex != "479e" {
		t.Error("should return the right port number", hex)
	}

	if len(port) != 2 {
		t.Error("port number should have the right length", 2)
	}
}

func TestNewMessage(t *testing.T) {
	message := NewMessage("testnet", "version")
	if message.StartString == nil {
		t.Error("start string should not be nil")
	}

	if hex := fmt.Sprintf("%x", message.StartString); hex != "f9beb4d9" {
		t.Error("start string should have right value", hex)
	}

	if length := len(message.StartString); length != 4 {
		t.Error("start string should have the right length", length)
	}

	if message.CommandName == nil {
		t.Error("command name should not be nil")
	}

	if hex := fmt.Sprintf("%x", message.CommandName); hex != "76657273696f6e0000000000" {
		t.Error("command name should have the right value", hex)
	}

	if length := len(message.CommandName); length != 12 {
		t.Error("command name should have the right length", length)
	}

	if message.PayloadSize == nil {
		t.Error("payload size should not be nil")
	}

	if hex := fmt.Sprintf("%x", message.PayloadSize); hex != "00000000" {
		t.Error("payload size should have the right value", hex)
	}

	if length := len(message.PayloadSize); length != 4 {
		t.Error("payload size should have the right length", length)
	}

	if message.Checksum == nil {
		t.Error("checksum should not be nil")
	}

	if hex := fmt.Sprintf("%x", message.Checksum); hex != "5df6e0e2" {
		t.Error("checksum should have the right value", hex)
	}

	if length := len(message.Checksum); length != 4 {
		t.Error("checksum should have the right length", length)
	}
}

func TestNewVersionMessage(t *testing.T) {
	version := NewVersionMessage(true)
	if version == nil {
		t.Error("should return a version message")
	}

	vVersion := fmt.Sprintf("%x", version.Version)
	if vVersion != "7f110100" {
		t.Error("version number has not the right value", vVersion)
	}

	if len(version.Version) != 4 {
		t.Error("version number has not the right size", len(version.Version))
	}

	vServices := fmt.Sprintf("%x", version.Services)
	if vServices != "0100000000000000" {
		t.Error("should return value for Full Node", vServices)
	}

	version = NewVersionMessage(false)
	vServices = fmt.Sprintf("%x", version.Services)
	if vServices != "0000000000000000" {
		t.Error("should return value for Full Node", vServices)
	}

	if len(version.Services) != 8 {
		t.Error("should return the right length for services", len(version.Services))
	}

	timestamp := version.Timestamp
	if hex := fmt.Sprintf("%x", timestamp); hex == "" {
		t.Error("should return the right timestamp", hex)
	}

	if len(timestamp) != 8 {
		t.Error("should return the right timestamp size", len(timestamp))
	}

	fromIpPort := version.FromIpPort
	if hex := fmt.Sprintf("%x", fromIpPort); hex != "000000000000000000000000c0a8026a208d" {
		t.Error("should return the right from ip/port", hex)
	}

	if len(fromIpPort) != 18 {
		t.Error("should return the right from ip/port size", len(fromIpPort))
	}

	toIpPort := version.ToIpPort
	if hex := fmt.Sprintf("%x", toIpPort); hex != "000000000000000000000000c0a8026a479d" {
		t.Error("should return the right from ip/port", hex)
	}

	if len(toIpPort) != 18 {
		t.Error("should return the right from ip/port size", len(toIpPort))
	}

	nonce := version.Nonce
	if hex := fmt.Sprintf("%x", nonce); hex != "0000000000000000" {
		t.Error("should return the right nonce", hex)
	}

	userAgent := version.UserAgent
	if hex := fmt.Sprintf("%x", userAgent); hex != "2f6f6e65636f696e2d6274633a302e302e312f" {
		t.Error("should return the right user agent", hex)
	}

}

//func TestConnection(t *testing.T) {
//	peer, err := Connect("localhost", "18333")
//	defer peer.Connection.Close()
//	if err != nil {
//		t.Error("should connect successfully")
//	}
//}
