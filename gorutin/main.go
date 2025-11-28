package main

import (
	"fmt"
	"os"
	"time"
)

func homework(n int) {
	for i := 1; i <= n; i++ {
		fmt.Printf("[gorutine %d] message %s-%d", i, "messss",i)
	}
}

func main() {
	n:=os.Args[0]
	homework(n)
	time.Sleep(5 * time.Second)
}
