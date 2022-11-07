package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID      int32           `json:"accountId"`
	Balance decimal.Decimal `json:"balance"`
}

type Service struct {
	ID   int32
	Name string
}

type Reservation struct {
	ID           int32
	AccountId    int32
	ServiceId    int32
	OrderId      int32
	TotalCost    decimal.Decimal
	State        ReserveState
	RecordTime   time.Time
	BalanceAfter decimal.Decimal
}

type StatementElem struct {
	RecordTime   time.Time       `json:"recordTime"`
	TransferType string          `json:"transferType"`
	Amount       decimal.Decimal `json:"amount"`
	Description  string          `json:"description"`
}

type AccountantReportElem struct {
	ServiceName   string          `json:"serviceName"`
	TotalReceived decimal.Decimal `json:"totalReceived"`
}
