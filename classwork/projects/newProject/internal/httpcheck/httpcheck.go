package httpcheck

import (
	"errors"
	"io"
	"net/http"
	"time"
)


func httpCheck(url string, timeout time.Duration) Result{
	/*
		провеяем ссылку и выносим ответ
	*/
	if(url==""){
		return nil, Errorf("ссылка не была объявлена")
	}

	var result Result
	result.URL = url

	start := time.Now()

    client := &http.Client{
        Timeout: timeout,
    }

	// "Упаковываем" наши данные

    resp, err := client.Get(url)
    if err != nil {
		//проверка на ошибки
		result.TTFB = time.Since(start)
		result.Finalize()

		result.Status = ""
		result.OK = ""
		result.SizeBytes = 0
		result.Error = err.Error()
		return result
    }
    defer resp.Body.Close()
    
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Обработка ошибки при выполнении запроса
		result.TTFB = time.Since(start)
		result.Finalize()

		result.Status = ""
		result.OK = ""
		result.SizeBytes = 0
		result.Error = err.Error()
		return result
	}

	// Заполняем результаты успешного запроса
	result.TTFB = time.Since(start)
	result.Finalize()

	result.Status = resp.Status
	result.OK = resp.Status
	result.SizeBytes = int64(len(body))
	result.Location = resp.Request.URL.String()

	return result
}