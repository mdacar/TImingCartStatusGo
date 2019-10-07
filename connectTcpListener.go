package main

import (
	"bufio"
	"net"
	"time"
)

type ConnectTcpListener struct {
	IsListening   bool
	LastHeartbeat time.Time
	IPAddress     string
	Port          string
}

func (l ConnectTcpListener) StartListening(ipAddress string, port string) {
	l.IPAddress = ipAddress
	l.Port = port
	conn, err := net.Dial("tcp", ipAddress+":"+port)
	if err != nil {
		l.IsListening = false
		return
	}
	l.IsListening = true
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')

	}
}

func

type ReaderOutput interface

type ReaderHeartbeat struct {
	Time time.Time
}

type TagRead struct {
	TagId     string
	Timestamp string
	Antenna   string
}

type ReaderData string
