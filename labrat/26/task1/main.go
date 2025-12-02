package main

import "fmt"

func main() {
	ch := make(chan string)
	go func(ch chan string) {
		ch <- "hello world!"
	}(ch)
	fmt.Println(<-ch)

}
