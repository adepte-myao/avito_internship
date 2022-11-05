package storage

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (repo *AccountRepository) GetAccount(tx *sql.Tx, id int32) (models.Account, error) {
	var acc models.Account
	var balance string
	err := tx.QueryRow(
		"SELECT id, balance FROM accounts WHERE id = $1", id,
	).Scan(&acc.ID, &balance)

	if err != nil {
		return models.Account{}, err
	}

	// balance format: "123.45 P"
	regBalance := regexp.MustCompile(`[^0-9,]`)
	balance = regBalance.ReplaceAllString(balance, "")

	balance = strings.Replace(balance, ",", ".", 1)

	acc.Balance, err = decimal.NewFromString(balance)
	if err != nil {
		return models.Account{}, err
	}

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
	if value.IsNegative() {
		return errors.New("the value to add cannot be negative")
	}

	// decimal uses '.' to separate fractional part
	// postgres in RUS local uses ','
	strValue := strings.Replace(value.String(), ".", ",", -1)

	_, err := tx.Exec(
		`UPDATE accounts 
			SET balance = balance + $1 
			WHERE id = $2`,
		strValue, id,
	)
	return err
}

func (repo *AccountRepository) DecreaseBalance(tx *sql.Tx, id int32, value decimal.Decimal) error {
	if value.IsNegative() {
		return errors.New("the value to sub cannot be negative")
	}

	// decimal uses '.' to separate fractional part
	// postgres in RUS local uses ','
	intValue := strings.Replace(value.String(), ".", ",", -1)

	_, err := tx.Exec(
		`UPDATE accounts 
			SET balance = balance - $1 
			WHERE id = $2`,
		intValue, id,
	)
	return err
}
