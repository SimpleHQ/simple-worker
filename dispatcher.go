package worker

// WorkerQueue is the list of items to be processed.
var WorkerQueue chan chan interface{}

// StartDispatcher takes items from the Queue and sends them to a Worker.
func StartDispatcher(nworkers int, nqueue int, taskProcessor func(interface{})) chan interface{} {
	queue := make(chan interface{}, nqueue)
	WorkerQueue = make(chan chan interface{}, nworkers)

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
