package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var (
	onceCloser   sync.Once
	closer *Closer
)

type cancelFunc func(ctx context.Context) error

// Closer structure for managing shutdown functions
type Closer struct {
	funcs []cancelFunc
	mu    sync.RWMutex
}

// GetCloser returns pointer to Closer instance
func GetCloser() *Closer {
	onceCloser.Do(func() {
		closer = &Closer{}
	})

	return closer
}

// Add registering shutdown handlers
func (c *Closer) Add(fn cancelFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, fn)
}

// Close runs registered handlers
func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var wg sync.WaitGroup
	msgs := make([]string, 0, len(c.funcs))
	complete := make(chan bool)

	for _, fn := range c.funcs {
		wg.Add(1)

		go func(ctx context.Context, wg *sync.WaitGroup) {
			defer wg.Done()

			if err := fn(ctx); err != nil {
				msgs = append(msgs, err.Error())
			}
		}(ctx, &wg)
	}

	go func() {
		wg.Wait()
		complete <- true
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}
