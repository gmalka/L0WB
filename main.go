package main

import (
	"context"
	"fmt"
	"l0wb/store/cash/inmemory"
	"l0wb/store/database/postgresdb"
	"l0wb/transport/natstransport"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	db, err := postgresdb.NewPostgresDatabase("localjost", "5432", "postgres", "postgres", "postgres", "disable")
	if err != nil {
		log.Printf("cant connect db: %v\n", err)
		return
	}

	cash := inmemory.NewInmemoryCasher()
	n := natstransport.NewNatsHasher(cash, db)

	go n.RunNats(ctx, nats.DefaultURL)

	for {
		select {
		case <-sig:
			fmt.Println("Server down")
			cancel()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}