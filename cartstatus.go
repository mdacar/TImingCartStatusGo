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
	return currentStatus
}

func GetInternetStatus() string {
	url := "http://www.google.com/generate_204"
	_, err := http.Get(url)
	if err != nil {
		return "Down"
	}

	return "Online"
}
