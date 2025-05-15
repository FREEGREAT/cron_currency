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
	"github.com/spf13/viper"
	"scrapper.go/internal/handler"
	"scrapper.go/internal/service"
	"scrapper.go/internal/storage/postgres"
	postgreConnect "scrapper.go/pkg/postgreSQL"
	"scrapper.go/pkg/utils"
)

func main() {

	if err := utils.InitConfig(); err != nil {
		panic("Failed to load config")
	}

	postgresqlClient, err := postgreConnect.NewClient(context.TODO(), viper.GetInt("postgre.connectAttempt"), postgreConnect.StorageConfig{
		Host:     viper.GetString("postgre.host"),
		Port:     viper.GetString("postgre.port"),
		Username: viper.GetString("postgre.username"),
		Password: viper.GetString("postgre.password"),
		Database: viper.GetString("postgre.dbname"),
		SSLMode:  viper.GetString("postgre.sslmode"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to db %v", err)
	}

	currStorage := postgres.NewCurrencyRepository(postgresqlClient)
	pairStorage := postgres.NewPairRepository(postgresqlClient)

	scrapService := service.NewScrapService()
	storageService := service.NewStorageService(currStorage, pairStorage)

	go cronJob(storageService, scrapService)

	router := httprouter.New()

	hand := handler.NewHandler(storageService)
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

func cronJob(storageService service.StorageService, scrapService service.ScrapService) {
	c := cron.New()
	c.AddFunc("5 * * * *", func() {

		pairs, _ := storageService.GetAllPairs(context.Background())
		for _, pair := range pairs {
			log.Printf("\n\nFetching rate for %s/%s", pair.Base, pair.Quote)
			rate, err := scrapService.FetchRate(pair.Base, pair.Quote)
			if err != nil {
				fmt.Printf("Error fetching: %v", err)
				continue
			}

			if err := storageService.SaveRate(context.Background(), int64(pair.ID), rate, time.Now()); err != nil {
				fmt.Printf("Error SaveRate %v", err)
			}
			if err := storageService.DeleteOldRates(context.Background(), int64(pair.ID)); err != nil {
				fmt.Printf("Error clear old rates %v", err)
			}
		}
	})
	c.Start()

}
