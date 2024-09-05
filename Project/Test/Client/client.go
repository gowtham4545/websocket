package client

import (
	"bufio"
	"fmt"
	"net"
	server "test/Server"
	"time"
)

var err error

func CreateClient() net.Conn {
	conn, err := connect(server.Host)
	if err != nil {
		fmt.Print(err)
		return nil
	}
	return conn
}

func connect(ip string) (net.Conn, error) {
	return net.Dial("tcp", ip)
}

func Receive(conn *net.Conn) (string, error) {
	(*conn).SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(*conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return message, nil
}

func Send(conn *net.Conn, message string) error {
	_, err = (*conn).Write([]byte(message + "\n"))
	return err
}
