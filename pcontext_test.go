package pcontext_test

import (
	"context"
	"log"
	"time"

	"github.com/northbright/pcontext"
)

func Example() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	// Creates a progress context from parent.
	pctx := pcontext.WithProgress(ctx)

	// Run a work goroutine.
	go func(pctx pcontext.Context) {
		for i := 0; i <= 100; i += 10 {
			time.Sleep(time.Millisecond * 10)
			// Set progress in work goroutine.
			pctx.SetProgress(100, int64(i))
		}
		log.Printf("work goroutine exited")
	}(pctx)

	// Read progress data until the channel is closed.
	for pd := range pctx.Progress() {
		percent := pcontext.ComputePercent(pd.Total, pd.Current)
		log.Printf("progress: total: %v, current: %v, percent: %v", pd.Total, pd.Current, percent)
	}

	// The channel will be closed and for range loop will exit
	// when the context is done.
	log.Printf("task is done")

	// Output:
}
