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

To use simple-worker, you first need a `Processor`. A processor is simply a function that takes a `WorkRequest` that it will then process.

```Go
import (
    "fmt"
    "github.com/simplehq/simple-worker"
)

func ProcessTask(job worker.WorkRequest) {
    fmt.Println("Job processed")
}
```

To start a `dispatcher`, you simply call `StartDispatcher`.

```Go
maxWorkers := 10
maxQueue := 100

queue := worker.StartDispatcher(params.MaxWorkers, workQueue, ProcessTask)
```

`StartDispatcher` returns a `chan WorkRequest` for you to start adding jobs for processing.

```Go
queue := worker.StartDispatcher(params.MaxWorkers, workQueue, ProcessTask)

job := worker.WorkRequest{
    Type: 1,
    Message: "Oh Boy"
}

queue <- job
```

Future Plans
============

- WorkRequest interface for custom job structure
- Dispatcher interface
