package main

import (
	"fmt"
	"sync"
)

func factorial(n int) int{
	summ:=1;
	for i:=1; i<=n; i++{
		summ*=i;
	}
	return summ
}

func CalculateFactorial(i int, number int, results []int, wg *sync.WaitGroup){
	defer wg.Done()

	results[i]=factorial(number)
}

func main(){
	var wg sync.WaitGroup;
	numbers := []int{5, 2, 10, 7}
	results := make([]int, len(numbers))

	for i,number := range numbers{
		wg.Add(1)
		go CalculateFactorial(i,number,results,&wg)
	}

	wg.Wait()

	for i,number := range numbers{
		fmt.Printf("Факториал %d = %d\n",number, results[i])
	}

}