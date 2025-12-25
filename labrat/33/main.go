package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"sync"
)

type FileResult struct {
	Path string
	Size int64
}

type FileInfoProcessor interface {
	Analyze(path string, info fs.DirEntry) FileResult
}

type FIProcessor struct{}

func main() {
	var wg sync.WaitGroup

	processor := FIProcessor{}

	// Каналы для передачи результатов и ошибок
	results := make(chan FileResult)
	errors := make(chan error)

	// Срез для хранения всех результатов
	var allResults []FileResult

	// Семафор для ограничения количества одновременных горутин (буферизированный канал)
	semaphore := make(chan struct{}, 10)

	// Каталоги для сканирования
	directories := []string{
		"F:\\system files\\steamapps\\common\\tModLoader\\hodytrapl\\SaveData\\Mods",
		"F:\\system files\\steamapps\\common\\SteamVR",
		"F:\\system files\\steamapps\\common\\Terraria",
	}

	// Запускаем горутину для сбора результатов
	go func() {
		for result := range results {
			allResults = append(allResults, result)
		}
	}()

	// Обработка ошибок
	go func() {
		for err := range errors {
			fmt.Printf("Ошибка: %v\n", err)
		}
	}()

	// Запускаем обход каждого каталога
	for _, dir := range directories {
		wg.Add(1)
		go func(dirPath string) {
			defer wg.Done()

			// Захватываем слот семафора
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Освобождаем слот

			err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					errors <- fmt.Errorf("ошибка доступа к %s: %w", path, err)
					return nil
				}

				if !d.IsDir() {
					// Анализируем файл и отправляем результат
					result := processor.Analyze(path, d)
					results <- result
				}

				return nil
			})

			if err != nil {
				errors <- fmt.Errorf("ошибка обхода директории %s: %w", dirPath, err)
			}
		}(dir)
	}

	// Ждем завершения всех горутин обхода
	wg.Wait()

	// Закрываем каналы
	close(results)
	close(errors)

	// Сортируем и выводим результаты
	if len(allResults) > 0 {
		SortFileResultsBySize(allResults)
	} else {
		fmt.Println("Файлы не найдены")
	}
}

func (fip FIProcessor) Analyze(path string, info fs.DirEntry) FileResult {
	fileinfo, err := info.Info()
	if err != nil {
		return FileResult{Path: path, Size: 0}
	}
	return FileResult{Path: path, Size: fileinfo.Size()}
}

func SortFileResultsBySize(results []FileResult) {
	if len(results) == 0 {
		fmt.Println("Нет данных для сортировки")
		return
	}

	// Создаем копию для сортировки
	resultsSorted := make([]FileResult, len(results))
	copy(resultsSorted, results)

	// Сортируем по убыванию размера
	sort.Slice(resultsSorted, func(i, j int) bool {
		return resultsSorted[i].Size > resultsSorted[j].Size
	})

	topN := 5
	if len(resultsSorted) < topN {
		topN = len(resultsSorted)
	}

	fmt.Printf("TOP %d самых больших файлов:\n", topN)
	for i := 0; i < topN; i++ {
		fmt.Printf("%d место:\t путь: %s\t размер: %s\n",
			i+1,
			resultsSorted[i].Path,
			bytesToHumanReadable(resultsSorted[i].Size))
	}
}

func bytesToHumanReadable(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	var unitIndex int
	fsize := float64(size)

	for fsize >= 1024 && unitIndex < len(units)-1 {
		fsize /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", fsize, units[unitIndex])
}
