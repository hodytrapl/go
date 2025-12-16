package main

import (
	"fmt"
	"strings"
	"sync"
)

func ConcurrentWordCount(sentences []string) map[string]int {
	result := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, sentence := range sentences {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			words := strings.Fields(s)

			mu.Lock()
			for _, word := range words {
				result[word]++
			}
			mu.Unlock()
		}(sentence)
	}

	wg.Wait()

	return result
}

func main() {
	text := []string{
		"quick brown fox",
		"lazy dog",
		"quick brown fox jumps",
		"jumps over lazy dog",
	}
	wordCount := ConcurrentWordCount(text)

	for word, count := range wordCount {
		fmt.Printf("%s: %d\n", word, count)
	}
}
