package inmemory

import (
	"fmt"
	"l0wb/models"
	"l0wb/store/cash"
	"sync"
)

type inmemory struct {
	store map[string][]byte
	m *sync.Mutex
}

func NewInmemoryCasher() cash.Casher {
	return inmemory{
		store: make(map[string][]byte, 10),
		m: &sync.Mutex{},
	}
}

func (m inmemory) Add(order models.Order) error {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[order.OrderUID]; !ok {
		m.store[order.OrderUID] = order.Order
		return nil
	} else {
		return fmt.Errorf("entry with key %s already exists", order.OrderUID)
	}
}

func (m inmemory) Get(OrderUID string) (models.Order, error) {
	m.m.Lock()
	defer m.m.Unlock()
	if k, ok := m.store[OrderUID]; !ok {
		return models.Order{}, fmt.Errorf("could not find an entry with key %s in the cache", OrderUID)
	} else {
		return models.Order{
			OrderUID: OrderUID,
			Order: k,
		}, nil
	}
}