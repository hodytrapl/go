package main

import (
    "encoding/json"
    "log"
    "net/http"
    "fmt"
    "strconv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

type Task struct{
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Done    	bool   `json:"done"`
}

var tasks = []Task{
	{ID:1, Title:"lee", Done:false},
	{ID:2, Title:"lggfgde", Done:true},
}

func adminOnly(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        if r.Header.Get("X-Admin-Key") !="secret123"{
            http.Error(w,"Forbidden",http.StatusForbidden)
            return
        }
        next.ServeHTTP(w,r)
    })
}

func getTaskByID(w http.ResponseWriter, r *http.Request){
	idStr:=chi.URLParam(r,"id")
    id,err:=strconv.Atoi(idStr)
    if err!=nil{
        http.Error(w,"invalid id", http.StatusBadRequest)
        return
    }

	for _, t :=range tasks{
		if t.ID == id{
			w.Header().Set("Content-Type","application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}

	http.Error(w, "status not found ", http.StatusNotFound)
}

func getAllTasks(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(tasks)
}

func deleteTaskByID(w http.ResponseWriter, r *http.Request){
	idStr:=chi.URLParam(r,"id")
    id,err:=strconv.Atoi(idStr)
    if err!=nil{
        http.Error(w,"invalid id", http.StatusBadRequest)
        return
    }

	index := -1
	for i, t := range tasks {
		if t.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request){
    var task Task
    err:=json.NewDecoder(r.Body).Decode(&task)
    if err!=nil{
        http.Error(w,"ttt json:"+err.Error(),http.StatusBadRequest)
        return
    }
    if task.Title==""{
        http.Error(w,"title is req", http.StatusBadRequest)
        return
    }

    maxID := 0
    for _, t := range tasks {
        if t.ID > maxID {
            maxID = t.ID
        }
    }
    task.ID = maxID + 1
    tasks = append(tasks, task)

    fmt.Printf("задача получена:%+v\n",task)

    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusCreated)

    if err:=json.NewEncoder(w).Encode(task);err!=nil{
        fmt.Printf("error encodind response: %v",err)
        return
    }
}

func putTaskByID(w http.ResponseWriter, r *http.Request){
    idStr:=chi.URLParam(r,"id")
    id,err:=strconv.Atoi(idStr)
    if err!=nil{
        http.Error(w,"invalid id", http.StatusBadRequest)
        return
    }

    var incoming Task
    if err:=json.NewDecoder(r.Body).Decode(&incoming);err!=nil{
        http.Error(w,err.Error(),http.StatusBadRequest)
        return
    }

    for i,t :=range tasks{
        if t.ID==id{
            tasks[i].Title =incoming.Title
            tasks[i].Done=incoming.Done
            
            w.Header().Set("Content-Type", "application/json")
            _ =json.NewEncoder(w).Encode(tasks[i])
            return
        }
    }

    http.Error(w, "task not found",http.StatusNotFound)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Service is running"))
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
    info := map[string]string{
        "author":  "Ваше имя",
        "version": "1.0",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    json.NewEncoder(w).Encode(info)
}

func main() {

    r:=chi.NewRouter()
    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)

    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/health", healthHandler)
        r.Get("/info", infoHandler)
        r.Post("/tasks", CreateTaskHandler)
        r.Get("/tasks", getAllTasks)
        r.Get("/tasks/{id}", getTaskByID)
        r.Put("/tasks/{id}", putTaskByID)
        r.With(adminOnly).Delete("/tasks/{id}", deleteTaskByID)
    })
	
    log.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}