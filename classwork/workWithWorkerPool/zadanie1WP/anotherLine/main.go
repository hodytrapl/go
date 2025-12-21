package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID     int
	Number int
}

type Result struct {
	jobID     int
	InputNum  int
	Square    int
	perimeter int
	workerID  int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

		square := job.Number * job.Number
		perimeter := 4 * job.Number

		results <- Result{
			jobID:     job.ID,
			InputNum:  job.Number,
			Square:    square,
			perimeter: perimeter,
			workerID:  id,
		}
	}
	fmt.Printf("Воркер %d завершил работу\n", id)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var wg sync.WaitGroup
	inputs := []int{1, 5, 12, 5, 3, 8, 9}

	jobs := make(chan Job)
	results := make(chan Result)

	const workers = 3

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		fmt.Println("start sending task...")
		for i, num := range inputs {
			jobs <- Job{ID: i + 1, Number: num}
		}
		close(jobs)
		fmt.Println("задач нету")
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\nРезультаты:")
	for result := range results {
		fmt.Printf("Число: %d, Квадрат: %d\n", result.InputNum, result.Square)
	}
}