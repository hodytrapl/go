package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var wg sync.WaitGroup
	//var mu sync.

	iterations := 1000
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			//mu.Lock()
			counter++
			//mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("wait ", iterations)
	fmt.Println("get ", counter)
}
