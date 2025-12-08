package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10000; i++ {
		go func() {
			fmt.Printf("Goroutine %d says hello\n", i)
		}()
	}

	time.Sleep(10000 * time.Millisecond)

}

//чем ближе к нулю, тем  хаотичнее числа, времени надо гдето в 2.5 раза больше
