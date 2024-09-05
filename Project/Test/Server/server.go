package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	IP   string
	id   int
}

var listener net.Listener
var clients = make(map[int]*Client)
var Host = "localhost:9000"
var status = false
var rstrt = false
var err error
var count int

func Status() bool {
	return status
}

func RunServer() {
	status = false
	Start(Host)
	go Accept()
	for status {
	}
	defer StopServer()
}

func Start(ip string) {
	listener, err = net.Listen("tcp", ip)
	if err != nil {
		fmt.Println("\n\nError Starting Server\n", err)
		return
	}
	fmt.Println("\nServer started at ", listener.Addr().String())
	status = true
}

func restart() {
	time.Sleep(10 * time.Second)
	fmt.Print("Restarting...")
	Start(Host)
}

func Accept() {
	if !status {
		err = fmt.Errorf("server not listening")
		return
	}
	var connection net.Conn
	for status {
		connection, err = listener.Accept()
		if err != nil {
			continue
		}
		count++
		client := &Client{conn: connection, id: count}
		client.IP = connection.RemoteAddr().String()
		clients[count] = client
		go handleConnection(client)
		if !status && rstrt {
			restart()
		}
	}
	if !status && rstrt {
		restart()
	}
}

func StopServer() {
	listener.Close()
	rstrt = false
	status = false
	time.Sleep(5 * time.Second)
}

func handleConnection(client *Client) {
	defer client.conn.Close()
	defer fmt.Printf("Client %d (%s) exited..\n", client.id, client.IP)
	reader := bufio.NewReader(client.conn)
	fmt.Printf("Client %d (%s) entered the chat..\n", client.id, client.IP)

	for status {
		_, err := reader.Peek(1)
		if err == nil {
			message, _ := reader.ReadString('\n')
			fmt.Printf("Client %.3d: %s\n", client.id, strings.TrimSpace(message))
			continue
		} else {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected.\n", client.IP)
			} else {
				fmt.Printf("Error reading from %s: %v\n", client.IP, err)
			}
			delete(clients, client.id)
			return
		}
	}
}

func Send(val interface{}, message string) error {
	var conn net.Conn
	if !status {
		return fmt.Errorf("server teminated")
	}
	if reflect.TypeOf(val) == reflect.TypeOf(int(0)) {
		conn = clients[val.(int)].conn
		fmt.Println(val.(int))
	} else {
		conn = val.(net.Conn)
	}
	_, err = conn.Write([]byte(message + "\n"))
	return err
}

func GetErrorStatus() error {
	return err
}
