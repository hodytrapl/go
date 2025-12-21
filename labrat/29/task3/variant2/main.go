package main

import ( 
	"fmt"
	"sync"
)

type BankAccount struct{
	mu 		sync.Mutex
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
	first,second= from,to
	if (uintptr(unsafe.Pointer(first))>uintptr(unsafe.Pointer(second))){
		first,second= second,first
	}

	first.mu.Lock()
	second.mu.Lock()
	defer first.mu.Unlock()
	defer second.mu.Unlock()

	if from.balance < amount {
		return fmt.Errorf("дэнэг мало!")
	}

	first.balance-=amount
	second.balance+=amount

	return nil
}

func(b *BankAccount) getBalance() int{
	return b.balance
}


func main() {
	user := BankAccount{balance:0}
	user2 := BankAccount{balance:1500}

	var wg sync.WaitGroup

	for i:=0;i<1000;i++ {
		wg.Add(1)
		go func(value int){
			defer wg.Done()
			user.Deposit(value)
		}(1)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		user.Transfer(&user2, 455)
	}()

	go func() {
		defer wg.Done()
		if(err:=user2.Transfer(&user, 55);err!=nil){

		}
	}()

	wg.Wait()
	
	fmt.Printf("Баланс после 1000 пополнений: %d\n", user.getBalance())
	fmt.Printf("Баланс после 1000 пополнений: %d\n", user2.getBalance())
}
