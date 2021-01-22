// Package pbar provides a simple progress bar that includes display of records-per-second (RPS) and an ETA
package pbar

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// The ETA and OPS measurements use a moving average over the last 25 samples
const sampleSize = 25

// Bar maintains the progressbar's state
type Bar struct {
	total     uint64
	startTime time.Time
	lastTime  time.Time
	progress  uint64
	width     int

	// Tracks the last N states so we can average the records-per-second and produce a slightly
	// less jumpy number
	samples     [sampleSize]float32
	sampleIndex int
}

// Options contains the configurable options that may be passed in to New
type Options struct {
	// The total number of records that will be processed
	Total uint64

	// The width of the progress bar; if omitted will default to terminal width
	Width int

	// The time at which operations started; if omitted, defaults to time.Now()
	Start time.Time
}

// New constructs and returns a new Bar instance
func New(opts Options) *Bar {
	if opts.Total <= 0 {
		panic("Cannot create a progressbar with a 0 total")
	}

	if opts.Width == 0 {
		opts.Width = terminalWidth()
	}

	if opts.Start == (time.Time{}) {
		opts.Start = time.Now()
	}

	bar := Bar{
		total:     opts.Total,
		lastTime:  opts.Start,
		startTime: opts.Start,
		width:     opts.Width,
	}

	return &bar
}

// Tick advances the internal state to the given progress, allowing the progress bar
// to measure how much progress has been made since the previous call to `Tick`.
// Useful when you want to set the bar's progress to an absolute value.
//
// Example:
//	for i := 0; i < 100; i++ {
//		// Set the total number of operations performed since the loop started (i)
//		bar.Tick(i)
//	}
func (p *Bar) Tick(progress uint64) {
	elapsed := time.Since(p.lastTime)
	delta := float32(progress - p.progress)

	// operations per second
	ops := delta / float32(elapsed.Seconds())

	p.samples[p.sampleIndex] = ops
	p.sampleIndex++
	if p.sampleIndex >= sampleSize {
		p.sampleIndex = 0
	}

	p.lastTime = time.Now()
	p.progress = progress
}

// TickDelta advances the internal state of the progress bar by the given amount.
// Useful when you want to increment the bar by a specific number.
//
// Example:
//	for i := 0; i < 100; i += 5 {
//		// Increment by the number of operations performed this iteration (5)
//		bar.TickDelta(1)
//	}
func (p *Bar) TickDelta(delta uint64) {
	p.Tick(p.progress + delta)
}

// Progress returns a string containing the rendered progress bar,
// with leading ANSI control characters to blank out the previous progress bar so that
// the bar appears to advance in-place.
func (p *Bar) Progress() string {
	buf := strings.Builder{}
	fmt.Fprintf(&buf, "\033[2K\033[%dD", p.width)

	ops := p.avg()
	// How many records were processed this tick
	percent := float32(p.progress) / float32(p.total)
	remainingSeconds := float32(p.total-p.progress) / ops
	eta := time.Until(time.Now().Add(time.Second * time.Duration(remainingSeconds)))
	fmt.Fprintf(&buf, "%d of %d (%d%%) - OPS: %.2f - ETA: %s",
		p.progress, p.total, int(percent*100), ops, eta.Round(time.Second))

	p.renderBar(&buf, percent)
	return buf.String()
}

// Summary displays a summary of the progress bar, after completion.
// For example:
//   "Completed 5000 of 5000 records in 5m12s"
func (p *Bar) Summary() string {
	buf := strings.Builder{}
	fmt.Fprintf(&buf, "\033[2K\033[%dD", p.width)
	fmt.Fprintf(&buf, "Completed %d records in %s\n", p.progress, time.Since(p.startTime).Round(time.Second))
	return buf.String()
}

func (p *Bar) renderBar(buf *strings.Builder, percent float32) {
	if percent < 0 {
		panic("sub-zero percent")
	}
	barWidth := p.width - buf.Len() - 3 // three chars for " []"
	barCompleteWidth := int(percent * float32(barWidth))
	barRemainingWidth := barWidth - barCompleteWidth

	buf.WriteString(" [")
	buf.Write(bytes.Repeat([]byte("="), barCompleteWidth))
	buf.Write(bytes.Repeat([]byte("-"), barRemainingWidth))
	buf.WriteString("]")
}

// avg returns the average over the last N samples
func (p *Bar) avg() float32 {
	var sum float32

	for _, n := range p.samples {
		sum += n
	}

	return sum / sampleSize
}

// Attempt to get the width of the terminal. Probably very non-portable, so if anything errors out we return 79
func terminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error finding terminal width:", err)
		return 79
	}

	parts := strings.Split(string(out), " ")

	width, err := strconv.Atoi(strings.TrimSpace(parts[1]))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error finding terminal width:", err)
		return 79
	}

	return width
}
