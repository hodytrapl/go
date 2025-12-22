# Паттерны Fan-out, Fan-in и Pipeline (Конвейер). Лабораторная работа.

#### Тема: 
Построение конвейера обработки данных.
#### Задание: 
Вам нужно реализовать систему обработки "заказов". Заказ — это просто число (ID). Система должна состоять из 3 стадий.
### Шаг 1. Подготовка данных (Структура и Интерфейс) 
Создайте файл main.go. Определите структуру Order и интерфейс. (Использование интерфейса — задание со звездочкой, если сложно — используйте просто struct, но попробуйте с интерфейсом).
```
type Order struct {
    ID     int
    Amount int    // сумма заказа
    Status string // "new", "processed", "paid"
}
// Задание со звездочкой: создайте интерфейс Processable
// type Processable interface { Process() }
```

### Шаг 2. Стадия 1: Генератор (Generator)
Напишите функцию generateOrders(count int) <-chan Order.
- Она должна создавать count заказов со статусом "new" и случайной суммой.
- Отправлять их в канал.
- Закрывать канал после завершения.
### Шаг 3. Стадия 2: Обработчик (Processor) - Fan-out
Напишите функцию processOrders(in <-chan Order) <-chan Order.
- Функция должна запускать несколько горутин (внутри себя или вызывать её несколько раз в main — на ваш выбор), чтобы имитировать параллельную обработку.
- Логика обработки: поменять статус на "processed".
- Подсказка: Если вы запускаете горутины внутри функции, не забудьте про sync.WaitGroup для закрытия выходного канала.
### Шаг 4. Стадия 3: Фильтр (Filter)
Напишите функцию filterOrders(in <-chan Order, minAmount int) <-chan Order.
- Она читает заказы.
- Пропускает дальше только те, у которых Amount > minAmount.
- Остальные отбрасывает.
### Шаг 5. Сборка (Main) В функции main:
1.	Запустите генератор (50 заказов).
2.	Передайте канал в обработчик.
3.	Передайте канал фильтру (сумма > 100).
4.	Выведите итоговые заказы в консоль.
## Критерии успеха:
1.	Программа не падает с deadlock.
2.	В конце выводятся только заказы с суммой > 100 и статусом "processed".
3.	Используются каналы для передачи данных.


# code
```
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
```

