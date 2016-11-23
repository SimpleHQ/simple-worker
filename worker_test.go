package worker_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/simplehq/simple-worker"
)

func TestProcessMultipleJobs(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	dispatcher := worker.NewDispatcher(100, func(payload worker.Payload) {
		// Processed
		wg.Done()
	})
	dispatcher.StartWorkers(1)

	go func() {
		for i := 0; i < 10; i++ {
			// Add to queue
			dispatcher.Queue <- i
		}
	}()

	wg.Wait()
}

func TestDispatcherStopsWorkers(t *testing.T) {
	var wg sync.WaitGroup

	dispatcher := worker.NewDispatcher(100, func(payload worker.Payload) {})

	// Between 1 & 6 workers
	rand.Seed(time.Now().Unix())
	numWorkers := rand.Intn(5) + 1

	workers := []worker.BaseWorker{}
	for i := 1; i <= numWorkers; i++ {
		workers = append(workers, worker.NewWorker(i, nil, nil))
	}
	dispatcher.Workers = workers

	wg.Add(numWorkers)

	go func() {
		for _, worker := range dispatcher.Workers {
			<-worker.QuitChan
			wg.Done()
		}
	}()

	dispatcher.StopWorkers()
	wg.Wait()
}

func TestStopWorker(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	theWorker := worker.NewWorker(1, nil, nil)

	theWorker.Start()

	go func() {
		// Once quitted, we are done
		<-theWorker.QuitChan
		wg.Done()
	}()

	theWorker.Stop()

	// Wait until we receive done request
	wg.Wait()
}
