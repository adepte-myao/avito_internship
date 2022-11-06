package main

import (
	"fmt"
	"os"

	"github.com/adepte-myao/avito_internship/internal/config"
	"github.com/adepte-myao/avito_internship/internal/handlers"
	"github.com/adepte-myao/avito_internship/internal/server"
	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/adepte-myao/avito_internship/internal/storage"
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

	db, err := storage.NewPostgresDb(&cfg.Store, logger)
	if err != nil {
		logger.Fatal(err)
		return
	}
	repository := storage.NewSQLRepository(db)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service, logger)
	router := handler.InitRoutes()

	server := server.NewServer(&cfg, logger, router)

	// Handlers registration

	err = server.Start()
	if err != nil {
		logger.Error(err)
		return
	}
}
