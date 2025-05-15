package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron"
	"scrapper.go/internal/handler"
	"scrapper.go/internal/service"
	"scrapper.go/internal/storage/postgres"
	dbC "scrapper.go/pkg/db"
)

func main() {
	client, err := dbC.NewClient(context.Background(), 3, dbC.StorageConfig{
		Host:     "localhost",
		Port:     "5430",
		Username: "user",
		Password: "user",
		Database: "subscribe_db",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}

	currStorage := postgres.NewCurrencyRepository(client)

	pairStorage := postgres.NewPairRepository(client)

	go func() {
		c := cron.New()
		c.AddFunc("5 * * * *", func() {

			pairs, _ := pairStorage.GetAllPairs(context.Background())
			for _, pair := range pairs {
				log.Printf("Fetching rate for %s/%s", pair.Base, pair.Quote)
				rate, err := service.FetchRate(pair.Base, pair.Quote)
				if err != nil {
					fmt.Printf("Error fetching: %v", err)
					continue
				}
				fmt.Printf("Rate: %f", rate)
				if err := currStorage.SaveRate(context.Background(), int64(pair.ID), rate, time.Now()); err != nil {
					fmt.Printf("Error SaveRate %v", err)
				}
				if err := currStorage.DeleteOldRates(context.Background(), int64(pair.ID)); err != nil {
					fmt.Printf("Error clear old rates %v", err)
				}
			}
		})
		c.Start()
	}()

	router := httprouter.New()

	hand := handler.NewHandler(currStorage, pairStorage)
	hand.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	log.Print("Starting application")
	listener, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("Server is listening")
	log.Fatalln(server.Serve(listener))
}
