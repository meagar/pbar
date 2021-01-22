package pbar_test

import (
	"fmt"
	"time"

	"github.com/meagar/pbar"
)

func ExampleBar_Progress() {
	bar := pbar.New(pbar.Options{
		Width: 80,
		Total: 360,
		Start: time.Now().Add(-5 * time.Second),
	})

	// Advance the progress bar to the halfway point
	bar.Tick(180)

	// Strip leading ANSI control characters
	out := bar.Progress()[9:]

	fmt.Println(out)
	// Output: 180 of 360 (50%) - OPS: 1.44 - ETA: 2m5s [==============--------------]
}

func ExampleBar_Summary() {
	bar := pbar.New(pbar.Options{
		Total: 100,
		Start: time.Now().Add(-10 * time.Second),
		Width: 80,
	})

	bar.Tick(100)

	// Strip leading ANSI control characters
	out := bar.Summary()[9:]

	fmt.Println(out)
	// Output: Completed 100 records in 10s
}
