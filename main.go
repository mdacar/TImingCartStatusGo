package main

import (
	"fmt"
	"time"
)

func main() {
	var cartStatusChecker cartstatus
	for {
		curStatus := cartStatusChecker.GetStatus()

		fmt.Printf("Internet: %v\n", curStatus.internetStatus)
		fmt.Printf("Reader: %v\n", curStatus.readerStatus)
		time.Sleep(5 * time.Second)
	}
}
