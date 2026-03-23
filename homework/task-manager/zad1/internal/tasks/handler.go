package tasks

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"zad1/internal/taskstore"
)

type Handler struct{
	store *taskstore.TaskStore
}

func newHandler(store *taskstore.TaskStore) *Handler{
	return &Handler{store:store}
}

func getTaskByID(w http.ResponseWriter, r *http.Request){
	idStr:=chi.URLParam(r,"id")
    id,err:=strconv.Atoi(idStr)
    if err!=nil{
        http.Error(w,"invalid id", http.StatusBadRequest)
        return
    }

    tasksMu.RLock() // потокобезопасное чтение
    defer tasksMu.RUnlock()


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
    tasksMu.RLock()
    defer tasksMu.RUnlock()

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(tasks)

}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request){
	idStr:=r.PathValue("id")
	id,_=strconv.Atoi(idStr)

	if h.store.Delete(id){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task deleted"))

	} else {
		http.Error(w, "Task not found", http.StatusNotFound)

	}
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request){
    var task Task
    
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if task.Title == "" {
        http.Error(w, "Title is required", http.StatusBadRequest)
        return
    }

    tasksMu.Lock()
    defer tasksMu.Unlock()

    task.ID = nextID
    nextID++

    tasks = append(tasks, task)

    if err := store.SaveTasks(tasks); err != nil {
        http.Error(w, "Failed to save tasks", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(task)


}

func putTaskByID(w http.ResponseWriter, r *http.Request){
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var incoming Task
    if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    tasksMu.Lock()
    defer tasksMu.Unlock()

    for i := range tasks {
        if tasks[i].ID == id {
            tasks[i].Title = incoming.Title
            tasks[i].Done = incoming.Done

            if err := store.SaveTasks(tasks); err != nil {
                http.Error(w, "Failed to save tasks", http.StatusInternalServerError)
                return
            }

            w.Header().Set("Content-Type", "application/json")
            _ = json.NewEncoder(w).Encode(tasks[i])
            return
        }
    }
    http.Error(w, "Task not found", http.StatusNotFound)

}

func calcNextID(ts []Task) int {
    maxID := 0
    for _, t := range ts {
        if t.ID > maxID {
            maxID = t.ID
        }
    }
    return maxID + 1
}