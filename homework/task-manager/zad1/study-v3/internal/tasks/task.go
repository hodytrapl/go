package tasks

// Task - модель задачи
//
// Хранятся в памяти для скорости и сериализуется в JSON
type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Done     bool   `json:"done"`
	Priority string `json:"priority"`
}

type CreateTaskRequest struct {
	Title    string `json:"title" validate:"required,max=100"`
	Done     bool   `json:"done"`
	Priority string `json:"priority" validate:"required,oneof=low medium high"`
}
