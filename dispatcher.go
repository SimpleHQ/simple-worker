package worker

// WorkerQueue is the list of items to be processed.
var WorkerQueue chan chan WorkRequest

// StartDispatcher takes items from the Queue and sends them to a Worker.
func StartDispatcher(nworkers int, nqueue int, taskProcessor func(WorkRequest)) chan WorkRequest {
	queue := make(chan WorkRequest, nqueue)
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	for i := 0; i < nworkers; i++ {
		// Starting n workers
		worker := NewWorker(i+1, WorkerQueue, taskProcessor)
		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-queue:
				// Job added to global queue
				go func() {
					worker := <-WorkerQueue

					// Add job to to worker
					worker <- job
				}()
			}
		}
	}()

	return queue
}
