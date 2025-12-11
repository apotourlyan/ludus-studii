package syncutil

import (
	"sync"
)

type Counter interface {
	Next(shouldReset func() bool) int64
}

type counter struct {
	mu    sync.Mutex
	count int64
}

func NewCounter() Counter {
	return &counter{}
}

func (c *counter) Next(shouldReset func() bool) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if shouldReset() {
		c.count = 0
	} else {
		c.count++
	}

	return c.count
}
