package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type status struct {
	internetStatus    string
	readerStatus      string
	lastRead          time.Time
	batteryPercentage int
	last5Reads        string
}

type cartstatus status

var lastHeartbeatChan chan time.Time
var lastHeartbeat time.Time
var lastConnectAttempt time.Time
var listener ConnectTcpListener
var lastTagReadTime time.Time

type tagRead struct {
	antenna  int
	tagID    string
	unixTime string
}

var last5TagReads []tagRead

//NewCartStatus works as a constructor
func NewCartStatus() *cartstatus {
	//Init the channel so we can start listening
	lastHeartbeatChan = make(chan time.Time)

	//go GetReaderStatus()
	cs := new(cartstatus)

	cs.internetStatus = GetInternetStatus()
	cs.readerStatus = "Unknown"
	cs.lastRead = lastTagReadTime

	var last5Reads strings.Builder
	for _, tr := range last5TagReads {
		last5Reads.WriteString(tr.tagID)
	}
	cs.last5Reads = last5Reads.String()
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
			//fmt.Println(latestHeartbeat)
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
	if listener.Status != Connecting && listener.Status != Listening {
		//Wait 10 seconds before trying to reconnect
		if time.Now().Sub(lastConnectAttempt).Seconds() > float64(10) {
			go listener.StartListening("192.168.1.79", "14150")
			lastConnectAttempt = time.Now()

			for {
				//fmt.Println("Checking for messages from the reader")
				select {
				case message := <-readerMessageReceived:
					//Heartbeat message
					if strings.ContainsAny(message, "*") {
						lastHeartbeatChan <- time.Now()
					}
					//Data message
					if strings.ContainsAny(message, ",") {
						HandleReaderData(message)
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

func HandleReaderData(data string) {
	lines := strings.Split(data, "\n")
	//loop through the lines we got back
	for _, line := range lines {
		parsedData := strings.Split(line, ",")
		//data is the correct length so proceed
		if len(parsedData) == 3 {
			var newTagRead tagRead
			newTagRead.antenna, _ = strconv.Atoi(parsedData[0])
			newTagRead.tagID = parsedData[1]
			newTagRead.unixTime = parsedData[2]

			//append to the last5TagReads slice, but first make sure we don't have more than 5
			if len(last5TagReads) == 5 {
				last5TagReads = last5TagReads[:len(last5TagReads)-1]
			}
			last5TagReads = append(last5TagReads, newTagRead)

			//not exact, should parse the unix time but this will be close enough
			lastTagReadTime = time.Now()
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
