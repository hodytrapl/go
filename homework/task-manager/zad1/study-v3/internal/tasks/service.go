package tasks

import(
	"context"
	"sync"
)

type Service struct{
	store *TaskStore

	mu sync.RWMutex
	tasks []Task
	nextID int
}

func NewService(ctx context.Context, store *TaskStore) (*Service,error){
	loaded,err:=store.LoadTasks(ctx)
	if err!=nil{
		return nil,err
	}

	return &Service{
		store: store,
		tasks: loaded,
		nextID: calcNextID(loaded)
	},nil
}

func (s *Service) ListTasks(ctx context.Context, delay time.Duration) ([]Task, error) {
	// Зачем начинать долгую операцию,
	// если контекст уже отменен
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if delay > 0 {
		if err := s.store.SimulateSlowIO(ctx, delay); err != nil {
			return nil, err
		}
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock() // Блокируе только на запись
	defer s.mu.RUnlock()

	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out, nil
}

// GetTask возвращает задачу по id.
func (s *Service) GetTask(ctx context.Context, id int) (Task, bool, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.tasks {
		if t.ID == id {
			return t, true, nil
		}
	}
	return Task{}, false, nil
}

// CreateTask создаёт задачу и сохраняет в файл.
func (s *Service) CreateTask(ctx context.Context, incoming Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	created := Task{
		ID:       s.nextID,
		Title:    incoming.Title,
		Done:     incoming.Done,
		Priority: incoming.Priority, // [Валидация]
	}

	// Готовим новый список, но НЕ коммитим в память, пока не сохранили на диск.
	candidate := make([]Task, 0, len(s.tasks)+1)
	candidate = append(candidate, s.tasks...)
	candidate = append(candidate, created)

	if err := s.store.SaveTasks(ctx, candidate); err != nil {
		return Task{}, err
	}

	s.tasks = candidate
	s.nextID++
	return created, nil
}

// UpdateTask обновляет задачу по id и сохраняет в файл.
func (s *Service) UpdateTask(ctx context.Context, id int, incoming Task) (Task, bool, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return Task{}, false, nil
	}

	updated := s.tasks[idx]
	updated.Title = incoming.Title
	updated.Done = incoming.Done

	candidate := make([]Task, len(s.tasks))
	copy(candidate, s.tasks)
	candidate[idx] = updated

	if err := s.store.SaveTasks(ctx, candidate); err != nil {
		return Task{}, false, err
	}

	s.tasks = candidate
	return updated, true, nil
}

// DeleteTask удаляет задачу по id и сохраняет в файл.
func (s *Service) DeleteTask(ctx context.Context, id int) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i := range s.tasks {
		if s.tasks[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return false, nil
	}

	candidate := make([]Task, 0, len(s.tasks)-1)
	candidate = append(candidate, s.tasks[:idx]...)
	candidate = append(candidate, s.tasks[idx+1:]...)

	if err := s.store.SaveTasks(ctx, candidate); err != nil {
		return false, err
	}

	s.tasks = candidate
	return true, nil
}

func calcNextID(ts []taskstore.Task) int {
    maxID := 0
    for _, t := range ts {
        if t.ID > maxID {
            maxID = t.ID
        }
    }
    return maxID + 1
}