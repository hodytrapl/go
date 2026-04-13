//zad1/internal/taskstore/storage.go
package taskstore

import (
    "encoding/json"
    "os"
    "sync"
    "strings"
)

func (ts *TaskStore) saveToFile() {
    // собираем все задачи в слайс
    tasksSlice := make([]Task, 0, len(ts.tasks))
    for _, t := range ts.tasks {
        tasksSlice = append(tasksSlice, t)
    }
    _ = ts.SaveTasks(tasksSlice)
}

func (ts *TaskStore) LoadTasks() ([]Task, error) {
    ts .mu.RLock() // Блокируем ТОЛЬКО для чтения (другие тоже могут читать)
    defer ts .mu.RUnlock()

    data, err := os.ReadFile(ts.filename)
    if err != nil {
        if os.IsNotExist(err) {
            // // Если файла нет, это не ошибка, возвращаем пустой список
            return []Task{}, nil
        }
        return nil, err
    }

    // Пустой файл — это не ошибка, просто нет задач
    if len(strings.TrimSpace(string(data))) == 0 {
        return []Task{}, nil
    }

    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }
    return tasks, nil
}
/*реализовать crud*/

func (ts *TaskStore) GetAllTasks() []Task {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    tasks := make([]Task, 0, len(ts.tasks))
    for _, task := range ts.tasks {
        tasks = append(tasks, task)
    }
    return tasks
}

func (ts *TaskStore) GetTaskByID(id int) (Task, bool) {
    ts.mu.Lock()
    defer ts.mu.Unlock()
    task, ok := ts.tasks[id]
    return task, ok
}

func (ts *TaskStore) Create(task Task) Task {
    ts.mu.Lock()
    defer ts.mu.Unlock()

    task.ID = ts.nextID
    ts.nextID++
    ts.tasks[task.ID] = task
    ts.saveToFile()
    return task
}

func (ts *TaskStore) Update(id int, task Task) bool {
    ts.mu.Lock()
    defer ts.mu.Unlock()

    if _, ok := ts.tasks[id]; !ok {
        return false
    }
    task.ID = id
    ts.tasks[id] = task
    ts.saveToFile()
    return true
}

func (ts *TaskStore) Delete(id int) bool{
    ts.mu.Lock()
    defer ts.mu.Unlock()

    if _,ok:=ts.tasks[id];ok{
        delete(ts.tasks,id)
        return true
    }
    return false
}