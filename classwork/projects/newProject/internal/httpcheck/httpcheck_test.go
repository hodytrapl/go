package httpcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpCheck_Success(t *testing.T) {
	// Создаем локальный тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	result := HttpCheck(ts.URL, 2*time.Second)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}
	if result.SizeBytes != 11 {
		t.Errorf("Expected size 11, got: %d", result.SizeBytes)
	}
	if result.Status != "200 OK" {
		t.Errorf("Expected status '200 OK', got: %s", result.Status)
	}
	if result.TTFB <= 0 {
		t.Errorf("TTFB should be greater than 0")
	}
}

func TestHttpCheck_EmptyURL(t *testing.T) {
	result := HttpCheck("", 2*time.Second)
	if result.Error != "ссылка не была объявлена" {
		t.Errorf("Expected error for empty URL, got: %s", result.Error)
	}
}

func TestHttpCheck_InvalidURL(t *testing.T) {
	// Некорректный протокол
	result := HttpCheck("http://[::1]:NamedPort", 2*time.Second)
	if result.Error == "" {
		t.Errorf("Expected error for invalid URL, but got none")
	}
}
