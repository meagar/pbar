package pbar_test

import (
	"testing"

	"github.com/meagar/pbar"
)

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
