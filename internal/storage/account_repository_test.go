package storage_test

import (
	"database/sql"
	"testing"

	"github.com/adepte-myao/avito_internship/internal/storage"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func getAccountRepoTxCommitFunc(t *testing.T) (*storage.AccountRepository, *sql.Tx, func(*sql.Tx)) {
	storageStruct, closeFunc := storage.TestStore(t, databaseURL)

	repo := storage.NewAccountRepository()
	txHelper := storage.NewTransactionHelper(storageStruct)

	tx, err := txHelper.BeginTransaction()
	assert.NoError(t, err)

	return repo, tx, func(tx *sql.Tx) {
		txHelper.CommitTransaction(tx)
		closeFunc()
	}
}

func TestAccountRepository_CreateAccount(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 15
	err := repo.CreateAccount(tx, accountId)

	assert.NoError(t, err)
}

func TestAccountRepository_GetAccount_Success(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 16
	repo.CreateAccount(tx, accountId)

	account, err := repo.GetAccount(tx, accountId)

	assert.NoError(t, err)
	assert.NotNil(t, account)
	if account.ID != 16 || !account.Balance.Equal(decimal.NewFromInt(0)) {
		t.Fatal("account values must be: 16, 0; given:", account.ID, ",", account.Balance)
	}
}

func TestAccountRepository_GetAccount_FailOnWrongID(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 17
	var wrongAccountId int32 = 18
	repo.CreateAccount(tx, accountId)

	_, err := repo.GetAccount(tx, wrongAccountId)

	assert.Error(t, err)
}

func TestAccountRepository_IncreaseBalance_Success(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 19
	repo.CreateAccount(tx, accountId)

	valueToAdd := decimal.NewFromFloat(12345.67)
	err := repo.IncreaseBalance(tx, accountId, valueToAdd)

	assert.NoError(t, err)

	// Check balance
	account, err := repo.GetAccount(tx, accountId)
	assert.NoError(t, err)

	areEqual := account.Balance.Equal(valueToAdd) // Start balance was 0
	assert.True(t, areEqual)
}

func TestAccountRepository_IncreaseBalance_FailOnNegativeValue(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 20
	repo.CreateAccount(tx, accountId)
	valueToAdd := decimal.NewFromFloat(-12345.67)

	err := repo.IncreaseBalance(tx, accountId, valueToAdd)

	assert.Error(t, err)

	// Check balance
	account, err := repo.GetAccount(tx, accountId)

	areEqual := account.Balance.Equal(decimal.NewFromInt(0))
	assert.True(t, areEqual)

	assert.NoError(t, err)
}

func TestAccountRepository_DecreaseBalance_Success(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 21
	repo.CreateAccount(tx, accountId)
	valueToSub := decimal.NewFromFloat(12345.67)
	repo.IncreaseBalance(tx, accountId, valueToSub)

	err := repo.DecreaseBalance(tx, accountId, valueToSub)

	assert.NoError(t, err)

	// Check balance
	account, err := repo.GetAccount(tx, accountId)
	assert.NoError(t, err)

	// Must be 0 because we sub the same value
	areEqual := account.Balance.Equal(decimal.NewFromInt(0))
	assert.True(t, areEqual)
}

func TestAccountRepository_DecreaseBalance_FailOnNegativeValue(t *testing.T) {
	repo, tx, commitFunc := getAccountRepoTxCommitFunc(t)
	defer commitFunc(tx)

	var accountId int32 = 22
	repo.CreateAccount(tx, accountId)
	valueToSub := decimal.NewFromFloat(-12345.67)
	repo.IncreaseBalance(tx, accountId, valueToSub)

	err := repo.DecreaseBalance(tx, accountId, valueToSub)
	assert.Error(t, err)
}
