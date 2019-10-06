package main

import (
	"net/http"
	"time"
)

type status struct {
	internetStatus    string
	readerStatus      string
	lastRead          time.Time
	batteryPercentage int
}

type cartstatus status

func (cs cartstatus) GetStatus() status {
	var currentStatus status
	currentStatus.internetStatus = GetInternetStatus()
	currentStatus.readerStatus = GetReaderStatus()
	return currentStatus
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
func GetReaderStatus() string {


	return "Unknown"
}
