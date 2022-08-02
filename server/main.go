package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {

	fmt.Println("Launching server...")

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
	}
}
