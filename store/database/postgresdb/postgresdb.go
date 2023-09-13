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
	_, err := n.db.Exec("INSERT INTO Orders VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		return fmt.Errorf("can not insert order with id %s: %v", order.OrderUID, err)
	}

	_, err = n.db.Exec("INSERT INTO Deliveries VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		order.Delivery.Name, order.OrderUID, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("can not insert delivery with name %s: %v", order.Delivery.Name, err)
	}

	_, err = n.db.Exec("INSERT INTO Payments VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		order.Payment.Transaction, order.OrderUID, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("can not insert payment with transaction %s: %v", order.Payment.Transaction, err)
	}

	for _, v := range order.Items {
		_, err = n.db.Exec("INSERT INTO Items VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
			v.ChrtID, order.OrderUID, v.TrackNumber, v.Price, v.Rid, v.Name, v.Sale, v.Size, v.TotalPrice, v.NmID, v.Brand, v.Status)
		if err != nil {
			return fmt.Errorf("can not insert item with id %d: %v", v.ChrtID, err)
		}
	}

	return nil
}

func (n *noteRepository) Get(OrderUID string) (models.Order, error) {
	order := models.Order{}
	err := n.db.Select(&order, `SELECT * FROM Orders LEFT JOIN
		Deliveries d ON o.order_uid = d.order_id
	LEFT JOIN
		Payments p ON o.order_uid = p.order_id
	LEFT JOIN
		Items i ON o.order_uid = i.order_id
	WHERE order_uid=$1`, OrderUID)
	if err != nil {
		return models.Order{}, fmt.Errorf("can not get order with id %s: %v", OrderUID, err)
	}

	return order, nil
}

func (n *noteRepository) GetAll() ([]models.Order, error) {
	orders := []models.Order{}
	err := n.db.Select(&orders, `SELECT * FROM Orders LEFT JOIN
		Deliveries d ON o.order_uid = d.order_id
	LEFT JOIN
		Payments p ON o.order_uid = p.order_id
	LEFT JOIN
		Items i ON o.order_uid = i.order_id`)
	if err != nil {
		return nil, fmt.Errorf("can not get all orders: %v", err)
	}

	return orders, nil
}
