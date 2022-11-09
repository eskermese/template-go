// golang.org/x/sync/errgroup -> Group

// Package workers provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
package workers

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	cancel func()

	wg sync.WaitGroup
	mu sync.Mutex

	errs *multierror.Error
}

// GroupWithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func GroupWithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel()
	}

	return g.errs.ErrorOrNil()
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; all errors will be
// accumulated and returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.mu.Lock()
			g.errs = multierror.Append(g.errs, err)
			g.mu.Unlock()

			if g.cancel != nil {
				g.cancel()
			}
		}
	}()
}
