package main

import (
	"fmt"
	"os"

	"github.com/adepte-myao/avito_internship/internal/config"
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

	server := server.NewServer(&cfg, logger, router)

	server.RegisterHandler("/ping", server.Ping)

	err = server.Start()
	if err != nil {
		logger.Error(err)
		return
	}
}
