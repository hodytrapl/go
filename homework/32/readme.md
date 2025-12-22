# Тема: Конвейер обработки логов (имитация).
#### Цель:
Закрепить навыки построения Pipeline, Fan-out и Fan-in на задаче, похожей на реальную.
#### Задание: 
Напишите программу, которая имитирует чтение файлов, их фильтрацию и подсчет данных. Вам понадобятся пакеты fmt, strings, time, sync.
#### Требования к реализации:
1.	Источник (Source): Функция, которая принимает срез строк (имена файлов) и выдает их в канал по одному.

	- Пример списка: ["data.txt", "image.png", "info.txt", "backup.zip"].
2.	Стадия 1 (Filter): Функция, которая читает имена файлов из канала, проверяет расширение. Если файл оканчивается на .txt, передает его дальше. Остальные игнорирует.

3.	Стадия 2 (Processing - Fan-out): Функция, которая принимает имя файла и возвращает "результат обработки" (структуру с именем файла и числом — размером или количеством строк).

	- Запустите 3 воркера (горутины), выполняющих эту функцию.
	- Внутри функции добавьте time.Sleep (случайно от 100мс до 500мс), чтобы имитировать чтение с диска.
	- Возвращайте случайное число в качестве результата.
4.	Сборщик (Fan-in): Объедините результаты от воркеров в один канал и выведите общую сумму (например, "Всего обработано строк: 1540").

#### Что нужно сдать:
- Ссылка на Git или файл .go.
- Небольшой комментарий в коде: где именно у вас Fan-out, а где Fan-in.
#### Критерии приемки:
- go run main.go выполняется без ошибок и дедлоков.
- Программа завершается сама после обработки всех файлов.
- Видна работа воркеров (можно добавить принты "Воркер 1 начал файл data.txt").

> [!WARNING]
> не решил прорблему deadlock

# code 

```
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"sync"
)

type Datafile struct{
	nameFile string;
	number int;
}

func Source(files []string,channel chan<-string, wg *sync.WaitGroup){
	defer wg.Done()
	for _, file := range files {
		channel <- file
	}
	close(channel)
}

func Filter(in <-chan string, out chan<- string, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(out)
	
	var substr = ".txt"

	for file :=range in{
		if(strings.Contains(file, substr)){
			out<-file
		}
	}
}

func Processing(in <-chan string, out chan<- Datafile, wg *sync.WaitGroup){
	defer wg.Done()
	for file := range in {
		sleepTime := time.Duration(rand.Intn(401)+100) * time.Millisecond
		time.Sleep(sleepTime)
		
		result := Datafile{
			nameFile: file,
			number:   rand.Intn(1000) + 1,
		}
		out <- result
	}
}

func FanIn(in <-chan Datafile,out chan<- int, wg *sync.WaitGroup){
	defer wg.Done()
	defer close(out)
	total:=0
	for file :=range in{
		total+=file.number
	}
	out <- total
}

func main(){
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	str:= []string{"data.txt", "image.png", "info.txt", "backup.zip"}

	sourceChan := make(chan string, 10)
	filterChan := make(chan string, 10)
	resultChan := make(chan Datafile, 10)
	resultTotalChan := make(chan int, 1)

	wg.Add(1)
	go Source(str,sourceChan,&wg)

	wg.Add(1)
	go Filter(sourceChan,filterChan,&wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go Processing(filterChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	wg.Add(1)
	go FanIn(resultChan, resultTotalChan, &wg)

	total := <-resultTotalChan
	fmt.Println("всего обработано строк:", total)

	wg.Wait()
}
```