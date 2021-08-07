package pcontext

import (
	"context"
	"sync"
	"sync/atomic"
)

type PContext interface {
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

func (pc *pContext) createProgressCh() (chan ProgressData, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	select {
	// Work done.
	case <-pc.Done():
		p, _ := pc.progress.Load().(chan ProgressData)
		if p == nil {
			p = closedchan
			pc.progress.Store(p)
		} else {
			close(p)
		}

		return p, true

	default:
	}

	// Work goroutine is running.
	p, _ := pc.progress.Load().(chan ProgressData)
	if p == nil {
		p = make(chan ProgressData)
		pc.progress.Store(p)
	}

	return p, false
}

func (pc *pContext) Progress() <-chan ProgressData {
	p, _ := pc.createProgressCh()
	return p
}

func (pc *pContext) SetProgress(total, current int64) {
	if total <= 0 || current < 0 {
		return
	}

	p, closed := pc.createProgressCh()
	if !closed {
		p <- ProgressData{total, current}
	}
}

/*
func WithProgress(ctx context.Context) PContext {
	return &pContext{context.Context: ctx, mu: sync.Mutex{}, progress: atomic.Value{}}
}
*/
