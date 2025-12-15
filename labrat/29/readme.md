# Лабораторная работа sync.Mutex и конкурентный доступ к данным.
Тема: Создание потокобезопасного банковского счета.
## Задание 1. Базовая конкурентность (Обязательно)
1.	Создайте файл bank.go.
2.	Объявите структуру BankAccount с полем Balance int.
3.	Напишите функцию main, где создается счет с балансом 0.
4.	Запустите цикл на 1000 итераций. В каждой итерации запускайте горутину, которая:
	-	Добавляет к балансу 1 (Deposit).
	-	Используйте sync.WaitGroup, чтобы дождаться всех.
5.	Запустите программу. Убедитесь, что баланс в конце НЕ равен 1000 (или равен случайно, если повезет).
6.	Добавьте в структуру sync.Mutex.
7.	Реализуйте метод Deposit(amount int), который безопасно меняет баланс, используя Lock и Unlock.
8.	Измените main так, чтобы горутины вызывали этот метод.
9.	Критерий успеха: Программа всегда выводит ровно 1000.
## Задание 2. Потокобезопасный справочник (Обязательно)
1.	Создайте структуру PhoneBook, которая содержит внутри map[string]string (имя -> телефон) и sync.RWMutex.
2.	Реализуйте метод Set(name, phone string), использующий Lock.
3.	Реализуйте метод Get(name string) string, использующий RLock.
4.	В main создайте книгу.
5.	Запустите 1 горутину, которая в бесконечном цикле (или 100 раз) пишет новые телефоны.
6.	Запустите 5 горутин, которые параллельно читают телефоны.
7.	Критерий успеха: Программа не падает с ошибкой "fatal error: concurrent map read and map write". (Без мьютексов Go аварийно завершит программу при конкурентном доступе к карте).
## Задание 3. "Со звёздочкой" (Дополнительно)
Реализуйте метод Transfer(to *BankAccount, amount int) для перевода денег между счетами из Задания 1.
-	Вам нужно заблокировать оба счета (и с которого списываем, и на который зачисляем).
-	Подумайте: В каком порядке блокировать мьютексы, если А переводит Б, а Б переводит А одновременно? (Это может вызвать Deadlock, но пока просто попробуйте реализовать блокировку двух мьютексов).



# code
### task1

```
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

```

### task2

```
package main

import (
	"fmt"
	"sync"
	"math/rand"
	"strconv"
	"time"
)

type PhoneBook struct{
	mu 		sync.RWMutex
	wg      sync.WaitGroup
	contactList map[string]string 
}

func (pb *PhoneBook) Set(name, phone string){
	pb.mu.Lock()
	defer pb.mu.Unlock()
	fmt.Println("vnosim: ",name)
	pb.contactList[name]=phone
}

func (pb *PhoneBook) Get(name string) string{
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	return pb.contactList[name]
}

func generateName(id int) string {
	firstNames := []string{"Алексей", "Мария", "Иван", "Ольга", "Дмитрий", "Анна", "Сергей", "Елена", "Андрей", "Наталья"}
	lastNames := []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Васильев", "Павлов", "Семенов", "Голубев"}
	
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	
	return fmt.Sprintf("%s %s_%d", firstName, lastName, id)
}

func generatePhoneNumber(regionCode string) string {
	number := ""
	for i := 0; i < 10; i++ {
		number += strconv.Itoa(rand.Intn(10))
	}
	return regionCode + " " + number
}

func main() {
	rand.Seed(time.Now().UnixNano())
	book := PhoneBook{contactList: make(map[string]string)};

	
	for i:=0;i<1000;i++{
		book.wg.Add(1)
		go func(id int){
			defer book.wg.Done()
			name := generateName(id)
			number := generatePhoneNumber("+7")
			book.Set(name,number)
		}(i)
	}

	book.wg.Wait()

	count := 0
	for name, phone := range book.contactList {
		fmt.Printf("%s: %s\n", name, phone)
		count++
		if count >= 10 {
			break
		}
	}
}

```


### task3

Вариант 1

```

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


```