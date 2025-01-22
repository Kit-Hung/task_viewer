package model

import "time"

// Task 任务模型
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskDependency 任务依赖关系模型
type TaskDependency struct {
	ID        int       `json:"id"`
	SourceID  int       `json:"source_id"`
	TargetID  int       `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskResponse 任务响应模型
type TaskResponse struct {
	Task         *Task   `json:"task"`
	Dependencies []*Task `json:"dependencies"`
}
