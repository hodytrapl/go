package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Order struct {
	ID       int
	Customer string
	Amount   int
}

func ShopOrder(order Order, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Обработка заказа %d для %s...\n", order.ID, order.Customer)

	duration := time.Duration(rand.IntN(3)) * time.Second
	time.Sleep(duration)

	fmt.Printf("Заказ %d на сумму %d успешно оформлен\n", order.ID, order.Amount)
}

func main() {
	orders := []Order{
		{ID: 101, Customer: "Kostya", Amount: 2354},
		{ID: 102, Customer: "alfedov", Amount: 6475262},
		{ID: 103, Customer: "lololoska", Amount: 634},
		{ID: 104, Customer: "genadiy", Amount: 63426},
		{ID: 105, Customer: "erdjan", Amount: 345646},
	}

	var wg sync.WaitGroup

	fmt.Println("Запускаем систему обработки заказов")
	for _, order := range orders {
		wg.Add(1)
		go ShopOrder(order, &wg)
	}

	wg.Wait()

	fmt.Println("Все заказы успешно обработаны")
}
