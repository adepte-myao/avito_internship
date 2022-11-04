package storage

import (
	"database/sql"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type AccountRepository struct {}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (repo *AccountRepository) GetAccount(tx *sql.Tx, id int32) (models.Account, error) {
	var acc models.Account
	err := tx.QueryRow(
		"SELECT id, balance FROM accounts WHERE id = $1", id,
	).Scan(&acc.ID, &acc.Balance)

	return acc, err
}

func (repo *AccountRepository) CreateAccount(tx *sql.Tx, id int32) error {
	_, err := tx.Exec(
		"INSERT INTO accounts (id, balance) VALUES ($1, $2)",
		id, 0,
	)
	return err
}

func (repo *AccountRepository) IncreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error {
	_, err := tx.Exec(
		`UPDATE accounts 
			SET balance = balance + $1 
			WHERE id = $2`,
		value, id,
	)
	return err
}

func (repo *AccountRepository) DecreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error {
	_, err := tx.Exec(
		`UPDATE accounts 
			SET balance = balance - $1 
			WHERE id = $2`,
		value, id,
	)
	return err
}
