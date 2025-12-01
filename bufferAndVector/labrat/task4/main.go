package main

import "fmt"

func generate(ch chan<- int) {
	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch) // ответ: будет деадлок
}

func square(in <-chan int, out chan<- int) {
	for num := range in {
		out <- num * num
	}
	close(out)
}

func print(out <-chan int, done chan bool) {
	for v := range out {
		fmt.Println(v)
	}
	done <- true
}

func main() {
	chan1 := make(chan int)
	chan2 := make(chan int)
	done := make(chan bool)
	go generate(chan1)
	go square(chan1, chan2)
	go print(chan2, done)

	<-done
}
