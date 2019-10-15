package main

import (
	"fmt"
	"time"
)

func main() {
	cartStatusChecker := NewCartStatus()
	for {
		cartstatus := cartStatusChecker.GetStatus()

		fmt.Printf("Internet: %v\n", cartstatus.internetStatus)
		fmt.Printf("Reader: %v\n", cartstatus.readerStatus)
		time.Sleep(5 * time.Second)
		fmt.Printf("Last Read: %v\n", cartstatus.lastRead)
		fmt.Printf("Last 5 Reads: %v\n", cartstatus.last5Reads)
		time.Sleep(5 * time.Second)
	}
}
