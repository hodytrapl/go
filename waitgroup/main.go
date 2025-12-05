package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Рабочий %d начал работу\n", id)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Рабочий %d закончил работу\n", id)
	
}

func main() {
	//channel := make(chan int)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}

	wg.Wait()

	fmt.Println("done")
}
