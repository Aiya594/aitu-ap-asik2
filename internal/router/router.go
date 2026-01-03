package router

import (
	"net/http"

	"github.com/Aiya594/aitu-ap-asik2/internal/api"
)

func Router(h *api.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", h.CreateTask)

	mux.HandleFunc("GET /tasks", h.GetAllTasks)

	mux.HandleFunc("GET /tasks/{id}", h.GetTask)

	mux.HandleFunc("GET /stats", h.GetStatistics)

	return mux
}
