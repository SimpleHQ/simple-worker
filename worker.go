package worker

// Worker is what processes our tasks
type Worker interface {
	Start()
	Stop()
}

// BaseWorker listens for work requests and processes the correct task.
type BaseWorker struct {
	ID          int
	Work        Queue
	WorkerQueue chan Queue
	QuitChan    chan bool
	Processor   func(Payload)
}

// NewWorker builds new worker ready to accept tasks.
func NewWorker(id int, workerQueue chan Queue, processor func(Payload)) BaseWorker {
	return BaseWorker{
		ID:          id,
		Work:        make(Queue),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		Processor:   processor,
	}
}

// Start initiates a Worker. Listens for the work requests.
func (w BaseWorker) Start() {
	go func() {
		for {
			// Paylaod
			w.WorkerQueue <- w.Work

			select {
			case job := <-w.Work:
				// Process a task
				w.Processor(job)
			case <-w.QuitChan:
				// Qutting worker
				return
			}
		}
	}()
}

// Stop will stop the worker
func (w BaseWorker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
