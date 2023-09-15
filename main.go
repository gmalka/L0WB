package main

import (
	"context"
	"fmt"
	"l0wb/services/orderservice"
	"l0wb/store/cash/inmemory"
	"l0wb/store/database/postgresdb"
	"l0wb/transport/natstransport"
	"l0wb/transport/resttransport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	db, err := postgresdb.NewPostgresDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_TABLE"), os.Getenv("DB_SSLMODE"))
	if err != nil {
		log.Printf("cant connect db: %v\n", err)
		return
	}

	cash := inmemory.NewInmemoryCasher()

	orderservice, err := orderservice.NewOrderService(db, cash)
	if err != nil {
		log.Printf("cant connect db: %v\n", err)
		return
	}

	n := natstransport.NewNatsHasher(orderservice)

	go func() {
		err := n.RunNats(ctx, fmt.Sprintf("nats://%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")))
		if err != nil {
			log.Printf("Nats error: %v\n", err)
			sig <- syscall.SIGINT
		}
	}()

	h := resttransport.NewHandler(orderservice)

	serv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", os.Getenv("URL"), os.Getenv("PORT")),
		Handler: h.Init(),
	}

	go func() {
		err := serv.ListenAndServe()
		if err != nil {
			log.Printf("Nats error: %v\n", err)
			sig <- syscall.SIGINT
		}
	}()

	for {
		select {
		case <-sig:
			fmt.Println("Server down")
			serv.Shutdown(context.TODO())
			cancel()
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// docker run --name nats-server --rm -p 4222:4222 synadia/nats-server:nightly -js
