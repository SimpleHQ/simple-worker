package worker

type Job interface {
	Start()
	Stop()
}

// Worker listens for work requests and processes the correct task.
type Worker struct {
	ID          int
	Work        chan interface{}
	WorkerQueue chan chan interface{}
	QuitChan    chan bool
	Processor   func(interface{})
}

// NewWorker builds new worker ready to accept tasks.
func NewWorker(id int, workerQueue chan chan interface{}, processor func(interface{})) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan interface{}),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		Processor:   processor,
	}

	return worker
}

// Start initiates a Worker. Listens for the work requests.
func (w Worker) Start() {
	go func() {
		for {
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

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
