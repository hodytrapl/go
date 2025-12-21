package main

func main(){
	fast :=make(chan string)
	slow :=make(chan string)

	go func(){
		time.Sleep(3*time.Second)
		slow<-"im slowing..."
	}()

	go func(){
		time.Sleep(3*time.Second)
		fast<-"im slowing..."
	}()

	select{
	case msg1:=<-fast:
		fmt.Println("fast win.",msg1)
	case msg2:=<-slow:
		fmt.Println("slow win.",msg2)
	case <-time.After(2*time.Second):
		fmt.Println("no winner")
	
	}
}