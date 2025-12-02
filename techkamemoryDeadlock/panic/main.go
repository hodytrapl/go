package main

import "fmt"

func main() {

	ch := make(chan int)
	close(ch)
	//ch<-1 - panic
	val, ok := <-ch
	fmt.Println(val, " ", ok)
}
