package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		fmt.Println("Горутина 1: Привет!")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Горутина 2: мир!")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Горутина 3: golang!")
	}()
	wg.Wait()
}
