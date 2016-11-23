package worker

// Worker is what processes our tasks
type Worker interface {
	Start()
	Stop()
}

// BaseWorker listens for work requests and processes the correct task.
type BaseWorker struct {
	// A basic identifier
	ID int
	// Work is a channel of the current job being processed.
	// Once it is empty, it means the worker is free to pick up another job.
	Work chan Payload
	// WorkerQueue is the buffered channel of jobs being processed
	// It contains each workers Work field
	WorkerQueue chan chan Payload
	// QuitChan is used to stop the worker
	QuitChan chan bool
	// Processor is called when a job is processed
	Processor func(Payload)
}

// NewWorker builds new worker ready to accept tasks.
func NewWorker(id int, workerQueue chan chan Payload, processor func(Payload)) BaseWorker {
	return BaseWorker{
		ID:          id,
		Work:        make(chan Payload),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		Processor:   processor,
	}
}

// Start initiates a Worker. Listens for the work requests.
func (w BaseWorker) Start() {
	go func() {
		for {
			// Place the worker in the worker queue
			w.WorkerQueue <- w.Work

			select {
			case job := <-w.Work:
				// Process a task
				w.Processor(job)
				continue
			case <-w.QuitChan:
			}
			break
		}
	}()
}

// Stop will stop the worker
func (w BaseWorker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
