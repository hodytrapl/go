package main

import "fmt"

func producer(ch chan<- int) {
	for i := 0; i <= 5; i++ {
		fmt.Println("sending...")
		ch <- i
	}
	close(ch)
	fmt.Println("producer done")
}

func consumer(ch <-chan int, done chan<- bool) {
	for v := range ch {
		fmt.Println("get: ", v)
	}
	fmt.Println("consumer done")
	done<-true
}

func main() {
	ch := make(chan int)
	done := make(chan bool)
	go producer(ch)
	go consumer(ch,done)
	<-done
	
	fmt.Println("main done")
}
