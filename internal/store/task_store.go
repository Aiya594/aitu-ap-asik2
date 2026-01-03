package store

import (
	"errors"
	"strings"

	"github.com/Aiya594/aitu-ap-asik2/internal/model"
)

type TaskStore struct {
	Repo      *Repository[string, model.Task]
	Submitted int
	Completed int
}

func NewTaskStore(repo *Repository[string, model.Task]) *TaskStore {
	return &TaskStore{
		Repo: repo,
	}
}

func (s *TaskStore) CreateTask(task *model.Task) (model.TaskResponse, error) {

	if strings.TrimSpace(task.Payload) == "" {
		return model.TaskResponse{}, errors.New("payload is empty")
	}
	s.Submitted++
	s.Repo.Post(task.ID, *task)
	return model.TaskResponse{
		ID:     task.ID,
		Status: task.Status,
	}, nil
}

func (s *TaskStore) GetTask(id string) (model.Task, bool) {

	return s.Repo.Get(id)
}

func (s *TaskStore) GetAllTasks() []model.TaskResponse {
	tasks := s.Repo.GetAll()

	var result []model.TaskResponse
	for _, task := range tasks {
		tsk := model.TaskResponse{
			ID:     task.ID,
			Status: task.Status,
		}
		result = append(result, tsk)
	}

	return result
}

func (s *TaskStore) UpdateStatus(id string, status model.TaskStatus) bool {

	return s.Repo.Update(id, func(t *model.Task) error {
		t.Status = status

		if status == model.StatusDone {
			s.Submitted--
			s.Completed++
		}

		return nil
	})
}

func (s *TaskStore) Statistics() model.Statistics {
	return model.Statistics{
		Submitted:  s.Submitted,
		Completed:  s.Completed,
		InProgress: (s.Submitted + s.Completed) - s.Completed,
	}
}
