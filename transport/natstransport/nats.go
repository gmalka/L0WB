package natstransport

import (
	"context"
	"encoding/json"
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

			v := models.OrderStruct{}
			err = json.Unmarshal(msg.Data(), &v)
			if err != nil || v.OrderUID == nil {
				log.Printf("Incorrect input data: %v\n", string(msg.Data()))
				msg.Ack()
				break
			}

			data, err := json.Marshal(v)
			if err != nil {
				log.Printf("Incorrect input data: %v\n", string(msg.Data()))
				msg.Ack()
				break
			}

			err = n.s.Add(models.Order{
				OrderUID: *v.OrderUID,
				Order:    data,
			})
			if err != nil {
				log.Printf("cant add order: %v\n", err)
			}

			msg.Ack()
		}
	}
}
