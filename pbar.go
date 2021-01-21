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

// How many samples we've collected
const sampleSize = 10

// Bar maintains the progressbar's state
type Bar struct {
	total     int64
	startTime time.Time
	lastTime  time.Time
	progress  int64
	width     int

	// Tracks the last N states so we can average the records-per-second and produce a slightly
	// less jumpy number
	samples     [sampleSize]float64
	sampleIndex int
}

// Options contains the configurable options that may be passed in to New
type Options struct {
	// The total number of records that will be processed
	Total int64

	// The width of the progress bar; if omitted will default to terminal width
	Width int
}

// New constructs and returns a new Bar instance
func New(opts *Options) *Bar {
	width := opts.Width

	if width == 0 {
		width = terminalWidth()
	}

	bar := Bar{
		total:     opts.Total,
		startTime: time.Now(),
		width:     terminalWidth(),
	}

	return &bar
}

// Tick advances the internal state of the progress bar by the given amount, allowing the progress bar
// to measure how much progress has been made since the previous call to `Tick`.
func (p *Bar) Tick(progress int64) {
	elapsed := time.Since(p.lastTime)
	records := float64(progress - p.progress)

	rps := records / elapsed.Seconds()

	p.samples[p.sampleIndex] = rps
	p.sampleIndex++
	if p.sampleIndex >= sampleSize {
		p.sampleIndex = 0
	}

	p.lastTime = time.Now()
	p.progress = progress
}

// Progress returns a string containing the rendered progress bar,
// with leading ANSI control characters to blank out the previous progress bar so that
// the bar appears to advance in-place.
func (p *Bar) Progress() string {
	buf := strings.Builder{}
	buf.WriteString("\033[2K")                        // Clear to start of line
	buf.WriteString(fmt.Sprintf("\033[%dD", p.width)) // Move cursor to start of line

	rps := p.avg()
	// How many records were processed this tick
	percent := float64(p.progress) / float64(p.total)
	remainingSeconds := float64(p.total-(p.progress)) / rps
	eta := time.Until(time.Now().Add(time.Second * time.Duration(remainingSeconds)))
	fmt.Fprintf(&buf, "%d of %d (%d%%) - RPS: %.2f - ETA: %s",
		p.progress, p.total, int(percent*100), rps, eta.Round(time.Second))

	p.renderBar(&buf, percent)
	return buf.String()
}

// Summary displays a summary of the progress bar, after completion.
func (p *Bar) Summary() string {
	buf := strings.Builder{}
	buf.WriteString("\033[2K")                        // Clear to start of line
	buf.WriteString(fmt.Sprintf("\033[%dD", p.width)) // Move cursor to start of line
	fmt.Fprintf(&buf, "Completed %d records in %v\n", p.progress, time.Since(p.startTime))
	return buf.String()
}

func (p *Bar) renderBar(buf *strings.Builder, percent float64) {
	barWidth := p.width - buf.Len() - 3 // three chars for " []"
	buf.WriteString(" [")
	barCompleteWidth := int(percent * float64(barWidth))
	barRemainingWidth := barWidth - barCompleteWidth
	buf.Write(bytes.Repeat([]byte("="), barCompleteWidth))
	buf.Write(bytes.Repeat([]byte("-"), barRemainingWidth))
	buf.WriteString("]")
}

// avg returns the
func (p *Bar) avg() float64 {
	var sum float64

	for _, n := range p.samples {
		sum += n
	}

	return sum / sampleSize
}

func terminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 79
	}

	parts := strings.Split(string(out), " ")

	width, err := strconv.Atoi(strings.TrimSpace(parts[1]))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 79
	}

	return width
}
