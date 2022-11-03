package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID      int32           `json:"id"`
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

type CustomTransfer struct {
	ID             int64
	AccountId      int32
	OtherAccountId int32
	Type           TransferType
	Amount         decimal.Decimal
	RecordTime     time.Time
	BalanceAfter   decimal.Decimal
}
