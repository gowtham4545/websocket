package tests

import (
	"net"
	client "test/Client"
	server "test/Server"
	"testing"
	"time"
)

func TestClientConnect(t *testing.T) {
	go server.RunServer()
	time.Sleep(5 * time.Second)
	conn := client.CreateClient()
	if conn == nil {
		t.Fatalf("Unable to connect to server.")
	}
	close(t, conn)
	server.StopServer()
}

func TestClientSendMessage(t *testing.T) {
	go server.RunServer()
	time.Sleep(5 * time.Second)
	conn := client.CreateClient()
	err := client.Send(&conn, "Hi")
	if err != nil {
		t.Errorf("Unable to send message")
	}
	close(t, conn)
	server.StopServer()
}

// func TestClientRecieveMessage(t *testing.T) {
// 	go server.RunServer()
// 	time.Sleep(5 * time.Second)
// 	conn := client.CreateClient()
// 	err := server.Send(1, "Hello")
// 	if err != nil {
// 		t.Errorf("Server unable to send message")
// 	}
// 	_, err = client.Receive(&conn)
// 	if err != nil {
// 		t.Errorf("Unable to recieve")
// 	}
// 	// close(t, conn)
// 	// server.StopServer()
// }

func close(t *testing.T, conn net.Conn) {
	if err := conn.Close(); err != nil {
		t.Errorf("Failed to close connection: %v", err)
	}
}
