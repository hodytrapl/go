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
