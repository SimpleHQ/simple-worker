package worker

// Payload is what the task will contain
type Payload interface{}

// WorkerQueue is the queue for the worker
var WorkerQueue chan chan Payload

// Dispatcher prepares our workers
type Dispatcher interface {
	StartWorkers(int) []BaseWorker
	StopWorkers()
	AddJob(Payload)
}

// BaseDispatcher is our base dispatcher type
type BaseDispatcher struct {
	// Processor is called when a job is processed
	Processor func(Payload)
	// Queue is where we send our jobs
	Queue chan Payload
	// Workers is our slice of worker interfaces
	Workers []BaseWorker
}

// NewDispatcher creates a BaseDispatcher ready for starting
func NewDispatcher(maxQueue int, processor func(Payload)) *BaseDispatcher {
	return &BaseDispatcher{
		Processor: processor,
		Queue:     make(chan Payload, maxQueue),
	}
}

// StartWorkers takes items from the queue and adds them to the worker queue
func (dispatcher *BaseDispatcher) StartWorkers(maxWorkers int) []BaseWorker {
	// WorkerQueue contains the workers
	WorkerQueue = make(chan chan Payload, maxWorkers)

	workers := []BaseWorker{}

	// Create the amount of Workers requested
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(i+1, WorkerQueue, dispatcher.Processor)

		// Start the worker
		worker.Start()

		// Add worker to our worker list
		workers = append(workers, worker)
	}

	dispatcher.Workers = workers

	// Goroutine to wait for jobs
	go func() {
		for {
			select {
			case job := <-dispatcher.Queue:
				// Job has been popped off our dispatcher queue

				go func() {
					// Get the next available worker
					worker := <-WorkerQueue

					// Send the job to the worker
					worker <- job
				}()
			}
		}
	}()

	return workers
}

// AddJob will add a job to the queue to be processed
// It is merely a helper method. You should really just push to the channel.
func (dispatcher *BaseDispatcher) AddJob(payload Payload) {
	dispatcher.Queue <- payload
}

// StopWorkers will tell all workers to stop
func (dispatcher *BaseDispatcher) StopWorkers() {
	for _, worker := range dispatcher.Workers {
		worker.Stop()
	}
}
