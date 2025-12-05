package main

import (
    "flag"
    "fmt"
    "sync"
)

const (
    messagePerGoroutine = 3
)

func printMessage(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < messagePerGoroutine; i++ {
        fmt.Printf("[goroutine %d] message %d\n", id, i)
    }
}

func main() {
    countFlag := flag.Int("count", 2, "сколько горутин запустить")
    flag.Parse()
    count := *countFlag
    var wg sync.WaitGroup
    wg.Add(count)

    for i := 0; i < count; i++ {
        go printMessage(i, &wg)
    }

    wg.Wait()
}
