package main

import ( 
	"fmt"
	"sync"
)

type BankAccount struct{
	mu 		sync.Mutex
	wg      sync.WaitGroup
	balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	defer b.wg.Done()
	b.balance += amount
}


func main() {
	user := BankAccount{balance:0}

	for i:=0;i<1000;i++ {
		user.wg.Add(1)
		go func(value int){
			user.Deposit(value)
		}(1)
	}

	user.wg.Wait()
	
	fmt.Printf("Баланс после 1000 пополнений: %d\n", user.balance)
}
