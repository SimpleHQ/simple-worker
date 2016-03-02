package worker

// Payload is what the task will contain
type Payload interface{}

// Queue has jobs waiting to be processed
type Queue chan Payload

// WorkerQueue is the queue for the worker
var WorkerQueue chan Queue

// Dispatcher prepares our workers
type Dispatcher interface {
	Start(int) []BaseWorker
	AddJob(Payload)
}

// BaseDispatcher is our base dispatcher type
type BaseDispatcher struct {
	Payload   Payload
	Processor func(Payload)
	Queue     Queue
}

// NewDispatcher creates a BaseDispatcher ready for starting
func NewDispatcher(maxQueue int, processor func(Payload)) BaseDispatcher {
	return BaseDispatcher{
		Processor: processor,
		Queue:     make(Queue, maxQueue),
	}
}

// Start takes items from the queue and adds them to the worker queue
func (dispatcher BaseDispatcher) Start(maxWorkers int) []BaseWorker {
	WorkerQueue = make(chan Queue, maxWorkers)
	workers := []BaseWorker{}

	// Create the amount of Workers requested
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(i+1, WorkerQueue, dispatcher.Processor)

		// Start the worker
		worker.Start()

		// Add worker to our dispatchers list
		workers = append(workers, worker)
	}

	// Goroutine to wait for jobs
	go func() {
		for {
			select {
			case job := <-dispatcher.Queue:
				// Goroutine to add job to workers queue
				go func() {
					worker := <-WorkerQueue
					worker <- job
				}()
			}
		}
	}()

	return workers
}

// AddJob will add a job to the queue to be processed
func (dispatcher *BaseDispatcher) AddJob(payload Payload) {
	dispatcher.Queue <- payload
}
