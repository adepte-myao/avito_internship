package storage

import "database/sql"

type TransactionHelper struct {
	db *sql.DB
}

func NewTransactionHelper(db *sql.DB) *TransactionHelper {
	return &TransactionHelper{
		db: db,
	}
}

func (helper *TransactionHelper) BeginTransaction() (*sql.Tx, error) {
	return helper.db.Begin()
}

func (helper *TransactionHelper) RollbackTransaction(tx *sql.Tx) {
	tx.Rollback()
}

func (helper *TransactionHelper) CommitTransaction(tx *sql.Tx) {
	tx.Commit()
}
