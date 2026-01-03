package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Aiya594/aitu-ap-asik2/internal/api"
	"github.com/Aiya594/aitu-ap-asik2/internal/model"
	"github.com/Aiya594/aitu-ap-asik2/internal/router"
	"github.com/Aiya594/aitu-ap-asik2/internal/store"
	"github.com/Aiya594/aitu-ap-asik2/internal/worker"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	repo := store.NewRepository[string, model.Task]()
	taskStore := store.NewTaskStore(repo)
	log.Println("TaskStore created")

	workerPool := worker.NewWorkerPool(taskStore, 15)
	workerPool.Start(ctx, 2)
	log.Println("Workerpool started")

	monitor := worker.NewMonitor(taskStore)
	monitor.Start(ctx)
	log.Println("Monitor started")

	handler := api.NewHandler(taskStore, workerPool)
	mux := router.Router(handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Starting server on http://localhost:8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Println("Server shutdown error:", err)
	}

	workerPool.Wait()
	monitor.Wait()

	log.Println("Server gracefully stopped")

}
