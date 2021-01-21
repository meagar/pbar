package pbar

import (
	"fmt"
)

func ExampleBar_Progress() {
	bar := New(&Options{
		Width: 80,
		Total: 360,
	})

	// Advance the progress bar to the halfway point
	bar.Tick(180)

	// Strip leading ANSI control characters
	out := bar.Progress()[9:]

	fmt.Println(out)
	// Output: 180 of 360 (50%) - RPS: 0.00 - ETA: -1s [==============--------------]
}
