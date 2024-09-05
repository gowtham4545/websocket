package main

import (
	"net"
	"testing"
)

func connect(t *testing.T, ip string) net.Conn {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		t.Errorf("Unable to connect to server.")
	}
	return conn
}

func message(t *testing.T, conn *net.Conn, message string) {
	_, err := (*conn).Write([]byte(message + "\n"))
	if err != nil {
		t.Fatalf("Error sending message: %s", err)
	}
}

func TestClientConnect(t *testing.T) {
	conn := connect(t, "localhost:8180")
	conn.Close()
}

func TestClientMessage(t *testing.T) {
	conn := connect(t, "localhost:8180")
	message(t, &conn, "Hello!")
	conn.Close()
}
