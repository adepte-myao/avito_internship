package services

import (
	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
)

type Account interface {
	GetBalance(dto dtos.GetBalanceDto) (models.Account, error)
	Deposit(depDto dtos.DepositAccountDto) error
	Withdraw(wdDto dtos.WithdrawAccountDto) error
}

type Reservation interface {
	MakeReservation(resDto dtos.ReservationDto) error
	AcceptReservation(resDto dtos.ReservationDto) error
	CancelReservation(resDto dtos.ReservationDto) error
}

type Service struct {
	Account
	Reservation
}

func NewService(repo *storage.SQLRepository) *Service {
	return &Service{
		Reservation: NewReservationer(repo),
	}
}
