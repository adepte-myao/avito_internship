package services

import (
	"errors"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
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

func (serv *Accounter) GetBalance(dto dtos.GetBalanceDto) (models.Account, error) {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return models.Account{}, err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	Account, err := serv.Account.GetAccount(tx, dto.AccountId)
	if err != nil {
		return models.Account{}, errors.New("Account with given ID does not exist")
	}

	serv.TxHelper.CommitTransaction(tx)

	return Account, nil
}

func (serv *Accounter) Deposit(depDto dtos.DepositAccountDto) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	_, err = serv.Account.GetAccount(tx, depDto.AccountId)
	if err != nil {
		// TODO: can't be other errors except no Account?
		err := serv.Account.CreateAccount(tx, depDto.AccountId)
		if err != nil {
			return err
		}
	}

	err = serv.Account.IncreaseBalance(tx, depDto.AccountId, depDto.Value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Accounter) Withdraw(wdDto dtos.WithdrawAccountDto) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	Account, err := serv.Account.GetAccount(tx, wdDto.AccountId)
	if err != nil {
		return errors.New("Account does not exist")
	}

	if Account.Balance.LessThan(wdDto.Value) {
		return errors.New("not enough money")
	}

	err = serv.Account.DecreaseBalance(tx, wdDto.AccountId, wdDto.Value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}
