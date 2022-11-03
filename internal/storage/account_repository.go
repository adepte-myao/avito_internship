package storage

import (
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/shopspring/decimal"
)

type AccountRepository struct {
	storage *Storage
}

func NewAccountRepository(storage *Storage) *AccountRepository {
	return &AccountRepository{
		storage: storage,
	}
}

func (repo *AccountRepository) GetAccount(id int32) (models.Account, error) {
	var acc models.Account
	err := repo.storage.db.QueryRow(
		"SELECT id, balance FROM accounts WHERE id = $1", id,
	).Scan(&acc.ID, &acc.Balance)

	return acc, err
}

func (repo *AccountRepository) CreateAccount(id int32) error {
	_, err := repo.storage.db.Exec(
		"INSERT INTO accounts (id, balance) VALUES ($1, $2)",
		id, 0,
	)
	return err
}

func (repo *AccountRepository) IncreaseBalance(id int32, value decimal.Decimal) error {
	_, err := repo.storage.db.Exec(
		`UPDATE accounts 
			SET balance = balance + $1 
			WHERE id = $2`,
		value, id,
	)
	return err
}

func (repo *AccountRepository) DecreaseBalance(id int32, value decimal.Decimal) error {
	_, err := repo.storage.db.Exec(
		`UPDATE accounts 
			SET balance = balance - $1 
			WHERE id = $2`,
		value, id,
	)
	return err
}
