# Домашнее задание
#### Тема: 
Обработка данных с использованием Worker Pool.
#### Цель: 
Закрепить навык создания конвейера: Данные -> Очередь -> Воркеры -> Результат.
#### Легенда:
Вы разрабатываете систему обработки изображений (условно). Вам поступает список размеров сторон квадратов, нужно посчитать их площади и периметры параллельно, так как вычислений очень много.
## Задание:
Создайте файл homework_pool.go.
1. Входные данные:
В main создайте слайс чисел: inputs := []int{1, 5, 12, 5, 3, 8, 9}.
2. Задача (Job):
Структура Job должна содержать число из слайса.
3. Воркеры:
Запустите 3 воркера. Каждый воркер:
    - Принимает число.
    -  Считает квадрат числа (площадь).
    - Отправляет результат в канал результатов.
    - Бонус: добавляет time.Sleep (случайное время), чтобы сбить порядок вывода.
4. Результат:
В main вы должны принять все ответы из канала результатов и вывести их в консоль.
Критерии приемки (минимум на зачет):
Программа компилируется (go run ...).
Используется sync.WaitGroup для контроля завершения воркеров.
Все числа из входного слайса обработаны.
Программа корректно завершается (нет ошибки deadlock).
Задание со звёздочкой (по желанию):a
Попробуйте реализовать это так, чтобы main не знал заранее, сколько будет задач.
### Схема:
go func() { ... отправка в jobs; close(jobs) }() — отдельная горутина-генератор.
go func() { wg.Wait(); close(results) }() — отдельная горутина-наблюдатель.
В main только цикл range по каналу results.
Удачи! Это основа для нашего будущего проекта загрузчика сайтов.


# code

```
package main

import (
	"fmt"
	"sync"
)

type Job struct{
	number int
}

type Result struct {
	job Job
	square int
}

func main(){
	var wg sync.WaitGroup
	inputs := []int{1, 5, 12, 5, 3, 8, 9}

	jobs := make(chan Job,len(inputs))
	results := make(chan Result, len(inputs))

	numWorkers := 3

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup){
			defer wg.Done()
			fmt.Printf("Воркер %d начал работу\n", id)
			for job  := range jobs{
				result:=job.number*job.number
				results<-Result{job: job, square: result}
			}
			
			fmt.Printf("Воркер %d завершил работу\n", id)
		}(i,jobs,results, &wg)
	}

	for _, num := range inputs {
		jobs <- Job{number: num}
	}

	close(jobs)

	wg.Wait()
	
	close(results)

	fmt.Println("\nРезультаты:")
	for result := range results {
		fmt.Printf("Число: %d, Квадрат: %d\n", result.job.number, result.square)
	}

}

```