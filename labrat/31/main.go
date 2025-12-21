package main

import (
	"fmt"
	"time"
)

func newsFeed(channel chan<- string) {
	for i := 1; i <= 10; i++ {
		channel <- fmt.Sprintf("Новости из ленты новостей #%d", i)
		time.Sleep(800 * time.Millisecond)
	}
	close(channel)
}

func socialMedia(channel chan<- string) {
	for i := 1; i <= 10; i++ {
		channel <- fmt.Sprintf("блог #%d", i)
		time.Sleep(800 * time.Millisecond)
	}
	close(channel)
}

func main() {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go newsFeed(channel1)
	time.Sleep(1 * time.Second)
	go socialMedia(channel2)

	timeout := time.After(5 * time.Second)

	for {
		select {
		case msg1 := <-channel1:
			println("Новость:", msg1)
		case msg2 := <-channel2:
			println("Соцсети:", msg2)
		case <-timeout:
			println("Время вышло, выключаемся")
			return
		default:
			time.Sleep(500 * time.Millisecond)
			println(".")
		}
	}
}
