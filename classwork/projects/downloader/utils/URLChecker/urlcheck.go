package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"text/tabwriter"
	"encoding/json"
	"strconv"
)

type Data struct {
	/*
		данные ссылки
	*/
	URL 		string 	`json:"URL"` 		 //ссылка
	HTTPStatus 	string 	`json:"HTTPStatus"`	 //статус код
	Status  	string 	`json:"Status"`		 //просто статус
	Ttfb_ms 	string 	`json:"Ttfb_ms"`	 //время до первого байта 
	Size_bytes 	string 	`json:"Size_bytes"`	 //количество прочитанных байт 
	isSubstring	string 	`json:"isSubstring"` //найдена ли подстрока
	Error	 	string 	`json:"Error"`		 //и текст ошибки (если есть)
}

func ReadUrls(filename string) ([]string,error){
	/*
		читаем файл со ссылками и возращаем их список
		делаем еще проверку есть ли схема https://
	*/
	file, err := os.Open(filename)

	if(err!=nil){
		return nil, err
	}
	defer file.Close()
	
	var urls []string

	//добавляем схему если её нету
	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        url := strings.TrimSpace(scanner.Text())
        if url != "" {
			if url[0:8]!="https://" && url[0:7]!="http://"{
				url="https://"+url
			}
            urls = append(urls, url)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

	return urls ,nil
}


func checkURL(
	url string, timeout time.Duration,
	substr string,
	results chan<- Data,
	) {
	/*
		проверяем ссылку на работоспособность и выводим её готовый пакет структуры для таблицы и подобного
	*/

    start := time.Now()

    client := &http.Client{
        Timeout: timeout,
    }
	// "Упаковываем" наши ссылки в горутину, для вывода таблиц
    resp, err := client.Get(url)
    if err != nil {
		elapsed := time.Since(start)
		results <-Data{
		URL:			url,
		HTTPStatus:		"",
		Status:			"",
		Ttfb_ms:		fmt.Sprintf("%d", elapsed.Milliseconds()),
		Size_bytes:		fmt.Sprintf("%d", 0),
		isSubstring:	fmt.Sprintf("%t", false),
		Error:			err.Error(),
		}

        
        return
    }
    defer resp.Body.Close()
    
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		elapsed := time.Since(start)
		results <- Data{
			URL:   url,
			HTTPStatus:		fmt.Sprintf("%d", resp.StatusCode),
			Status:			resp.Status,
			Ttfb_ms:		fmt.Sprintf("%d", elapsed.Milliseconds()),
			Size_bytes:		fmt.Sprintf("%d", len(body)),
			isSubstring:	fmt.Sprintf("%t", false),
			Error: err.Error(),
		}
		return
	}

	found := false
	if substr != "" {
		found = strings.Contains(string(body), substr)
	}

	elapsed := time.Since(start)

	data :=Data{
		URL:			url,
		HTTPStatus:		fmt.Sprintf("%d", resp.StatusCode),
		Status:			resp.Status,
		Ttfb_ms:		fmt.Sprintf("%d", elapsed.Milliseconds()),
		Size_bytes:		fmt.Sprintf("%d", len(body)),
		isSubstring:	fmt.Sprintf("%t", found),
		Error:			fmt.Sprintf("%v", nil),
	}

    results <- data
}

func SaveToJson(results []Data, jsonFLag string) error{
	/*
		сохраняем наши файлы в json документе
	*/
	file, err := os.Create(jsonFLag+".json")
	if err!=nil{
		return err;
	}
	defer file.Close()
//проверяем маленькие ссылки, поэтому без маршала

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func PrintTable(results []Data){
	/*
		Создаёт ту самую таблицу и выводит ее в cmd
	*/
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w,
		"URL\tHTTP\tSTATUS\tTTFB(ms)\tSIZE(bytes)\tSUBSTRING\tERROR",
	)

	for _, r := range results {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.URL,
			r.HTTPStatus,
			r.Status,
			r.Ttfb_ms,
			r.Size_bytes,
			r.isSubstring,
			r.Error,
		)
	}

	w.Flush()
}

func main(){
	/*ядро*/
	timeoutStr := flag.String("t", "10s", "Timeout for HTTP requests")
	json := flag.String("json", "results", "save files in json,")
	contains := flag.String("contains", "", "Find substrings")
	workers := flag.Int("workers", 10, "Find substrings")	
    flag.Parse()

	var timeout time.Duration

	urls, err := ReadUrls("urls.txt")
    if err != nil {
        fmt.Printf("Error read URLs: %v\n", err)
        os.Exit(1)
    }

	if len(urls) == 0 {
        fmt.Println("No URLs found in urls.txt.")
        return
    }

	//обработчик ошибок на ввод времени
	if ms, err := strconv.Atoi(*timeoutStr); err == nil {
		timeout = time.Duration(ms) * time.Millisecond
	} else {
		timeout, err = time.ParseDuration(*timeoutStr)
		if err != nil {
			fmt.Printf("Invalid timeout value\n")
			os.Exit(1)
		}
	}

	fmt.Printf("%d URLs timeout %v...\n\n", len(urls), timeout)

	//создаем потоки для обработки ссылок, взависимости от флага workers
	jobs := make(chan string)
    results := make(chan Data, len(urls))
	var wg sync.WaitGroup

	for i:=0;i<*workers;i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			for url := range jobs {		
				checkURL(
					url,
					timeout,
					*contains,
					results,
				)
			}
		}()
	}

	//заносим все ссылки в горутину
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	//ждём ответа от всех и закрываем канал
    go func() {
		wg.Wait()
        close(results)
    }()

	//обработка таблицы и вывод всего в json
	var allResults []Data

	for r := range results {
		allResults = append(allResults, r)
	}
	PrintTable(allResults)

	SaveToJson(allResults,*json)
}