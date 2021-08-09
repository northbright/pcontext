package pcontext

import (
	"context"
	//"log"
	"sync"
	"sync/atomic"
)

type Context interface {
	context.Context
	Progress() <-chan ProgressData
	SetProgress(total, current int64)
}

type ProgressData struct {
	Total   int64
	Current int64
}

type pContext struct {
	context.Context
	mu       sync.Mutex
	progress atomic.Value
	done     bool
}

var closedchan = make(chan ProgressData)

func init() {
	close(closedchan)
}

func (pctx *pContext) createProgressCh() chan ProgressData {
	p := pctx.progress.Load()
	if p != nil {
		return p.(chan ProgressData)
	}

	pctx.mu.Lock()
	defer pctx.mu.Unlock()
	p = pctx.progress.Load()
	if p == nil {
		p = make(chan ProgressData)
		pctx.progress.Store(p)
	}

	return p.(chan ProgressData)
}

func (pctx *pContext) Progress() <-chan ProgressData {
	return pctx.createProgressCh()
}

func (pctx *pContext) SetProgress(total, current int64) {
	if total <= 0 || current < 0 {
		return
	}

	pctx.mu.Lock()
	if pctx.done {
		pctx.mu.Unlock()
		return // already canceled
	}

	p := pctx.createProgressCh()
	p <- ProgressData{total, current}
	pctx.mu.Unlock()
}

func WithProgress(ctx context.Context) Context {
	pctx := &pContext{Context: ctx}

	go func() {
		for {
			select {
			case <-pctx.Done():
				pctx.mu.Lock()
				pctx.done = true
				p, _ := pctx.progress.Load().(chan ProgressData)
				if p == nil {
					pctx.progress.Store(closedchan)
				} else {
					close(p)
				}
				pctx.mu.Unlock()
				return
			}
		}
	}()

	return pctx
}
