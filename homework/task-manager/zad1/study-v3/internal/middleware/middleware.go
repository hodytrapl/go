package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging логирует метод, путь и время выполнения запроса
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Логируем после выполнения
		log.Printf("%s %s served in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// BasicAuth проверяет наличие прав администратора
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != "admin" || pass != "secret" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONHeaderMiddleware выставляет заголовок Content-Type для JSON-ответов
func JSONHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
