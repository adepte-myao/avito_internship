package main

import (
	"fmt"
	"os"

	"github.com/adepte-myao/avito_internship/internal/config"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/server"
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
	storage := storage.NewStorage(&cfg.Store, logger)
	if err = storage.Open(); err != nil {
		logger.Error(err)
		return
	}

	// Handlers initialization
	makeReservationHandler := handlers.NewMakeReservationHandler(logger, storage)
	acceptReservationHandler := handlers.NewAcceptReservationHandler(logger, storage)
	cancelReservationHandler := handlers.NewCancelReservationHandler(logger, storage)

	getBalanceHandler := handlers.NewGetBalanceHandler(logger, storage)
	depositHandler := handlers.NewDepositAccountHandler(logger, storage)
	withdrawHandler := handlers.NewWithdrawAccountHandler(logger, storage)

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
