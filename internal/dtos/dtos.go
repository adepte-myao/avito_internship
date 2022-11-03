package dtos

import "github.com/shopspring/decimal"

type ReservationDto struct {
	AccountId int32           `json:"accountId"`
	ServiceId int32           `json:"serviceId"`
	OrderId   int32           `json:"orderId"`
	TotalCost decimal.Decimal `json:"totalCost"`
}

type DepositAccountDto struct {
	AccountId int32           `json:"accountId"`
	Value     decimal.Decimal `json:"value"`
}

type GetAccountStatementDto struct {
	AccountId int32 `json:"accountId"`
}

type GetBalanceDto struct {
	AccountId int32 `json:"accountId"`
}

type MakeAccountantReportDto struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type MakeInternalTransferDto struct {
	SenderId   int32           `json:"senderId"`
	ReceiverId int32           `json:"receiverId"`
	Value      decimal.Decimal `json:"value"`
}

type WithdrawAccountDto struct {
	AccountId int32           `json:"accountId"`
	Value     decimal.Decimal `json:"value"`
}
