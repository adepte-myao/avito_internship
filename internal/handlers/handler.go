package handlers

import (
	"net/http"

	"github.com/adepte-myao/avito_internship/internal/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services *services.Service
	Logger   *logrus.Logger
}

func NewHandler(services *services.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		services: services,
		Logger:   logger,
	}
}

func (h *Handler) ping(rw http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Ping request received")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hello from balancer!"))
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", h.ping).Methods("GET")

	balance := router.PathPrefix("/balance").Subrouter()
	balance.HandleFunc("/", h.getBalance).Methods("GET")
	balance.HandleFunc("/deposit", h.deposit).Methods("POST")
	balance.HandleFunc("/withdraw", h.withdraw).Methods("POST")

	reservation := router.PathPrefix("/reservation").Subrouter()
	reservation.HandleFunc("/make", h.makeReservation).Methods("POST")
	reservation.HandleFunc("/accept", h.acceptReservation).Methods("POST")
	reservation.HandleFunc("/cancel", h.cancelReservation).Methods("POST")

	return router
}
