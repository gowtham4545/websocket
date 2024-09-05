package tests

import (
	client "test/Client"
	server "test/Server"
	"testing"
)

func TestStartServer(t *testing.T) {
	server.Start(server.Host)
	err := server.GetErrorStatus()
	if err != nil {
		t.Errorf("%v", err)
	}
	server.StopServer()
}

func TestServerAccepting(t *testing.T) {
	server.Start(server.Host)
	go server.Accept()
	client := client.CreateClient()
	if client == nil {
		t.Fatalf("Unable to create Client")
	}
	server.StopServer()
}

func TestServerSendMessage(t *testing.T) {
	server.Start(server.Host)
	go server.Accept()
	conn := client.CreateClient()
	err := server.Send(conn, "hello")
	if err != nil {
		t.Errorf("Unable to send message")
	}
	conn.Close()
	server.StopServer()
}

func TestServerSendMessageById(t *testing.T) {
	server.Start(server.Host)
	go server.Accept()
	conn := client.CreateClient()
	err := server.Send(1, "hello")
	if err != nil {
		t.Errorf("Unable to send message")
	}
	conn.Close()
	server.StopServer()
}
