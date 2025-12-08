package main

import (
	"fmt"
	"time"
)

func printNumbers(id int, n int) {
	for i := 1; i <= n; i++ {
		fmt.Println("Worker: ", id, "Number: ", i)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println("Start lab1")
	go printNumbers(1, 5)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("End lab1")

}
