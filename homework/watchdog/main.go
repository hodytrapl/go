package main

import (
	"time"
)

func worker(channel chan<-string){
	defer close(channel)
	for i:=0;i<10;i++{
		channel<-"ping"
		time.Sleep(1 * time.Second)
	}
	
	time.Sleep(3 * time.Second)
}

func main(){
	channel :=make(chan string,1)
	go worker(channel)

	timeout := time.After(2* time.Second)

	for {
		select {
		case msg1 := <-channel:
			println("живой: ", msg1)
		case <-timeout:
			println("авария, он сдох!")
			return
		}
	}
}