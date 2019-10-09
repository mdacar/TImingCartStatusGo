package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const ReaderEventName = "ReaderEvent"

var readerMessageReceived chan string

type ConnectTcpListener struct {
	IsListening   bool
	LastHeartbeat time.Time
	IPAddress     string
	Port          string
}

func (l ConnectTcpListener) StartListening(ipAddress string, port string) {
	//fmt.Println("Connecting...")
	readerMessageReceived = make(chan string)
	l.IPAddress = ipAddress
	l.Port = port
	conn, err := net.Dial("tcp", ipAddress+":"+port)
	if err != nil {
		fmt.Println(err)
		l.IsListening = false
		return
	}
	l.IsListening = true
	for {

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			l.IsListening = false
			fmt.Printf("Failed to connect: %v\n", err)
			break
		} else {
			readerMessageReceived <- message
			time.Sleep(1 * time.Second)
		}
	}
}

/*
type ReaderHeartbeat struct {
	Time time.Time
}

type TagRead struct {
	TagId     string
	Timestamp string
	Antenna   string
}
*/

//ReaderData contains the output from the reader.  Could be either a tag read or a heartbeat
type ReaderData string
