//zad1/cmd/main.go
package main

import (
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    "zad1/internal/middleware"
    "zad1/internal/tasks"
)

var (
    store *taskstore.TaskStore
)

func main() {
    store = taskstore.NewTaskStore("tasks.json")
    r:=chi.NewRouter()
    r.Use(middleware.JSONHeaderMiddleware)

    handler := tasks.NewHandler(store)
    r.Route("/api/v1/tasks", func(r chi.Router) {
        r.Get("/", handler.GetAllTasksHandler)
        r.Get("/{id}", handler.GetTaskByIDHandler)
        r.Post("/", handler.CreateTaskHandler)
        r.Put("/{id}", handler.PutTaskByIDHandler)
        r.With(middleware.BasicAuthMiddleware).Delete("/{id}",  handler.DeleteTaskByID)
    })
	
    log.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}