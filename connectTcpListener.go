package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting to listen")
	readerIp := "192.168.1.79:14150"

	conn, _ := net.Dial("tcp", readerIp)
	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
	}

}
