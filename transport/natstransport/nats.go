package natstransport

import (
	"context"
	"fmt"
	"l0wb/store/cash"
	"l0wb/store/database"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsHandler struct {
	casher cash.Casher
	store  database.Database
}

func NewNatsHasher(casher cash.Casher, store database.Database) NatsHandler {
	return NatsHandler{
		casher: casher,
		store: store,
	}
}

func (n NatsHandler)RunNats(ctx context.Context, url string) error {
	cont, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	nc, _ := nats.Connect(url)

	defer nc.Close()

	// Create a JetStream management interface
	js, _ := jetstream.New(nc)

	// Create a stream
	_, err := js.CreateStream(cont, jetstream.StreamConfig{
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
				break
			}

			fmt.Println(string(msg.Data()))
			msg.Ack()
		}
	}
}
