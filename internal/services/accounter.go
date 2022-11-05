package services

import (
	"errors"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
)

type Accounter struct {
	repo *storage.SQLRepository
}

func NewAccounter(repo *storage.SQLRepository) *Accounter {
	return &Accounter{
		repo: repo,
	}
}

func (serv *Accounter) GetBalance(dto dtos.GetBalanceDto) (models.Account, error) {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return models.Account{}, err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	account, err := serv.repo.Account.GetAccount(tx, dto.AccountId)
	if err != nil {
		return models.Account{}, errors.New("account with given ID does not exist")
	}

	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return account, nil
}

func (serv *Accounter) Deposit(depDto dtos.DepositAccountDto) error {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	_, err = serv.repo.Account.GetAccount(tx, depDto.AccountId)
	if err != nil {
		// TODO: can't be other errors except no account?
		err := serv.repo.Account.CreateAccount(tx, depDto.AccountId)
		if err != nil {
			return err
		}
	}

	err = serv.repo.Account.IncreaseBalance(tx, depDto.AccountId, depDto.Value)
	if err != nil {
		return err
	}

	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return nil
}

func (serv *Accounter) Withdraw(wdDto dtos.WithdrawAccountDto) error {
	tx, err := serv.repo.SQLTransactionHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.repo.SQLTransactionHelper.RollbackTransaction(tx)

	account, err := serv.repo.Account.GetAccount(tx, wdDto.AccountId)
	if err != nil {
		return errors.New("account does not exist")
	}

	if account.Balance.LessThan(wdDto.Value) {
		return errors.New("not enough money")
	}

	err = serv.repo.Account.DecreaseBalance(tx, wdDto.AccountId, wdDto.Value)
	if err != nil {
		return err
	}

	serv.repo.SQLTransactionHelper.CommitTransaction(tx)

	return nil
}
