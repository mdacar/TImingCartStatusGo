package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type status struct {
	internetStatus    string
	readerStatus      string
	lastRead          time.Time
	batteryPercentage int
}

type cartstatus status

var lastHeartbeatChan chan time.Time
var lastHeartbeat time.Time
var lastConnectAttempt time.Time
var listener ConnectTcpListener

//NewCartStatus works as a constructor
func NewCartStatus() *cartstatus {
	//Init the channel so we can start listening
	lastHeartbeatChan = make(chan time.Time)

	//go GetReaderStatus()
	cs := new(cartstatus)

	cs.internetStatus = GetInternetStatus()
	cs.readerStatus = "Unknown"
	return cs
}

func (cs cartstatus) GetStatus() cartstatus {

	go GetReaderStatus()
	cs.internetStatus = GetInternetStatus()
	latestHeartbeat := time.Now()
	nullTime := time.Time{}

	//Loop over the channel to get to the latest value
	for latestHeartbeat != nullTime {
		select {
		case latestHeartbeat = <-lastHeartbeatChan:
			lastHeartbeat = latestHeartbeat
		default:
			fmt.Println(latestHeartbeat)
			latestHeartbeat = time.Time{}
			//fmt.Println("\t\tNo heartbeat messages")
		}
	}

	if time.Now().Sub(lastHeartbeat).Seconds() < float64(20) {
		cs.readerStatus = "Reading"
	} else {
		cs.readerStatus = "Offline"
	}

	return cs
}

//GetReaderStatus checks to see if the RFID reader is transmitting the heartbeat over TCP port 14500
func GetReaderStatus() {
	if !listener.IsListening {
		//Wait 10 seconds before trying to reconnect
		if time.Now().Sub(lastConnectAttempt).Seconds() > float64(10) {
			go listener.StartListening("192.168.1.79", "14150")
			lastConnectAttempt = time.Now()

			for {
				//fmt.Println("Checking for messages from the reader")
				select {
				case message := <-readerMessageReceived:
					//fmt.Println("Message received: " + message)
					if strings.ContainsAny(message, "*") {
						//fmt.Println("\tSending the time to the lastHeartbeat channel...")
						lastHeartbeatChan <- time.Now()
					}
				default:
					//fmt.Println("No messages from the reader")
				}

				//fmt.Println("\tGetting reader status...")
				time.Sleep(3 * time.Second)
			}
		}
	}
}

//GetInternetStatus checks if the internet is connected and returns Down or Online
func GetInternetStatus() string {
	url := "http://www.google.com/generate_204"
	_, err := http.Get(url)
	if err != nil {
		return "Down"
	}

	return "Online"
}
