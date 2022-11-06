package storage_test

import (
	"database/sql"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/models"
	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func getTransferRepoTxCommitFunc(t *testing.T) (*storage.TransferRepository, *sql.Tx, func(*sql.Tx)) {
	storageStruct, closeFunc := storage.TestStore(t, databaseURL)

	repo := storage.NewTransferRepository()
	txHelper := storage.NewTransactionHelper(storageStruct)

	tx, err := txHelper.BeginTransaction()
	assert.NoError(t, err)

	accRepo := storage.NewAccountRepository()
	var accountId int32 = 15
	err = accRepo.CreateAccount(tx, accountId)
	assert.NoError(t, err)
	err = accRepo.CreateAccount(tx, accountId+1)
	assert.NoError(t, err)

	return repo, tx, func(tx *sql.Tx) {
		txHelper.CommitTransaction(tx)
		closeFunc()
	}
}

func TestTransferRepository_RecordExternalTransfer_Success(t *testing.T) {
	repo, tx, commitFunc := getTransferRepoTxCommitFunc(t)
	defer commitFunc(tx)

	err := repo.RecordExternalTransfer(tx, 15, models.Deposit, decimal.NewFromInt(100))
	assert.NoError(t, err)
}

func TestTransferRepository_RecordInternalTransfer_Success(t *testing.T) {
	repo, tx, commitFunc := getTransferRepoTxCommitFunc(t)
	defer commitFunc(tx)

	err := repo.RecordInternalTransfer(tx, 15, 16, decimal.NewFromInt(100))
	assert.NoError(t, err)
}
