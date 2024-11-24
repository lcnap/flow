package main

import (
	"log/slog"
	"sync"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	var wg sync.WaitGroup
	type hold struct {
		logger *slog.Logger
	}

	h := hold{}

	h.logger = NewLogger("logs/testreopen-old.log")

	wg.Add(1)
	go func() {

		for i := 0; i < 50_000; i++ {
			h.logger.Info("msg from old logger.", "count", i)
		}
		wg.Done()

	}()
	time.Sleep(200 * time.Millisecond)
	h.logger = NewLogger("logs/testreopen-new.log")

	wg.Add(1)
	go func() {

		for i := 0; i < 50_000; i++ {
			h.logger.Info("msg from new logger.", "count", i)
		}
		wg.Done()

	}()

	wg.Wait()

}
