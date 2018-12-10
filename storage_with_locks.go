package main

import "sync"

type ConcurrentMapWithLocks struct {
	sync.RWMutex
	items map[string]interface{}
}

type ConcurrentMapItemWithLocks struct {
	Key   string
	Value interface{}
}

func (cm *ConcurrentMapWithLocks) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()

	cm.items[key] = value
}

func (cm *ConcurrentMapWithLocks) Get(key string) (interface{}, bool) {
	cm.RLock()
	defer cm.RUnlock()

	value, ok := cm.items[key]

	return value, ok
}

func (cm *ConcurrentMapWithLocks) Iter() <-chan ConcurrentMapItemWithLocks {
	c := make(chan ConcurrentMapItemWithLocks)

	go func() {
		cm.RLock()
		defer cm.RUnlock()

		for k, v := range cm.items {
			c <- ConcurrentMapItemWithLocks{k, v}
		}
		close(c)
	}()

	return c
}
