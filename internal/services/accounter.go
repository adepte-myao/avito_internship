package services

import (
	"errors"

	"github.com/adepte-myao/avito_internship/internal/dtos"
	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
)

type Accounter struct {
	Account  storage.Account
	Transfer storage.Transfer
	TxHelper storage.SQLTransactionHelper
}

func NewAccounter(repo *storage.SQLRepository) *Accounter {
	return &Accounter{
		Account:  repo.Account,
		Transfer: repo.Transfer,
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

	err = serv.Transfer.RecordExternalTransfer(tx, accountId, models.Deposit, value)
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

	err = serv.Transfer.RecordExternalTransfer(tx, accountId, models.Withdraw, value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Accounter) InternalTransfer(senderId int32, recId int32, value decimal.Decimal) error {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	senderAccount, err := serv.Account.GetAccount(tx, senderId)
	if err != nil {
		return errors.New("sender account does not exist")
	}

	_, err = serv.Account.GetAccount(tx, recId)
	if err != nil {
		return errors.New("receiver account does not exist")
	}

	if senderAccount.Balance.LessThan(value) {
		return errors.New("not enough money")
	}

	err = serv.Account.DecreaseBalance(tx, senderId, value)
	if err != nil {
		return err
	}

	err = serv.Account.IncreaseBalance(tx, recId, value)
	if err != nil {
		return err
	}

	err = serv.Transfer.RecordInternalTransfer(tx, senderId, recId, value)
	if err != nil {
		return err
	}

	serv.TxHelper.CommitTransaction(tx)

	return nil
}

func (serv *Accounter) GetStatement(dto dtos.GetAccountStatementDto) ([]models.StatementElem, error) {
	tx, err := serv.TxHelper.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer serv.TxHelper.RollbackTransaction(tx)

	// TODO: if there are no statements, better to inform user about it
	statements, err := serv.Transfer.GetAccountStatements(tx, dto)
	if err != nil {
		return nil, err
	}

	serv.TxHelper.CommitTransaction(tx)

	return statements, nil
}
