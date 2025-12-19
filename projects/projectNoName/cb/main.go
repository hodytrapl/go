package main

import (
	"sync"
	"fmt"
	"time"
)

type SafeCounter struct {
	mu sync.Mutex
	total int
}

func(c *SafeCounter) Inc(amount int){
	c.mu.Lock()
	c.total+=amount
	c.mu.Unlock()
}

func main(){
	var wg sync.WaitGroup
	results:=make(chan int,3)

	counter:= SafeCounter()

	for i:=1;i<=3;i++{
		wg.Add(1)
		go func(id int){
			defer wg.Done()

			res:= id*10

			results <- res
			counter.Inc(1)
			fmt.Printf("worker %d done work.\n",i)
		}(i)
	}

	go func(){
		wg.Wait()
		close(results)
	}

	for res:=range results{
		fmt.Println("get result:",res)
	}
	fmt.Println("all worker did:" counter.total)
}