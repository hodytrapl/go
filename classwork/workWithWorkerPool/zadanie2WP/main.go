package main

import(
	"fmt"
	"time"
)

func worker(id int,quit<-chan bool){
	for{
		select{
		case<-quit:
			fmt.Printf("worker %d going home\n",id)
			return
		default:
			fmt.Printf("worker %d did work\n",id)
			time.Sleep(500*time.Millisecond)
		}
	}
}

func main(){
	quit:=make(chan bool)
	go worker(1,quit)

	time.Sleep(2*time.Second)

	fmt.Println("kinger: stop work!\n")
	close(quit)

	time.Sleep(1*time.Second)
}