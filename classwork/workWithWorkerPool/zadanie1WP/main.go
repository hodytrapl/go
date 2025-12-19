package main

import (
	"fmt"
	"sync"
)

type Job struct{
	number int
}

type Result struct {
	job Job
	square int
}

func main(){
	var wg sync.WaitGroup
	inputs := []int{1, 5, 12, 5, 3, 8, 9}

	jobs := make(chan Job,len(inputs))
	results := make(chan Result, len(inputs))

	numWorkers := 3

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup){
			defer wg.Done()
			fmt.Printf("Воркер %d начал работу\n", id)
			for job  := range jobs{
				result:=job.number*job.number
				results<-Result{job: job, square: result}
			}
			
			fmt.Printf("Воркер %d завершил работу\n", id)
		}(i,jobs,results, &wg)
	}

	for _, num := range inputs {
		jobs <- Job{number: num}
	}

	close(jobs)

	wg.Wait()
	
	close(results)

	fmt.Println("\nРезультаты:")
	for result := range results {
		fmt.Printf("Число: %d, Квадрат: %d\n", result.job.number, result.square)
	}

}