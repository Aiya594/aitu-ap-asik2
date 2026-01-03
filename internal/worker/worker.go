package worker

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Aiya594/aitu-ap-asik2/internal/model"
	"github.com/Aiya594/aitu-ap-asik2/internal/store"
)

type WorkerPool struct {
	queue chan string
	store *store.TaskStore
	wg    *sync.WaitGroup
}

func NewWorkerPool(store *store.TaskStore, qSize int) *WorkerPool {
	return &WorkerPool{
		queue: make(chan string, qSize),
		store: store,
		wg:    &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Start(ctx context.Context, workers int) {
	for i := 0; i < workers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	defer wp.wg.Done()
	log.Printf("worker %d starting", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d shutting down", id)
			return

		case taskID, exists := <-wp.queue:
			if !exists {
				log.Printf("worker %d shutting down", id)
				return
			}
			wp.store.UpdateStatus(taskID, model.StatusInProgress)
			time.Sleep(2 * time.Second)
			wp.store.UpdateStatus(taskID, model.StatusDone)

		}
	}

}

func (wp *WorkerPool) Enqueue(taskID string) error {
	//wp.queue <- taskID
	select {
	case wp.queue <- taskID:
		return nil
	default:
		return errors.New("worker pool overloaded")
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
