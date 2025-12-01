package main

import "fmt"

func generate(ch chan <-int){
	for i := 1; i <= 10; i++ {
		ch<-i
	}
	close(ch) // ответ: будет деадлок
}



func main() {
	ch := make(chan int)
	go generate(ch)
	for v :=range ch{
		fmt.Println(v)
	}
}
