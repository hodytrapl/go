package taskstore

import (
    "encoding/json"
    "os"
    "sync"
)

//Task отвечает за поля задач
type Task struct{
    ID          int    `json:"id"`
    Title       string `json:"title"`
    Done    	bool   `json:"done"`
}

// TaskStore отвечает за хранение задач
type Storage struct {
    mu       sync.RWMutex // Мьютекс для защиты доступа
    filename string       // Имя файла базы данных
}

type TaskStore struct{
    mu sync.Mutex
    tasks map[int]Task
    nextID int
}

// NewStorage создает новое хранилище
func NewTaskStore(filename string) *TaskStore {
    return &TaskStore{
        tasks: make(map[int]Task),
        nextID:1,
    }
}

// NewStorage создает новое хранилище
func NewStorage(filename string) *Storage {
    return &Storage{filename: filename}
}

// SaveTasks сохраняет задачи из слайса Tasks в файл
func (s *Storage) SaveTasks(tasks []Task) error {
    s.mu.Lock()         // Блокируем на запись
    defer s.mu.Unlock() // Обязательно разблокируем на выходе

    data, err := json.MarshalIndent(tasks, "", "   ") // Indent для красоты
    if err != nil {
        return err
    }
    // 0644 - права доступа (rw-r--r--)
    return os.WriteFile(s.filename, data, 0644)
}
func (s *Storage) LoadTasks() ([]Task, error) {
    s.mu.RLock() // Блокируем ТОЛЬКО для чтения (другие тоже могут читать)
    defer s.mu.RUnlock()

    data, err := os.ReadFile(s.filename)
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

func (ts *TaskStore) Delete(id int) bool{
    ts.mu.Lock()
    defer ts.mu.Unlock()

    if _,ok:=ts.tasks(id);ok{
        delete(ts.tasks,id)
    }
}