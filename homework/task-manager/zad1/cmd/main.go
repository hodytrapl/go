package main

import (
    "log"
    "net/http"
    "sync"
    "os"
    "github.com/go-chi/chi/v5"
    "zad1/internal/middleware"
    "zad1/internal/taskstore"
    
)

var (
    store *tasks.TaskStore
)

func main() {
    store = tasks.NewTaskStore("tasks.json")

    r:=chi.NewRouter()
    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)

    r.Route("/api/v1/tasks", func(r chi.Router) {
        r.Post("/", tasks.CreateTaskHandler(store))
        r.Get("/", tasks.GetAllTasksHandler(store))
        r.Get("/{id}", tasks.GetTaskByIDHandler(store))
        r.Put("/{id}", tasks.PutTaskByIDHandler(store))
        r.With(middleware.BasicAuthMiddleware).Delete("/{id}",  tasks.DeleteTaskByIDHandler(store))
    })
	
    log.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}