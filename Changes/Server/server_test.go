package server

import (
	// "bufio"

	"net"
	client "project/Client"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	host = "localhost:9000"
	listener, err := net.Listen("tcp", host)
	if err != nil {
		t.Errorf("Error Starting Server: %v\n", err)
		return
	}
	defer listener.Close()
	status = false
}

func TestServerAccept(t *testing.T) {
	host = "localhost:9000"
	listener, err := net.Listen("tcp", host)
	if err != nil {
		t.Errorf("Error Starting Server: %v\n", err)
		return
	}
	status = true
	defer listener.Close()
	go acceptConnections(listener)
	time.Sleep(2 * time.Second)
	clientConn, err := client.ConnectToServer(host)
	if err != nil {
		t.Fatalf("client unable to connect to server")
	}
	clientConn.Close()
	status = false
}

func TestServerRecieve(t *testing.T) {
	host = "localhost:9005"
	listener, err := net.Listen("tcp", host)
	if err != nil {
		t.Errorf("Error Starting Server: %v\n", err)
		return
	}
	status = true
	// connection, err := listener.Accept()
	// if err != nil {
	// 	t.Errorf("Error Accepting Connection: %v", err)
	// }
	// go handleNewConnection(connection)
	go acceptConnections(listener)
	time.Sleep(2 * time.Second)
	clientConn, err := client.ConnectToServer(host)
	if err != nil {
		t.Fatalf("client unable to connect to server")
	}
	client.SendMessageToServer(&clientConn, "hello")
	// mes, err := bufio.NewReader(cli).ReadString('\n')
	// if err != nil {
	// 	t.Errorf("Unable to read message from client")
	// }
	// if mes != "hello" {
	// 	t.Errorf("Expexted 'hello'; got %s\n, ", mes)
	// }
	// connection.Close()
	clientConn.Close()
	status = false
	listener.Close()
}
