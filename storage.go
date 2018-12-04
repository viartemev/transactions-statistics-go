package transactions_statistics

import "sync"

type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

type ConcurrentMapItem struct {
	Key   string
	Value interface{}
}

func (cm *ConcurrentMap) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()

	cm.items[key] = value
}

func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
	cm.RLock()
	defer cm.RUnlock()

	value, ok := cm.items[key]

	return value, ok
}

func (cm *ConcurrentMap) Iter() <-chan ConcurrentMapItem {
	c := make(chan ConcurrentMapItem)

	go func() {
		cm.RLock()
		defer cm.RUnlock()

		for k, v := range cm.items {
			c <- ConcurrentMapItem{k, v}
		}
		close(c)
	}()

	return c
}
