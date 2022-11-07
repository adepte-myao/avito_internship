package storage

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type TransferRepository struct{}

func NewTransferRepository() *TransferRepository {
	return &TransferRepository{}
}

func (repo *TransferRepository) RecordExternalTransfer(tx *sql.Tx, accId int32, ttype models.TransferType, amount decimal.Decimal) error {
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

func (repo *TransferRepository) RecordInternalTransfer(tx *sql.Tx, senderId int32, recId int32, amount decimal.Decimal) error {
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

func (repo *TransferRepository) GetAccountStatements(tx *sql.Tx, dto dtos.GetAccountStatementDto) ([]models.StatementElem, error) {
	var sortCriteria string
	if dto.FirstSortCriteria != "record_time" && dto.FirstSortCriteria != "amount" {
		sortCriteria = "record_time"
	} else {
		sortCriteria = dto.FirstSortCriteria
	}

	if dto.SecondSortCriteria != dto.FirstSortCriteria &&
		(dto.SecondSortCriteria == "record_time" || dto.SecondSortCriteria == "amount") {
		sortCriteria += ", " + dto.SecondSortCriteria
	}

	rows, err := tx.Query(
		`SELECT * 
		FROM (SELECT record_time, 'deposit'::transfer_type as "transfer_type", total_cost as "amount",
			concat('service: ', s.name, ', order ID: ', order_id) as "description"
			FROM reserves_history 
			JOIN services s on reserves_history.service_id = s.id
			WHERE account_id = $1 AND state = 'cancelled'::reserve_state
		UNION ALL
		SELECT record_time, 'withdraw'::transfer_type as "transfer_type", total_cost as "amount",
       		concat('service: ', s.name, ', order ID: ', order_id) as "description"
    		FROM reserves_history 
			JOIN services s on reserves_history.service_id = s.id
			WHERE account_id = $1 AND state = 'reserved'::reserve_state
		UNION ALL
		SELECT record_time, transfer_type, amount, 'external transfer' as "description"
			FROM external_transfers_history
			WHERE account_id = $1
		UNION ALL
		SELECT record_time, 'deposit'::transfer_type, amount,
			concat('internal transfer from ', sender_id) as "description"
			FROM internal_transfers_history
			WHERE receiver_id = $1
		UNION ALL
		SELECT record_time, 'withdraw'::transfer_type, amount,
			concat('internal transfer to ', receiver_id) as "description"
			FROM internal_transfers_history
			WHERE sender_id = $1
		) AS t
		ORDER BY $2 DESC
		LIMIT $3
		OFFSET $4`,
		dto.AccountId, sortCriteria, dto.PageSize, (dto.Page-1)*dto.PageSize)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statements := make([]models.StatementElem, 0)
	for rows.Next() {
		statement := models.StatementElem{}
		var balance string
		err = rows.Scan(&statement.RecordTime, &statement.TransferType, &balance, &statement.Description)
		if err != nil {
			return nil, err
		}

		// balance format: "123.45 P"
		regBalance := regexp.MustCompile(`[^0-9,]`)
		balance = regBalance.ReplaceAllString(balance, "")
		balance = strings.Replace(balance, ",", ".", 1)
		statement.Amount, err = decimal.NewFromString(balance)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statements, nil
}
