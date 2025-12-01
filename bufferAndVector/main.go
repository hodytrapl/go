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

/*func consumer(ch <-chan int) {
	for v := range ch {
		fmt.Println("get: ", v)
	}
	fmt.Println("consumer done")
}*/

func main() {
	ch := make(chan int)
	go producer(ch)
	//go consumer(ch)
	val,ok:=<-ch
	fmt.Println(val,ok)

	fmt.Println("main done")
}
