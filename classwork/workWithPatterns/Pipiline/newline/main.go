package main

func generator(limit int) <-chan int{
	out :=make(chan int)
	go func(){
		defer close(out)
		for i:=1;i<=limit;i++{
			out<-i
		}
	}()
	return out
}
//fan-out
func worker(id int,in <-chan int) <-chan int{
	out:=make(chan int)
	go func(){
		defer close(out)
		for n :=range in{
			fmt.Printf("worker %d processsing %d\n",id,n)
			out<-n*n
		}
	}()
	return out
}

//fan-in
func merge(cs ...<-chan int) <-chan int{
	out:=make(chan int)
	var wg sync.WaitGroup

	output:= func(c <-chan int){
		defer wg.Done()
		for n:=range c{
			out<-n
		}
	}

	wg.Add(len(cs))
	for _,c :=range cs{
		go output(c)
	}

	go func(){
		wg.Wait()
		close(out)
	}()

	return out
}

func main(){

	in:=generator(20)

	c1:=worker(1,in)
	c2:=worker(1,in)
	c3:=worker(1,in)

	out:=merge(c1,c2,c3)

	for result:=range out{
		fmt.Println("result: ",result)
	}
}