package main

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
var clients = make(map[int]*Client)
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
	listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("\n\nError Starting Server\n", err)
		return
	}

	fmt.Println("\nServer started at ", listener.Addr().String())
	defer listener.Close()

	go handleTerminal()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error Accepting Connection:", err)
			continue
		}
		mutex.Lock()
		count++
		client := &Client{conn: connection, id: count}
		client.IP = connection.RemoteAddr().String()
		clients[count] = client
		mutex.Unlock()
		go handleConnection(client)
	}
}

func handleConnection(client *Client) {
	defer client.conn.Close()
	reader := bufio.NewReader(client.conn)
	log.Printf("Client %d (%s) entered the chat..\n", client.id, client.IP)

	for {
		_, err := reader.Peek(1)
		if err == nil {
			tMutex.Lock()
			cMutex.Lock()
			message, _ := reader.ReadString('\n')
			log.Printf("Client %.3d: %s", client.id, strings.TrimSpace(message))
			tMutex.Unlock()
			cMutex.Unlock()
			continue
		} else {
			if err == io.EOF {
				log.Printf("Client %s disconnected.\n", client.IP)
			} else {
				log.Printf("Error reading from %s: %v\n", client.IP, err)
			}
			log.Printf("Client %d (%s) exited..", client.id, client.IP)
			mutex.Lock()
			delete(clients, client.id)
			mutex.Unlock()
			return
		}
	}
}

func handleTerminal() {
	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := reader.Peek(1)
		if err == bufio.ErrBufferFull {
			continue
		} else if err == nil {
			tMutex.Lock()
			cMutex.Lock()
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			command := strings.SplitN(text, " : ", 2)
			if len(command) < 2 {
				fmt.Println("Invalid command format. Use: <client_id> : <message>")
				continue
			}
			id, err := strconv.Atoi(command[0])
			if err != nil || clients[id] == nil {
				fmt.Println(fmt.Errorf("error: no such connection exists"))
				continue
			}
			conn := clients[id].conn
			n, err := conn.Write([]byte(command[1] + "\n"))
			if err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", clients[id].IP, err)
			} else {
				fmt.Printf("Sent %d bytes to (%v)\n", n, conn)
			}
			tMutex.Unlock()
			cMutex.Unlock()
		} else {
			fmt.Println(fmt.Errorf("server not taking inputs: %v", err))
			return
		}
	}
}
