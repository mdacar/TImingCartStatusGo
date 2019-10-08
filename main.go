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
	}
}
