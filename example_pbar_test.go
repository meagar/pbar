package pbar_test

import (
	"fmt"
	"time"

	"github.com/meagar/pbar"
)

// ExamplePBar demonstrates how to use pbar to create a progress bar in a simple for loop
func Example() {
	bar := pbar.New(pbar.Options{
		Total: 100,
		Width: 80,
	})

	for i := 0; i < 10; i++ {
		bar.TickDelta(1)
		fmt.Print(bar.Progress())
		time.Sleep(time.Second)
	}

	fmt.Print(bar.Summary())
}
