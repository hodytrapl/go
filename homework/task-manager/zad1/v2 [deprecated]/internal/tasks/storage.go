package tasks

import (
	"sync"
)

type  TaskStore struct{
	mu 		 sync.RWMutex
	fileName string
}

// NewStorage создает новое хранилище
func NewTaskStore(filename string) *TaskStore {
    return &TaskStore{
        filename: filename
    }
}

// SaveTasks сохраняет задачи из слайса Tasks в файл
func (ts *TaskStore) SaveTasks(ctx context.Context, tasks []Task) error {
	if err:=ctx.Err();err!=nil{
		return err
	}
	
    ts.mu.Lock()         // Блокируем на запись
    defer ts.mu.Unlock() // Обязательно разблокируем на выходе
	
    data, err := json.MarshalIndent(tasks, "", "   ") // Indent для красоты
    if err != nil {
        return err
    }

	if err:=ctx.Err();err!=nil{
		return err
	}

    // 0644 - права доступа (rw-r--r--)
    return os.WriteFile(ts.filename, data, 0644)
}

func (ts *TaskStore) LoadTasks(ctx context.Context) (tasks []Task, error) {
    if err:=ctx.Err();err!=nil{
		return err
	}
	
    ts.mu.RLock()         // Блокируем на запись
    defer ts.mu.RUnlock() // Обязательно разблокируем на выходе

	if err:=ctx.Err();err!=nil{
		return err
	}


    data, err := os.ReadFile(ts.filename)
    if err != nil {
        if os.IsNotExist(err) {
            // // Если файла нет, это не ошибка, возвращаем пустой список
            return []Task{}, nil
        }
        return nil, err
    }

	if err:=ctx.Err();err!=nil{
		return err
	}

    // Пустой файл — это не ошибка, просто нет задач
    if len(strings.TrimSpace(string(data))) == 0 {
        return []Task{}, nil
    }

    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }

	if err:=ctx.Err();err!=nil{
		return err
	}

    return tasks, nil
}

func (ts *TaskStore) SimulateSlowID(ctx context.Context, d time.Duration)err{
	if d<=0{
		return ctx.Err()
	}

	if err:=ctx.Err();err!=nil{
		return err
	}

	timer:=time.NewTimer(d)
	defer timer.Stop()

	select{
	case<-timer.C():
		return ctx.Err()
		case<-timer.Done():
		return ctx.Err()
	}
}