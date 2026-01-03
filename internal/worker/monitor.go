package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Aiya594/aitu-ap-asik2/internal/store"
)

type Monitor struct {
	store *store.TaskStore
	wg    sync.WaitGroup
}

func NewMonitor(store *store.TaskStore) *Monitor {
	return &Monitor{
		store: store,
		wg:    sync.WaitGroup{},
	}
}

func (m *Monitor) Start(ctx context.Context) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := m.store.Statistics()
				log.Printf("Statistics: submitted=%d, completed=%d, in_progress=%d\n",
					stats.Submitted, stats.Completed, stats.InProgress)
			case <-ctx.Done():
				log.Println("monitor stopped")
				return
			}
		}
	}()
}

func (m *Monitor) Wait() {
	m.wg.Wait()
}
