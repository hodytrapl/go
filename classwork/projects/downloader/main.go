package main

import (
    "bufio"
    "flag"
    "fmt"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"
)

func ReadUrls(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var urls []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := strings.TrimSpace(scanner.Text())
        if url != "" {
            urls = append(urls, url)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return urls, nil
}


func checkURL(url string, timeout time.Duration, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()
    start := time.Now()
    client := &http.Client{
        Timeout: timeout,
    }
    resp, err := client.Get(url)
    if err != nil {
        elapsed := time.Since(start)
        results <- fmt.Sprintf("[err] `%s` (%v) - %v", url, elapsed, err)
        return
    }
    defer resp.Body.Close()

    elapsed := time.Since(start)
    result := fmt.Sprintf("[%d] `%s` (%v)", resp.StatusCode, url, elapsed)
    results <- result
}

func main() {
    timeout := flag.Duration("t", 10*time.Second, "Timeout for HTTP requests")
    flag.Parse()

    urls, err := ReadUrls("urls.txt")
    if err != nil {
        fmt.Printf("Error read URLs: %v\n", err)
        os.Exit(1)
    }

    if len(urls) == 0 {
        fmt.Println("No URLs found in urls.txt.")
        return
    }

    fmt.Printf("%d URLs timeout %v...\n\n", len(urls), *timeout)

    var wg sync.WaitGroup
    results := make(chan string, len(urls))

    for _, url := range urls {
        wg.Add(1)
        go checkURL(url, *timeout, &wg, results)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        fmt.Println(result)
    }
}