//zad1/internal/tasks/handlers.go
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

func NewHandler(store *taskstore.TaskStore) *Handler{
	return &Handler{store:store}
}

/*реализовать crud*/

func (h *Handler) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
    tasks := h.store.GetAllTasks()
    json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    task, found := h.store.GetTaskByID(id)
    if !found {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(task)
}

func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task taskstore.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    createdTask := h.store.Create(task)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) PutTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    var task taskstore.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if success := h.store.Update(id, task); !success {
        http.Error(w, "Task not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(task)
}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request){
	idStr := chi.URLParam(r, "id")
    id, _ := strconv.Atoi(idStr)

	if h.store.Delete(id){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task deleted"))

	} else {
		http.Error(w, "Task not found", http.StatusNotFound)

	}
}