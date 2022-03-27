// Package pcontext provides a new context which derived from context.Context.
// It supports new feature(s): set / get progress.
package pcontext

import (
	"context"
	"sync"
	"sync/atomic"
)

// Context derives from context.Context, carries a progress.
// Work goroutines set the progress while other goroutines
// receive the progress data from the progress channel.
type Context interface {
	// Derive context.Context.
	context.Context
	// Progress returns the channel to receive progress data.
	Progress() <-chan ProgressData
	// SetProgress sets the progress in work goroutines.
	SetProgress(total, current int64)
}

// ProgressData contains the progress data.
type ProgressData struct {
	Total   int64
	Current int64
}

// pContext implements the Context interface.
type pContext struct {
	// Derived context.Context interface
	context.Context
	mu       sync.Mutex
	progress atomic.Value
	done     bool
}

// closedchan is a reusable closed channel.
var closedchan = make(chan ProgressData)

func init() {
	close(closedchan)
}

// ComputePercent computes the progress percent.
func ComputePercent(total, current int64) float32 {
	if total > 0 {
		return float32(float64(current) / (float64(total) / float64(100)))
	}
	return 0
}

// createProgressCh creates the progress channel dynamically.
// It'll create the channel the first createProgressCh is called.
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

// Progress returns the channel to receive progress data.
func (pctx *pContext) Progress() <-chan ProgressData {
	return pctx.createProgressCh()
}

// SetProgress sets the progress in work goroutines.
func (pctx *pContext) SetProgress(total, current int64) {
	if total <= 0 || current < 0 {
		return
	}

	pctx.mu.Lock()
	if pctx.done {
		pctx.mu.Unlock()
		return // already canceled
	}
	pctx.mu.Unlock()

	p := pctx.createProgressCh()
	p <- ProgressData{total, current}
}

// WithProgress returns a copy of parent with progress supported context.
func WithProgress(ctx context.Context) Context {
	pctx := &pContext{Context: ctx}

	go func() {
		<-pctx.Done()
		pctx.mu.Lock()
		pctx.done = true
		p, _ := pctx.progress.Load().(chan ProgressData)
		if p == nil {
			pctx.progress.Store(closedchan)
		} else {
			close(p)
		}
		pctx.mu.Unlock()
	}()

	return pctx
}
