package syncutil

import (
	"sync"
)

type Counter interface {
	Next(shouldReset func() bool) int64
}

type defaultCounter struct {
	mu    sync.Mutex
	count int64
}

func NewCounter() Counter {
	return &defaultCounter{}
}

func (c *defaultCounter) Next(shouldReset func() bool) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if shouldReset() {
		c.count = 0
	} else {
		c.count++
	}

	return c.count
}
