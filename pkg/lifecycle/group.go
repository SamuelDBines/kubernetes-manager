package lifecycle

import (
	"context"
	"sync"
)

type Group struct {
	mu     sync.Mutex
	actors []actor
}
type actor struct {
	start     func() error
	interrupt func(err error)
}

func (g *Group) Add(start func() error, interrupt func(err error)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.actors = append(g.actors, actor{
		start:     start,
		interrupt: interrupt,
	})
}

func (g *Group) Run() error {
	var wg sync.WaitGroup
	errc := make(chan error, len(g.actors))

	// start all
	for _, a := range g.actors {
		wg.Add(1)
		go func(a actor) {
			defer wg.Done()
			errc <- a.start()
		}(a)
	}
	err := <-errc

	g.mu.Lock()
	for _, a := range g.actors {
		go a.interrupt(err)
	}
	g.mu.Unlock()

	wg.Wait()
	return err
}

// Helper: derive a cancellable context & add cancelers to the group.
func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}
