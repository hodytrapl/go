package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// FileProcessor - интерфейс для обработки одного файла
type FileProcessor interface {
	Process(path string, entry fs.DirEntry) (int64, error)
}

// SizeCounter - реализация интерфейса, считающая размер
type SizeCounter struct{}

func (cs SizeCounter) Process(path string, entry fs.DirEntry) (int64, error) {
	info, err := entry.Info()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func walkConcurrentWithSem(path string, processor FileProcessor, sizes chan<- int64, wg *sync.WaitGroup) {
	// ВАЖНО!!! Мы уменьшаем счетчик wg, когда выходим из функции обхода ЭТОЙ папки
	// Но WaitGroup у нас будет глобальный для всех файлов
	// Лучше сделаем проще: walk  запускает горутины только для файлов

	entries, err := os.ReadDir(path)
	if err != nil {
		// Пропускаем ошибки доступа для простоты примера
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			// Рекурсивно идем дальше в том же потоке (или тоже в горутине, но это сложнее)
			walkConcurrentWithSem(fullPath, processor, sizes, wg)
		} else {
			// Для файла запускаем горутину
			wg.Add(1)
			go func(p string, e fs.DirEntry) {
				defer wg.Done()
				sem <- struct{}{}        // Заняли слот
				defer func() { <-sem }() // Освободили слот

				size, _ := processor.Process(p, e)
				sizes <- size
				// Отправляем результат в канал
			}(fullPath, entry)
		}

	}
}

var sem = make(chan struct{}, 20) //Ограничиваем до 20 одновременных чтений

func main() {
	path := "D:/SteamLibrary" // Или os.Args[1]
	processor := SizeCounter{}

	sizes := make(chan int64)
	var wg sync.WaitGroup

	//Запускаем обход в отдельной горутине, чтобы не блокировать main
	go func() {
		walkConcurrentWithSem(path, processor, sizes, &wg)
		wg.Wait()
		close(sizes) // Закрываем канал, когда все закончили
	}()

	var totalSize int64
	for size := range sizes {
		totalSize += size
	}

	fmt.Printf("Total size: %2f MB\n", float64(totalSize/1024/1024))
}
