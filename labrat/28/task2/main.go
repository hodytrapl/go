package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	domains := []string{"google.com", "yandex.ru", "github.com", "stackoverflow.com"}
	var wg sync.WaitGroup

	wg.Add(len(domains))

	for _, domain := range domains {
		
		go func(s string) {
			defer wg.Done()
			fmt.Printf("initializing check %s\n", s)
			duration := time.Duration(rand.Intn(3)) * time.Second
			time.Sleep(duration)
			fmt.Printf("init %s done\n", s)
		}(domain)
	}

	wg.Wait()
	fmt.Println("done")
}
