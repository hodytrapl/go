package main

import ( 
	"fmt"
	"sync"
)

type BankAccount struct{
	mu 		sync.Mutex

	id int
	balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.balance += amount

}
func (from *BankAccount) Transfer(to *BankAccount, amount int) error {
	if(from==to){
		return fmt.Errorf("низяшечки перевод самому себе!")
	}

	if from.balance < amount {
		return fmt.Errorf("дэнэг мало!")
	}

	first,second= from,to
	if (from.id>to.id ){
		first,second= to,from
	}

	first.mu.Lock()
	second.mu.Lock()
	defer first.mu.Unlock()
	defer second.mu.Unlock()

	first.balance-=amount
	second.balance+=amount

	return nil
}


func main() {
	user := BankAccount{id:0 balance:0}
	user2 := BankAccount{id:1 balance:1500}

	var wg sync.WaitGroup

	for i:=0;i<1000;i++ {
		wg.Add(1)
		go func(value int){
			defer wg.Done()
			user.Deposit(value)
		}(1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		user.Transfer(&user2, 455)
	}()

	wg.Wait()
	
	fmt.Printf("Баланс после 1000 пополнений: %d\n", user.balance)
	fmt.Printf("Баланс после 1000 пополнений: %d\n", user2.balance)
}
