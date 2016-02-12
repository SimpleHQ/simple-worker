simple-worker
======

simple-worker is a minimal Go library that allows you to create Queue->Worker systems easily.

Installation
============

```
go get github.com/simplehq/simple-worker
```

Usage
=====

To use simple-worker, you first need a `Processor`. A processor is simply a function that takes a `interface{}` that it will then process.

```Go
import (
    "fmt"
    "github.com/simplehq/simple-worker"
)

type WorkRequest struct {
    Message string
}

func ProcessTask(job WorkRequest) {
    fmt.Println(job.Message)
}
```

To start a `dispatcher`, you simply call `StartDispatcher`.

```Go
maxWorkers := 10
maxQueue := 100

queue := worker.StartDispatcher(params.MaxWorkers, workQueue, ProcessTask)
```

`StartDispatcher` returns a `chan interface{}` for you to start adding jobs for processing.

```Go
queue := worker.StartDispatcher(params.MaxWorkers, workQueue, ProcessTask)

job := WorkRequest{
    Message: "Oh Boy"
}

queue <- job
```

This would ouput the message `"Oh Boy"` once the job has been processed.

Future Plans
============

- Dispatcher interface
