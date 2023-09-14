package postgresdb

import (
	"fmt"
	"l0wb/models"
	"l0wb/store/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type noteRepository struct {
	db *sqlx.DB
}

func NewPostgresDatabase(host, port, user, password, dbname, sslmode string) (database.Database, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode))

	if err != nil {
		return nil, fmt.Errorf("can't connect to bd: %v", err)
	}

	return &noteRepository{
		db: db,
	}, nil
}

func (n *noteRepository) Add(order models.Order) error {
	_, err := n.db.Exec("INSERT INTO Orders VALUES ($1,$2)",
		order.OrderUID, order.Order)
		
	if err != nil {
		return fmt.Errorf("can not insert order with id %s: %v", order.OrderUID, err)
	}

	return nil
}

func (n *noteRepository) Get(OrderUID string) (models.Order, error) {
	order := models.Order{}
	err := n.db.Get(&order, `SELECT * FROM Orders WHERE order_uid=$1`, OrderUID)
	if err != nil {
		return models.Order{}, fmt.Errorf("can not get order with id %s: %v", OrderUID, err)
	}

	fmt.Println(order)

	return order, nil
}

func (n *noteRepository) GetAll() ([]models.Order, error) {
	orders := []models.Order{}
	err := n.db.Select(&orders, `SELECT * FROM Orders`)
	if err != nil {
		return nil, fmt.Errorf("can not get all orders: %v", err)
	}

	return orders, nil
}
