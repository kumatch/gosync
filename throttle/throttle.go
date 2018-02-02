package throttle

import (
	"sync"
	"time"
)

type Group struct {
	mu     sync.Mutex
	locker map[string]struct{}
}

func (g *Group) Do(key string, fn func() (interface{}, error), wait time.Duration) (v interface{}, err error, invoked bool) {
	g.mu.Lock()
	if g.locker == nil {
		g.locker = make(map[string]struct{})
	}

	if _, ok := g.locker[key]; ok {
		g.mu.Unlock()
		return
	}

	g.locker[key] = struct{}{}
	g.mu.Unlock()

	defer func() {
		time.Sleep(wait)
		g.mu.Lock()
		delete(g.locker, key)
		g.mu.Unlock()
	}()

	v, err = fn()
	invoked = true
	return
}
