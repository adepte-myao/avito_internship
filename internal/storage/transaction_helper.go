package storage

import "database/sql"

type TransactionHelper struct {
	storage *Storage
}

func NewTransactionHelper(storage *Storage) *TransactionHelper {
	return &TransactionHelper{
		storage: storage,
	}
}

func (helper *TransactionHelper) BeginTransaction() (*sql.Tx, error) {
	return helper.storage.db.Begin()
}

func (helper *TransactionHelper) RollbackTransaction(tx *sql.Tx) {
	tx.Rollback()
}

func (helper *TransactionHelper) CommitTransaction(tx *sql.Tx) {
	tx.Commit()
}
