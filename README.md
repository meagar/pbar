# pbar - A simple Progress bar for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/meagar/pbar.svg)](https://pkg.go.dev/github.com/meagar/pbar)

pbar is (yet another) simple progress bar for Go.
It's centered around the idea that you are probably iterating over a set of operations, and displays:

```
263094 of 8635491 (3%) - OPS: 2009.81 - ETA: 1h9m25s [=---------------------------------]
```

The fields are (left to right):

* Your progress against the total operations (`x of y`)
* The completion rate (`x%`)
* Your average operations per second (`OPS:`)
* The estimated time remaining (`ETA:`)
* A progress bar that scales to the width of the terminal

## Usage

1. Create a new bar with `pbar.New`, passing in the total number of records and (optionally) the maximum width of the bar, in case you don't want it to span the whole window. This is useful if the non-portal method used to get the terminal width fails on your platform.

2. Call `bar.Tick()` and pass in the current progress through your set of operations. You can call this inside a loop after each operation, or periodically each time you hit a specific number of records, or in a timer with `time.After` (See examples below).
   
    Internally `Tick` will measure the amount of time since the last call to `Tick` to determine the average operations-per-second

3. Call `fmt.Print(bar.Progress())` to display the progress bar

4. After completion, call `bar.Summary()` to display the total elapsed time.
