package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	langs := []string{"go", "java", "python", "rust"}
	var wg sync.WaitGroup
	fmt.Println("----start in errors----")

	for _, lang := range langs {
		wg.Add(1)
		go func(id string){
			defer wg.Done()
			time.Sleep(10*time.Millisecond)
			fmt.Printf("\t%s\n",id)
		}(lang)
	}
	wg.Wait()
	fmt.Println("    ----done----")
}