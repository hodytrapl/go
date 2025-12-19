# Домашнее задание
Тема: Конкурентный подсчет слов в предложениях.
### Зачем это нужно:
В реальных задачах часто нужно агрегировать статистику из множества источников. Мы закрепим работу с Mutex для защиты общей map.
Задача:
У вас есть срез строк (предложений). Нужно посчитать общее количество каждого слова во всех предложениях.
Например: ["hello world", "hello go"] -> {"hello": 2, "world": 1, "go": 1}.
### Требования:
1. Создайте файл homework.go.
2. Напишите функцию ConcurrentWordCount(sentences []string) map[string]int.
3. Функция должна запускать по одной горутине на каждое предложение.
4. Все горутины должны писать результат в одну общую карту map[string]int.
5. Используйте sync.Mutex, чтобы защитить карту от гонки данных.
6. Используйте sync.WaitGroup, чтобы дождаться обработки всех предложений перед возвратом результата.
7. (Опционально) Попробуйте использовать strings.Fields(sentence) для разбиения строки на слова.

### Входные данные для теста:
```
text := []string{
    "quick brown fox",
    "lazy dog",
    "quick brown fox jumps",
    "jumps over lazy dog",
}
```

### Критерии приёмки:

    • Код компилируется и запускается.
    • Результат подсчета верен (можно сверить с ручным подсчетом).
    • В коде есть Lock и Unlock вокруг записи в карту.
    • Используется defer wg.Done().
    
### Мини-эссе (3-5 предложений):
Ответьте в комментариях к коду или в текст.файле на вопрос:
«Почему в этой задаче мы использовали Mutex, а не каналы? В каком случае каналы были бы удобнее?»
# Code
```
package main

import (
	"fmt"
	"strings"
	"sync"
)

func ConcurrentWordCount(sentences []string) map[string]int {
	result := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, sentence := range sentences {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			words := strings.Fields(s)

			mu.Lock()
			for _, word := range words {
				result[word]++
			}
			mu.Unlock()
		}(sentence)
	}

	wg.Wait()

	return result
}

func main() {
	text := []string{
		"quick brown fox",
		"lazy dog",
		"quick brown fox jumps",
		"jumps over lazy dog",
	}
	wordCount := ConcurrentWordCount(text)

	for word, count := range wordCount {
		fmt.Printf("%s: %d\n", word, count)
	}
}

```

# эссе
Почему в этой задаче мы использовали Mutex, а не каналы? В каком случае каналы были бы удобнее?

потому что так легче работать.