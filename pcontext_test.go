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

	pctx := pcontext.WithProgress(ctx)

	go func(pctx pcontext.Context) {
		for i := 0; i <= 100; i += 10 {
			time.Sleep(time.Millisecond * 10)
			pctx.SetProgress(100, int64(i))
		}
		log.Printf("work goroutine exited")
	}(pctx)

	for pd := range pctx.Progress() {
		log.Printf("progress: total: %v, current: %v", pd.Total, pd.Current)
	}

	// The channel will be closed and for range loop will exit
	// when the context is done.
	log.Printf("task is done")

	// Output:
}
