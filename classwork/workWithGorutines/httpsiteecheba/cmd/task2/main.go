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
	go printNumbers(2, 5)
	go printNumbers(3, 5)
	go printNumbers(4, 5)
	go printNumbers(5, 5)

	// Временный rостыль: ждём, чтобы горутины успели закончить
	time.Sleep(1 * time.Second)

	fmt.Println("End lab1")

}
