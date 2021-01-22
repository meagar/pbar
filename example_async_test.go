package pbar_test

import (
	"fmt"
	"time"

	"github.com/meagar/pbar"
)

// ExampleAsync demonstrates how to use pbar an in asynchronous example
func ExampleAsync() {
	bar := pbar.New(pbar.Options{
		Total: 10,
		Width: 80,
	})

	done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			bar.TickDelta(1)
			time.Sleep(time.Second)
		}
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Print(bar.Summary())
			return
		case <-time.After(time.Duration(0.5 * float64(time.Second))):
			fmt.Print(bar.Progress())
		}
	}
}
