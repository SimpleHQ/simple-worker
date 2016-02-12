package worker_test

import (
	"testing"

	"github.com/simplehq/simple-worker"
)

type WorkRequest string

func TestWorker(t *testing.T) {
	expected := "Hi"

	dispatcher := worker.NewDispatcher(100, func(payload worker.Payload) {
		request := payload.(WorkRequest)
		if request != expected {
			panic(request + " does not match " + expected)
		}
	})

	dispatcher.Queue <- expected
}
