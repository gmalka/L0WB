package database

import "l0wb/models"

type Database interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
	GetAll() ([]models.Order, error)
}