package cash

import "l0wb/models"

type Casher interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
}