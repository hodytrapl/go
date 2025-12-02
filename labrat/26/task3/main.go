package main

import "fmt"

func generate(ch chan <-int){
	for i := 1; i <= 10; i++ {
		ch<-i
	}
	close(ch) // ответ: будет деадлок
}

func worker(id int,ch chan int, done chan bool){
	for v :=range ch{
		fmt.Println("Рабочий: ",id,"получил: ",v)
	}
	done <-true
}



func main() {
	ch := make(chan int)
	done := make(chan bool)
	go generate(ch)
	go worker(1,ch,done)
	go worker(2,ch,done)
	<-done
	<-done
	fmt.Println("done main")
}
