package worker

type Payload interface{}

// Queue has jobs waiting to be processed
type Queue chan Payload

// WorkerQueue is the queue for the worker
var WorkerQueue chan Queue

// Dispatcher prepares our workers
type Dispatcher interface {
	Start(maxWorkers int)
	Stop()
}

// BaseDispatcher is our base dispatcher type
type BaseDispatcher struct {
	Payload   Payload
	Workers   []Worker
	Processor func(Payload)
	Queue     Queue
}

// NewDispatcher creates a BaseDispatcher ready for starting
func NewDispatcher(maxQueue int, processor func(Payload)) BaseDispatcher {
	return BaseDispatcher{
		Processor: processor,
		Workers:   []Worker{},
		Queue:     make(Queue, maxQueue),
	}
}

// Start takes items from the queue and adds them to the worker queue
func (dispatcher BaseDispatcher) Start(maxWorkers int) {
	WorkerQueue = make(chan Queue, maxWorkers)

	// Create the amount of Workers requested
	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(i+1, WorkerQueue, dispatcher.Processor)

		// Add worker to our dispatchers list
		dispatcher.Workers = append(dispatcher.Workers, worker)

		// Start the worker
		worker.Start()
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
}

// Stop stops all known workers
func (dispatcher BaseDispatcher) Stop() {
	for _, worker := range dispatcher.Workers {
		worker.Stop()
	}
}
