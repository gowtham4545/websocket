package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	IP   string
	id   int
}

var host string
var count = 0
var status = false
var Clients = make(map[int]*Client)
var mutex = &sync.Mutex{}
var tMutex = &sync.Mutex{}
var cMutex = &sync.Mutex{}

func main() {
	if len(os.Args) > 1 {
		host = os.Args[1]
	} else {
		fmt.Println(fmt.Errorf("error: No host provided"))
		return
	}
	startServer(host)
}

func startServer(host string) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("\n\nError Starting Server\n", err)
		return
	}
	status = true
	defer listener.Close()
	fmt.Println("\nServer started at ", listener.Addr().String())

	go handleTerminal()

	acceptConnections(listener)
}

func acceptConnections(listener net.Listener) {
	for status {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error Accepting Connection:", err)
			continue
		}
		go handleNewConnection(connection)
	}
}

func handleNewConnection(connection net.Conn) {
	mutex.Lock()
	count++
	client := &Client{conn: connection, id: count}
	client.IP = connection.RemoteAddr().String()
	Clients[count] = client
	mutex.Unlock()
	go handleConnection(client)
}

func handleConnection(client *Client) {
	defer client.conn.Close()
	reader := bufio.NewReader(client.conn)
	log.Printf("Client %d (%s) entered the chat..\n", client.id, client.IP)

	for status {
		if peekError := checkClientConnection(reader); peekError != nil {
			handleClientDisconnect(client, peekError)
			return
		}
		handleClientMessage(client, reader)
	}
}

func checkClientConnection(reader *bufio.Reader) error {
	_, err := reader.Peek(1)
	return err
}

func handleClientMessage(client *Client, reader *bufio.Reader) {
	tMutex.Lock()
	cMutex.Lock()
	defer tMutex.Unlock()
	defer cMutex.Unlock()

	message, _ := reader.ReadString('\n')
	log.Printf("Client %.3d: %s", client.id, strings.TrimSpace(message))
}

func handleClientDisconnect(client *Client, err error) {
	if err == io.EOF {
		log.Printf("Client %s disconnected.\n", client.IP)
	} else {
		log.Printf("Error reading from %s: %v\n", client.IP, err)
	}
	log.Printf("Client %d (%s) exited..", client.id, client.IP)
	mutex.Lock()
	delete(Clients, client.id)
	mutex.Unlock()
}

func handleTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		if peekError := checkTerminalInput(reader); peekError != nil {
			fmt.Println(fmt.Errorf("server not taking inputs: %v", peekError))
			return
		}
		handleTerminalCommand(reader)
	}
}

func checkTerminalInput(reader *bufio.Reader) error {
	_, err := reader.Peek(1)
	return err
}

func handleTerminalCommand(reader *bufio.Reader) {
	tMutex.Lock()
	cMutex.Lock()
	defer tMutex.Unlock()
	defer cMutex.Unlock()

	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	command := strings.SplitN(text, " : ", 2)
	if len(command) < 2 {
		fmt.Println("Invalid command format. Use: <client_id> : <message>")
		return
	}
	id, err := strconv.Atoi(command[0])
	if err != nil || Clients[id] == nil {
		fmt.Println(fmt.Errorf("error: no such connection exists"))
		return
	}
	sendMessageToClient(id, command[1])
}

func sendMessageToClient(id int, message string) {
	conn := Clients[id].conn
	n, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Printf("Failed to send message to %s: %v\n", Clients[id].IP, err)
	} else {
		fmt.Printf("Sent %d bytes to (%v)\n", n, conn)
	}
}
