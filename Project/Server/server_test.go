package main

import (
	"bufio"
	"net"
	"testing"
)

func listen(t *testing.T, ip string) net.Listener {
	listener, err := net.Listen("tcp", ip)
	if err != nil {
		t.Errorf("\n\nError Starting Server\n%v\n", err)
	}
	return listener
}

func accept(t *testing.T, listener *net.Listener) net.Conn {
	conn, err := (*listener).Accept()
	if err != nil {
		t.Errorf("Error accepting connection: %v\n", err)
	}
	return conn
}

func message(t *testing.T, conn *net.Conn) {
	reader := bufio.NewReader(*conn)
	_, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Error reading messages:%v\n", err)
		return
	}
}

func TestListen(t *testing.T) {
	listener := listen(t, "localhost:8180")
	listener.Close()
}

func TestAccept(t *testing.T) {
	listener := listen(t, "localhost:8180")
	conn := accept(t, &listener)
	conn.Close()
	listener.Close()
}

func TestRecieveMessage(t *testing.T) {
	listener := listen(t, "localhost:8180")
	conn := accept(t, &listener)
	message(t, &conn)
	conn.Close()
	listener.Close()
}

