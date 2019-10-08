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
		//fmt.Println("Listening...")
		message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Println("message from reader: " + message)
		readerMessageReceived <- message
		//fmt.Println("message is sent to the channel, sleeping...")
		time.Sleep(1 * time.Second)
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
