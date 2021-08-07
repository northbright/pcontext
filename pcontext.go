package pcontext

import (
	"context"
	"log"
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
}

var closedchan = make(chan ProgressData)

func init() {
	close(closedchan)
}

func (pctx *pContext) createProgressCh() (chan ProgressData, bool) {
	pctx.mu.Lock()
	defer pctx.mu.Unlock()

	select {
	// Work done.
	case <-pctx.Done():
		log.Printf("pctx.Done()")
		p, _ := pctx.progress.Load().(chan ProgressData)
		if p == nil {
			p = closedchan
			pctx.progress.Store(p)
		}

		return p, true

	default:
	}

	// Work goroutine is running.
	p, _ := pctx.progress.Load().(chan ProgressData)
	if p == nil {
		p = make(chan ProgressData)
		pctx.progress.Store(p)
	}

	return p, false
}

func (pctx *pContext) Progress() <-chan ProgressData {
	p, _ := pctx.createProgressCh()
	return p
}

func (pctx *pContext) SetProgress(total, current int64) {
	if total <= 0 || current < 0 {
		return
	}

	p, closed := pctx.createProgressCh()
	if !closed {
		p <- ProgressData{total, current}
	}
}

func WithProgress(ctx context.Context) Context {
	pctx := &pContext{Context: ctx}

	go func() {
		for {
			select {
			case <-pctx.Done():
				pctx.mu.Lock()
				p, _ := pctx.progress.Load().(chan ProgressData)
				if p == nil {
					p = closedchan
					pctx.progress.Store(p)
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
