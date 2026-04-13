package tasks

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

// TaskStore отвечает за хранение задач в файле
type TaskStore struct {
	mu       sync.RWMutex // Мютекс для защиты доступа к файлу при I/O операциях
	fileName string
}

func NewTaskStore(filename string) *TaskStore {
	return &TaskStore{fileName: filename}
}

// SaveTasks
func(ts *TaskStore) SaveTasks(ctx context.Context, tasks []Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ts.mu.Lock()
	defer ts.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return os.WriteFile(ts.fileName, data, 0644)
}

// LoadTasks
func(ts *TaskStore) LoadTasks(ctx context.Context) (tasks []Task, error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Читаем файл
	if err := ctx.Err(); err != nil {
		return err
	}

	data, err := os.ReadFile(ts.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	// Если файла нет - это не ошибка, просто нет задач
	if len(strings.TrimSpace(string(data))) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	return tasks, nil
}

// Симуяция медленного ввода/вывода
func (ts *TaskStore) SimulateSlowIO(ctx context.Context, d time.Duration) err {
	if d <= 0 {
		return ctx.Err()
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C():
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}