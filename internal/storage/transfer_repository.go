package storage

import (
	"database/sql"
	"strings"
	"time"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type Transferer struct{}

func NewTransferer() *Transferer {
	return &Transferer{}
}

func (repo *Transferer) RecordExternalTransfer(tx *sql.Tx, accId int32, ttype models.TransferType, amount decimal.Decimal) error {
	_, err := tx.Exec(
		`INSERT INTO external_transfers_history (account_id, transfer_type, amount, record_time)
			VALUES ($1, $2, $3, $4)`,
		accId,
		ttype.String(),
		strings.Replace(amount.String(), ".", ",", -1),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Transferer) RecordInternalTransfer(tx *sql.Tx, senderId int32, recId int32, amount decimal.Decimal) error {
	_, err := tx.Exec(
		`INSERT INTO internal_transfers_history (sender_id, receiver_id, amount, record_time)
			VALUES ($1, $2, $3, $4)`,
		senderId,
		recId,
		strings.Replace(amount.String(), ".", ",", -1),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
