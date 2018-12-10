package main

type ConcurrentMapWithChannels struct {
	setValCh chan ConcurrentMapItemWithChannels
	getValCh chan map[string]interface{}
	items    map[string]interface{}
}

type ConcurrentMapItemWithChannels struct {
	Key   string
	Value interface{}
}

func NewConcurrentMapWithChannels() ConcurrentMapWithChannels {
	h := ConcurrentMapWithChannels{
		setValCh: make(chan ConcurrentMapItemWithChannels),
		getValCh: make(chan map[string]interface{}),
	}
	go h.mux()
	return h
}

func (cm *ConcurrentMapWithChannels) mux() {
	var item ConcurrentMapItemWithChannels
	for {
		select {
		case item = <-cm.setValCh:
			cm.items[item.Key] = item.Value
		case cm.getValCh <- cm.items:
		}
	}
}

func (cm *ConcurrentMapWithChannels) Set(item *ConcurrentMapItemWithChannels) {
	cm.setValCh <- *item
}

func (cm *ConcurrentMapWithChannels) Get(key string) interface{} {
	items := <-cm.getValCh
	return items[key]
}

func (cm *ConcurrentMapWithChannels) Iter() <-chan ConcurrentMapItemWithChannels {
	c := make(chan ConcurrentMapItemWithChannels)
	go func() {
		for k, v := range <-cm.getValCh {
			c <- ConcurrentMapItemWithChannels{k, v}
		}
		close(c)
	}()
	return c
}
