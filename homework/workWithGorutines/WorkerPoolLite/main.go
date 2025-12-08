package main

import (
	"fmt"
	"time"
)

type Task struct {
	ID          int
	TimeCreated string
}

func worker(id int, tasks <-chan Task) {
	for task := range tasks {
		fmt.Println("Воркер ", id, " начал задачу ", task)
		time.Sleep(time.Millisecond * 500)
		fmt.Println("Воркер ", id, " закончил задачу ", task)
	}

}

func main() {
	taskCh := make(chan Task, 10)
	go worker(1, taskCh)
	for i := 1; i <= 10; i++ {
		task := Task{
			ID:          i,
			TimeCreated: time.Now().Format("15:04:05.000"),
		}
		fmt.Println("Отправлена задача ", i)
		taskCh <- task
	}
	close(taskCh)
	time.Sleep(time.Second * 3)
}
