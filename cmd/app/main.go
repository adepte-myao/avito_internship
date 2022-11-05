package main

import (
	"fmt"
	"os"

	"github.com/adepte-myao/avito_internship/internal/config"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/server"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}
	defer f.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}

	logger := logrus.New()
	router := mux.NewRouter()

	db, err := storage.NewPostgresDb(&cfg.Store, logger)
	if err != nil {
		logger.Fatal(err)
		return
	}
	repository := storage.NewSQLRepository(db)
	service := services.NewService(repository)

	// Handlers initialization
	makeReservationHandler := handlers.NewMakeReservationHandler(logger, service)
	acceptReservationHandler := handlers.NewAcceptReservationHandler(logger, service)
	cancelReservationHandler := handlers.NewCancelReservationHandler(logger, service)

	getBalanceHandler := handlers.NewGetBalanceHandler(logger, service)
	depositHandler := handlers.NewDepositAccountHandler(logger, service)
	withdrawHandler := handlers.NewWithdrawAccountHandler(logger, service)

	server := server.NewServer(&cfg, logger, router)

	// Handlers registration
	server.RegisterHandler("/ping", server.Ping)

	server.RegisterHandler("/make-reservation", makeReservationHandler.Handle)
	server.RegisterHandler("/accept-reservation", acceptReservationHandler.Handle)
	server.RegisterHandler("/cancel-reservation", cancelReservationHandler.Handle)

	server.RegisterHandler("/balance", getBalanceHandler.Handle)
	server.RegisterHandler("/deposit", depositHandler.Handle)
	server.RegisterHandler("/withdraw", withdrawHandler.Handle)

	err = server.Start()
	if err != nil {
		logger.Error(err)
		return
	}
}
