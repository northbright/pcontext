package pcontext_test

import (
	"context"
	"fmt"
	"time"

	"github.com/northbright/pcontext"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	// Creates a progress context from parent and a channel to receive progress data.
	pctx, ch := pcontext.WithProgress(ctx)

	// Run a work goroutine.
	go func(pctx pcontext.Context) {
		for i := 0; i <= 100; i += 10 {
			time.Sleep(time.Millisecond * 10)
			// Set progress in work goroutine.
			pctx.SetProgress(100, int64(i))
		}
		fmt.Printf("work goroutine exited\n")
	}(pctx)

	// Read progress data until the channel is closed.
	for pd := range ch {
		percent := pcontext.ComputePercent(pd.Total, pd.Current)
		fmt.Printf("progress: total: %v, current: %v, percent: %v%%\n", pd.Total, pd.Current, percent)
	}

	// The channel will be closed and for range loop will exit
	// when the context is done.
	fmt.Printf("task is done\n")

	// Output:
	// progress: total: 100, current: 0, percent: 0%
	// progress: total: 100, current: 10, percent: 10%
	// progress: total: 100, current: 20, percent: 20%
	// progress: total: 100, current: 30, percent: 30%
	// progress: total: 100, current: 40, percent: 40%
	// progress: total: 100, current: 50, percent: 50%
	// progress: total: 100, current: 60, percent: 60%
	// progress: total: 100, current: 70, percent: 70%
	// progress: total: 100, current: 80, percent: 80%
	// progress: total: 100, current: 90, percent: 90%
	// work goroutine exited
	// progress: total: 100, current: 100, percent: 100%
	// task is done
}
