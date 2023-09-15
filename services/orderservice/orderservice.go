package orderservice

import (
	"fmt"
	"l0wb/models"
	"l0wb/store/cash"
	"l0wb/store/database"
	"log"
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
	log.Printf("Added order with id %v in cash", order.OrderUID)

	err = o.db.Add(order)
	if err != nil {
		return fmt.Errorf("can not add order to store: %v", err)
	}
	log.Printf("Added order with id %v in store", order.OrderUID)

	return nil
}

func(o orderService) Get(OrderUID string) (models.Order, error) {
	order, err := o.cash.Get(OrderUID)
	if err == nil {
		log.Printf("Returned order with id %v from cash", order.OrderUID)
		return order, nil
	}
	log.Printf("Cant find order with id %v in cash", order.OrderUID)

	order, err = o.db.Get(OrderUID)
	if err != nil {
		log.Printf("Cant find order with id %v in store", order.OrderUID)
		return models.Order{}, fmt.Errorf("can not get order: %v", err)
	}

	log.Printf("Returned order with id %v from store", order.OrderUID)
	return order, nil
}