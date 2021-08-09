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

	for {
		select {
		//case <-pctx.Done():
		//	log.Printf("timeout")
		case pd, ok := <-pctx.Progress():
			if ok {
				log.Printf("progress: total: %v, current: %v", pd.Total, pd.Current)
			} else {
				log.Printf("task is done")
				return
			}
		}
	}

	// Output:
}
