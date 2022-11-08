package services

import (
	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=service.go -destination=mock_service/mock.go

type Account interface {
	GetBalance(accountId int32) (models.Account, error)
	Deposit(accountId int32, value decimal.Decimal) error
	Withdraw(accountId int32, value decimal.Decimal) error
	InternalTransfer(senderId int32, recId int32, value decimal.Decimal) error
	GetStatement(dto dtos.GetAccountStatementDto) ([]models.StatementElem, error)
}

type Reservation interface {
	MakeReservation(resDto dtos.ReservationDto) error
	AcceptReservation(resDto dtos.ReservationDto) error
	CancelReservation(resDto dtos.ReservationDto) error
	GetAccountantReport(month int, year int) ([]models.AccountantReportElem, error)
}

type Service struct {
	Account
	Reservation
}

func NewService(repo *storage.SQLRepository) *Service {
	return &Service{
		Account:     NewAccounter(repo),
		Reservation: NewReservationer(repo),
	}
}
