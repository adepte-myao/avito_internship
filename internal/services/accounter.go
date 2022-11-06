package services

import (
	"errors"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
)

type Accounter struct {
	Account  storage.Account
	TxHelper storage.SQLTransactionHelper
}

func NewAccounter(repo *storage.SQLRepository) *Accounter {
	return &Accounter{
		Account:  repo.Account,
		TxHelper: repo.SQLTransactionHelper,
	}
}

func (serv *Accounter) GetBalance(accountId int32) (models.Account, error) {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return models.Account{}, err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	Account, err := serv.Account.GetAccount(tx, accountId)
	if err != nil {
		return models.Account{}, errors.New("Account with given ID does not exist")
	}

	serv.TxHelper.CommitTransaction(tx)

	return Account, nil
}

func (serv *Accounter) Deposit(accountId int32, value decimal.Decimal) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	_, err = serv.Account.GetAccount(tx, accountId)
	if err != nil {
		// TODO: can't be other errors except no Account?
		err := serv.Account.CreateAccount(tx, accountId)
		if err != nil {
			return err
		}
	}

	err = serv.Account.IncreaseBalance(tx, accountId, value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Accounter) Withdraw(accountId int32, value decimal.Decimal) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	Account, err := serv.Account.GetAccount(tx, accountId)
	if err != nil {
		return errors.New("account does not exist")
	}

	if Account.Balance.LessThan(value) {
		return errors.New("not enough money")
	}

	err = serv.Account.DecreaseBalance(tx, accountId, value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}
