//zad1/internal/middleware/middleware.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func LoggingMiddleware(hand http.Handler) http.Handler{
    /*функция сохраняет время обработки другой  функции*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        hand.ServeHTTP(w, r)
        log.Printf("[%s | %s] - %v", r.Method, r.URL, time.Since(start))
    })
}

func BasicAuthMiddleware(hand http.Handler) http.Handler {
    /*функция хранит все авторизацию юзера*/
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        name, pass, ok := r.BasicAuth()
        if !ok || name != "admin" || pass != "secret" {
            w.Header().Set("WWW-Authenticate", `Basic realm = "Restricted"`)
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        hand.ServeHTTP(w, r)
    })
}

func JSONHeaderMiddleware(hand http.Handler) http.Handler {
    /*функция переделывает сайт в json формат, а не text*/
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        hand.ServeHTTP(w, r)
    })
}
