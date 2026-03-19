package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "fmt"
)

type Task struct{
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost{
        http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
        return
    }

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

    task.ID=123
    task.CreatedAt = time.Now()
    task.Status="created"

    fmt.Printf("задача получена:%+v\n",task)

    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusCreated)

    if err:=json.NewEncoder(w).Encode(task);err!=nil{
        fmt.Printf("error encodind response: %v",err)
        return
    }

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

    mux:=http.NewServeMux()

    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/info", infoHandler)
    mux.HandleFunc("/tasks", CreateTaskHandler)
    mux.HandleFunc("POST /tasks", CreateTaskHandler)
    mux.HandleFunc("GET /tasks", CreateTaskHandler)
    // mux.HandleFunc("GET /tasks/{id}", CreateTaskHandler)
    // idStr:= r.PathValue("id")
	
    log.Println("Server is running on http://localhost:8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatal(err)
    }
}