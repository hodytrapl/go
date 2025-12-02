package main

import (
	"fmt"
	"time"
)

func logger(logs <-chan string) {
	for log := range logs {
		fmt.Println("[LOG]: ",log)
	}
}

func main() {
	logCh := make(chan string, 5)
	go logger(logCh)
	for i := 1; i <= 8; i++ {
        logCh <- fmt.Sprintf("Message %d", i)
    }
	time.Sleep(2 * time.Second)
}
