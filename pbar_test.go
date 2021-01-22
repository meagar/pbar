package pbar_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/meagar/pbar"
)

func TestPbar(t *testing.T) {
	bar := pbar.New(pbar.Options{
		Total: 10000,
		Width: 100,
	})

	for i := uint64(0); i < 10000; i++ {
		if i%100 == 0 {
			bar.Tick(i)
		}

		time.Sleep(time.Duration(0.01 * float64(time.Second)))
		fmt.Print(bar.Progress())
	}

}

func TestProgress(t *testing.T) {
	bar := pbar.New(pbar.Options{
		Total: 100,
		Width: 100,
	})

	for i := 0; i < 100; i++ {
		bar.TickDelta(1)
		bar.Progress()
	}
}

func BenchmarkProgress(b *testing.B) {
	bar := pbar.New(pbar.Options{
		Total: uint64(b.N),
		Width: 100,
	})

	for i := 0; i < b.N; i++ {
		bar.Tick(uint64(i))
		bar.Progress()
	}
}
