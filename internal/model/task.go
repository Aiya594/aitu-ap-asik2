package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID        string     `json:"id"`
	Payload   string     `json:"payload"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

type TaskResponse struct {
	ID     string     `json:"id"`
	Status TaskStatus `json:"status"`
}

func NewTask(payload string) *Task {
	id := uuid.New().String()
	return &Task{
		ID:        id,
		Payload:   payload,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}
