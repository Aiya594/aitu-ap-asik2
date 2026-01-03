package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Aiya594/aitu-ap-asik2/internal/model"
	"github.com/Aiya594/aitu-ap-asik2/internal/store"
	"github.com/Aiya594/aitu-ap-asik2/internal/worker"
)

type Handler struct {
	store *store.TaskStore
	wp    *worker.WorkerPool
}

func NewHandler(store *store.TaskStore, wp *worker.WorkerPool) *Handler {
	return &Handler{
		store: store,
		wp:    wp,
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}
	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}

	newTask := model.NewTask(task.Payload)

	err = h.wp.Enqueue(newTask.ID)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
	}

	taskResponse, err := h.store.CreateTask(newTask)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusBadRequest, fmt.Sprintf("error: %v", err))
	}
	log.Printf("Task with id=%v created", task.ID)

	writeJSON(w, http.StatusCreated, taskResponse)
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.GetAllTasks()
	log.Println("All tasks received")

	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, exists := h.store.GetTask(id)
	if !exists {
		log.Println("Task not found")
		writeJSON(w, http.StatusNotFound, nil)
	}
	log.Println("Task received")
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats := h.store.Statistics()
	writeJSON(w, http.StatusOK, stats)
}
