package main 

import (
	"fmt"
	"math/rand"
	//"sync"
)

type Order struct {
    ID     int
    Amount int    // сумма заказа
    Status string // "new", "processed", "paid"
}

func generateOrders(count int) <-chan Order{
	order:=make(chan Order,100)

	go func(){
		defer close(order)

		for i:=0;i<count;i++{
			order<-Order{
				ID:i,
				Amount:rand.Intn(200)+1,
				Status:"new",
			}
		}
	}()
	
	return order
}

func processOrders(in <-chan Order) <-chan Order{
	orders:=make(chan Order,100)
	
	go func(){
		defer close(orders)

		for order:=range in{
			order.Status="processed"
			orders<-order
		}
	}()
	return orders
}

func filterOrders(in <-chan Order, minAmount int) <-chan Order{
	result:=make(chan Order,100)
	go func(){
		defer close(result)
		for order :=range in{
			if order.Amount>minAmount{
				result<-order
			}
	}
	}()

	
	return result
}

func main(){
	gen:=generateOrders(50)

	handler:=processOrders(gen)
	results:=filterOrders(handler,100)

	for result:=range results{
		fmt.Printf("id: %d \namount: %d \nstatus: %s\n\n",result.ID,result.Amount,result.Status)
	}
}