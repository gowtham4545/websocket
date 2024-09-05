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
	"time"
)

var client string
var tMutex = &sync.Mutex{}
var cMutex = &sync.Mutex{}
var waitTime = 5

func main() {
	if len(os.Args) > 1 {
		client = os.Args[1]
		if len(os.Args) >= 3 {
			waitTime, _ = strconv.Atoi(os.Args[2])
		}
	} else {
		fmt.Println(fmt.Errorf("error: Server IP not provided"))
		return
	}
	for {
		conn, err := net.Dial("tcp", client)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			fmt.Printf("Retrying in %d seconds...\n", waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
			continue
		}
		fmt.Println("Connection established with server..")
		defer conn.Close()

		go handleTerminal(&conn)

		cReader := bufio.NewReader(conn)
		for {
			_, err := cReader.Peek(1)
			if err == nil {
				tMutex.Lock()
				cMutex.Lock()
				message, _ := cReader.ReadString('\n')
				message = strings.Trim(message, "\n")
				log.Println("Server:", message)
				tMutex.Unlock()
				cMutex.Unlock()
				continue
			} else {
				if err == io.EOF {
					fmt.Printf("Server disconnected: %v\n", err)
				} else {
					fmt.Printf("Error reading from server: %v\n", err)
				}
				break
			}
		}
		fmt.Println("Reconnecting to server...")
		time.Sleep(time.Second * time.Duration(waitTime))
	}
}

func handleTerminal(conn *net.Conn) {
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
			if text == "" {
				fmt.Println("Nothing to send. Enter text to send...")
				continue
			}
			n, err := (*conn).Write([]byte(text + "\n"))
			if err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", (*conn).RemoteAddr().String(), err)
			} else {
				fmt.Printf("Sent %d bytes to (%v)\n", n, conn)
			}
			tMutex.Unlock()
			cMutex.Unlock()
		} else {
			fmt.Println(fmt.Errorf("client not taking inputs: %v", err))
			return
		}
	}
}
