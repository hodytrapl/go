package main

import "fmt"

func generate(ch chan<- int, n int) {
	for i := 1; i < n; i++ {
		ch <- i
	}
	close(ch)
}

func main() {
	channel := make(chan int)
	go generate(channel, 10)
	sum := 0
	for i := range channel {
		sum += i
	}
	fmt.Println("sum := ",sum)
}