package main

import (
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

//NewCartStatus works as a constructor
func NewCartStatus() *cartstatus {
	//Init the channel so we can start listening
	lastHeartbeatChan = make(chan time.Time)

	go GetReaderStatus()
	cs := new(cartstatus)

	cs.internetStatus = GetInternetStatus()
	cs.readerStatus = "Unknown"
	return cs
}

func (cs cartstatus) GetStatus() cartstatus {

	cs.internetStatus = GetInternetStatus()
	//fmt.Println(lastHeartbeat)
	//fmt.Println(time.Now().Sub(lastHeartbeat).Seconds())

	select {
	case latestHeartbeat := <-lastHeartbeatChan:
		lastHeartbeat = latestHeartbeat
	default:
		//cs.readerStatus = "unknown"
	}

	if time.Now().Sub(lastHeartbeat).Seconds() < float64(20) {
		cs.readerStatus = "Reading"
	} else {
		cs.readerStatus = "Offline"
	}

	return cs
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

//GetReaderStatus checks to see if the RFID reader is transmitting the heartbeat over TCP port 14500
func GetReaderStatus() {
	var listener ConnectTcpListener

	go listener.StartListening("192.168.1.79", "14150")

	for {
		//fmt.Println("Checking for messages from the reader")
		select {
		case message := <-readerMessageReceived:
			//fmt.Println("Message received: " + message)
			if strings.ContainsAny(message, "*") {
				//fmt.Println("Sending the time to the lastHeartbeat channel")
				lastHeartbeatChan <- time.Now()
				//fmt.Println("Time is sent to the lastHeartbeat channel")
			}
		default:
			//fmt.Println("No messages from the reader")
		}

		//fmt.Println("GetReaderStatus is Sleeping...")
		time.Sleep(1 * time.Second)
	}
}
