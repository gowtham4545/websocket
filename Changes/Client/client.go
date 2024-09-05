package client

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
	startClient()
}

func startClient() {
	for {
		conn, err := ConnectToServer(client)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			fmt.Printf("Retrying in %d seconds...\n", waitTime)
			time.Sleep(time.Second * time.Duration(waitTime))
			continue
		}
		fmt.Println("Connection established with server..")

		go HandleTerminal(&conn)

		handleServerMessages(conn)
		fmt.Println("Reconnecting to server...")
		time.Sleep(time.Second * time.Duration(waitTime))
		conn.Close()
	}
}

func ConnectToServer(serverAddress string) (net.Conn, error) {
	return net.Dial("tcp", serverAddress)
}

func handleServerMessages(conn net.Conn) {
	cReader := bufio.NewReader(conn)
	for {
		if peekError := checkServerConnection(cReader); peekError != nil {
			handleServerDisconnect(peekError)
			break
		}
		handleIncomingMessage(cReader)
	}
}

func checkServerConnection(reader *bufio.Reader) error {
	_, err := reader.Peek(1)
	return err
}

func handleIncomingMessage(reader *bufio.Reader) {
	tMutex.Lock()
	cMutex.Lock()
	defer tMutex.Unlock()
	defer cMutex.Unlock()

	message, _ := reader.ReadString('\n')
	message = strings.Trim(message, "\n")
	log.Println("Server:", message)
}

func handleServerDisconnect(err error) {
	if err == io.EOF {
		fmt.Printf("Server disconnected: %v\n", err)
	} else {
		fmt.Printf("Error reading from server: %v\n", err)
	}
}

func HandleTerminal(conn *net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		if peekError := checkTerminalInput(reader); peekError != nil {
			fmt.Println(fmt.Errorf("client not taking inputs: %v", peekError))
			return
		}
		handleUserInput(conn, reader)
	}
}

func checkTerminalInput(reader *bufio.Reader) error {
	_, err := reader.Peek(1)
	return err
}

func handleUserInput(conn *net.Conn, reader *bufio.Reader) {
	tMutex.Lock()
	cMutex.Lock()
	defer tMutex.Unlock()
	defer cMutex.Unlock()

	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		fmt.Println("Nothing to send. Enter text to send...")
		return
	}
	SendMessageToServer(conn, text)
}

func SendMessageToServer(conn *net.Conn, message string) {
	n, err := (*conn).Write([]byte(message + "\n"))
	if err != nil {
		fmt.Printf("Failed to send message to %s: %v\n", (*conn).RemoteAddr().String(), err)
	} else {
		fmt.Printf("Sent %d bytes to (%v)\n", n, conn)
	}
}
