package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"sync"
)

type Datafile struct{
	nameFile string;
	number int;
}

func Source(files []string,channel chan<-string, wg *sync.WaitGroup){
	defer wg.Done()
	for _, file := range files {
		channel <- file
	}
	close(channel)
}

func Filter(in <-chan string, out chan<- string, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(out)
	
	var substr = ".txt"

	for file :=range in{
		if(strings.Contains(file, substr)){
			out<-file
		}
	}
}

func Processing(in <-chan string, out chan<- Datafile, wg *sync.WaitGroup){
	defer wg.Done()
	for file := range in {
		sleepTime := time.Duration(rand.Intn(401)+100) * time.Millisecond
		time.Sleep(sleepTime)
		
		result := Datafile{
			nameFile: file,
			number:   rand.Intn(1000) + 1,
		}
		out <- result
	}
}

func FanIn(in <-chan Datafile,out chan<- int, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(out)
	total:=0
	for file :=range in{
		total+=file.number
	}
	out <- total
}

func main(){
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	str:= []string{"data.txt", "image.png", "info.txt", "backup.zip"}

	sourceChan := make(chan string, 10)
	filterChan := make(chan string, 10)
	resultChan := make(chan Datafile, 10)
	resultTotalChan := make(chan int, 1)

	wg.Add(1)
	go Source(str,sourceChan,&wg)

	wg.Add(1)
	go Filter(sourceChan,filterChan,&wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go Processing(filterChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	wg.Add(1)
	go FanIn(resultChan, resultTotalChan, &wg)

	total := <-resultTotalChan
	fmt.Println("всего обработано строк:", total)

	wg.Wait()
}