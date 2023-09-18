package natstransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"l0wb/models"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsHandler struct {
	s Orderer
}

type Orderer interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
}

func NewNatsHasher(s Orderer) NatsHandler {
	return NatsHandler{
		s: s,
	}
}

func (n NatsHandler) RunNats(ctx context.Context, url string) error {
	var err error
	var nc *nats.Conn

	cont, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nc, err = nats.Connect(url)
	for err != nil {
		select {
		case <-cont.Done():
			if err != nil {
				return err
			}
		default:
			nc, err = nats.Connect(url)
		}
		time.Sleep(time.Second)
	}
	defer nc.Close()

	js, _ := jetstream.New(nc)

	_, err = js.CreateStream(cont, jetstream.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"ORDERS.*"},
		Retention: jetstream.WorkQueuePolicy,
	})

	if err != nil {
		return err
	}

	cons, err := js.CreateOrUpdateConsumer(cont, "ORDERS", jetstream.ConsumerConfig{
		Name:    "my_consumer",
		Durable: "my_consumer",
	})
	if err != nil {
		return err
	}

	iter, _ := cons.Messages(jetstream.PullMaxMessages(1))
	defer iter.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := iter.Next()
			if err != nil {
				log.Printf("stop?! %s\n", err)
				return nil
			}

			v := map[string]interface{}{}
			err = json.Unmarshal(msg.Data(), &v)
			if err != nil {
				log.Printf("Incorrect input data: %v\n", string(msg.Data()))
				msg.Ack()
				break
			}

			err = validateOrder(v)

			if err != nil {
				log.Printf("Incorrect input data: %v\n", err)
				msg.Ack()
				break
			}

			err = n.s.Add(models.Order{
				OrderUID: v["order_uid"].(string),
				Order:    msg.Data(),
			})
			if err != nil {
				log.Printf("cant add order: %v\n", err)
			}

			msg.Ack()
		}
	}
}

func validateOrder(v map[string]interface{}) error {
	if k, ok := v["items"]; !ok {
		return errors.New("cant find items in struct")
	} else {
		err := validateItems(k)
		if err != nil {
			return err
		}
	}

	if k, ok := v["delivery"]; !ok {
		return errors.New("cant find delivery in struct")
	} else {
		err := validateDelivery(k)
		if err != nil {
			return err
		}
	}

	if k, ok := v["payment"]; !ok {
		return errors.New("cant find payment in struct")
	} else {
		err := validatePayment(k)
		if err != nil {
			return err
		}
	}

	// count := 0
	// structType := reflect.TypeOf(models.OrderStruct{})

	// for i := 0; i < structType.NumField(); i++ {
	// 	field := structType.Field(i)

	// 	// Получите имя и тип поля
	// 	fieldName := field.Tag.Get("json")
	// 	fieldType := field.Type

	// 	if k, ok := v[fieldName]; !ok {
	// 		return fmt.Errorf("cant find field named %v", fieldName)
	// 	} else {
	// 		val := reflect.TypeOf(k)
	// 		fmt.Println(k)

	// 		if val != fieldType {
	// 			return fmt.Errorf("incorrect type: want %v, but got %v", val, fieldType)
	// 		}

	// 		fmt.Printf("%v - %v ::: %v\n", field.Name, val, fieldType)
	// 	}

	// 	count++
	// }

	val := v["l"]
	switch val.(type) {
	case float64, int:
		k := "false"
		fmt.Println("1:", k)
	case string:
		fmt.Println("ПОЧЕМУ СТРИНГ БЛЯТЬ?!")
	default:
		fmt.Println("2: неизвестный тип")
	}

	return nil
}

func validateItems(v interface{}) error {
	return nil
}

func validateDelivery(v interface{}) error {
	return nil
}

func validatePayment(v interface{}) error {
	return nil
}
