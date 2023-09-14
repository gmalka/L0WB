package orderservice

import (
	"fmt"
	"l0wb/models"
	"l0wb/store/cash"
	"l0wb/store/database"
)

type orderService struct {
	db   database.Database
	cash cash.Casher
}

type Orderer interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
}

func NewOrderService(db database.Database, cash cash.Casher) (Orderer, error) {
	orders, err := db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("can not create order service: %v", err)
	}

	for _, v := range orders {
		err = cash.Add(v)
		if err != nil {
			return nil, fmt.Errorf("can not create order service: %v", err)
		}
	}

	return orderService{
		db:   db,
		cash: cash,
	}, nil
}

func(o orderService) Add(order models.Order) error {
	err := o.cash.Add(order)
	if err != nil {
		return fmt.Errorf("can not add order to cash: %v", err)
	}

	err = o.db.Add(order)
	if err != nil {
		return fmt.Errorf("can not add order to store: %v", err)
	}

	return nil
}

func(o orderService) Get(OrderUID string) (models.Order, error) {
	order, err := o.cash.Get(OrderUID)
	if err == nil {
		return order, nil
	}

	order, err = o.db.Get(OrderUID)
	if err != nil {
		return models.Order{}, fmt.Errorf("can not get order: %v", err)
	}

	return order, nil
}