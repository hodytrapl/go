type Job interface{
	Process() error
}

func worker(in <-chan Job){
	for job :=range in{
		job.Process()
	}
}

func gen(nums ...int)<-chan int{
	out:=make(chan int)

	go func(){
		for _,n:=range nums{
			out<-n
		}
		close(out)
	}
	return func sq(in)
}

func sq(in<-chan int) <-chan int{
	out:=make(chan int)
	go func(){
		for n :=range in{
			out<-n*n
		}
		close(out)
	}
	return out
}