package main

func inputStream(n chan<- int) {
	for i := 0; i < 10; i++ {
		n <- i
	}
	close(n)
}

func doubler(n <-chan int, d chan<- int) {
	for v := range n {
		d <- v * 2
	}
	close(d)
}

printer := func(d <-chan int) {
	for v := range d {
		fmt.Println("Результат: ",v)
	}
	done <- true
}

func main(){
	firstChan := make(chan int)
	secondChan := make(chan int)
	done := make(chan bool)

	go inputStream(firstChan)
	go doubler(firstChan, secondChan)
	go printer(secondChan)	
	<-done
}