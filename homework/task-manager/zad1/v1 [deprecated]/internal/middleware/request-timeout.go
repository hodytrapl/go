package middleware

import (
	"context",
	"time"
)

func RequestTimeoutMiddleware(d time.Duration) func(http.Handler) http.handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			ctx,cancel :=context.WithTimeout(r.Context(),d)
			defer cancel()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}