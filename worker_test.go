package worker_test

import (
	"sync"
	"testing"

	"github.com/simplehq/simple-worker"
)

func TestProcessMultipleJobs(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	dispatcher := worker.NewDispatcher(100, func(payload worker.Payload) {
		// Processed
		wg.Done()
	})
	dispatcher.Start(1)

	go func() {
		for i := 0; i < 10; i++ {
			// Add to queue
			dispatcher.Queue <- i
		}
	}()

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
