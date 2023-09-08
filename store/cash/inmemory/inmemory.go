package inmemory

import (
	"fmt"
	"l0wb/models"
	"l0wb/store/cash"
)

type inmemory struct {
	store map[string]models.Order
}

func NewInmemoryCasher() cash.Casher {
	return inmemory{
		store: make(map[string]models.Order, 10),
	}
}

func (m inmemory) Add(order models.Order) error {
	if _, ok := m.store[order.OrderUID]; !ok {
		m.store[order.OrderUID] = order
		return nil
	} else {
		return fmt.Errorf("entry with key %s already exists", order.OrderUID)
	}
}

func (m inmemory) Get(OrderUID string) (models.Order, error) {
	if k, ok := m.store[OrderUID]; !ok {
		return models.Order{}, fmt.Errorf("could not find an entry with key %s in the cache", OrderUID)
	} else {
		return k, nil
	}
}

func (m inmemory) Delete(OrderUID string) error {
	if _, ok := m.store[OrderUID]; !ok {
		return fmt.Errorf("could not find an entry with key %s in the cache", OrderUID)
	}
	delete(m.store, OrderUID)
	return nil
}
