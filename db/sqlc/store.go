package db

import (
	"context"
	"database/sql"
	"fmt"
)

// The Store type contains a pointer to a Queries struct and a pointer to a sql.DB struct.
// @property {Queries}  - The `Store` struct has two properties:
// @property db - The `db` property is a pointer to a `sql.DB` object, which represents a database
// connection pool. It is used to execute SQL queries and interact with the database.
type Store struct {
	*Queries
	db *sql.DB
}

// The function creates a new instance of a Store struct with a given database connection and
// associated queries.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// This function `execTx` is used to execute a function within a database transaction. It takes a
// context and a function as input parameters. The function parameter is a function that takes a
// `*Queries` object as input and returns an error. The `*Queries` object is used to execute database
// queries within the transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	// run query within transaction
	q := New(tx)
	err = fn(q)

	// rollback transaction
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v\nrb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// The TransferTxParams type defines parameters for a transfer transaction, including the IDs of the
// accounts involved and the amount being transferred.
// @property {int64} FromAccountID - FromAccountID is an integer that represents the ID of the account
// from which the transfer is being made.
// @property {int64} ToAccountID - ToAccountID is an integer property that represents the unique
// identifier of the account to which the transfer is being made.
// @property {int64} Amount - The `Amount` property is an integer that represents the amount of a
// currency being transferred from one account to another. It could be a positive or negative value
// depending on whether the transfer is a deposit or a withdrawal.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// The TransferTxResult type represents the result of a transfer transaction, including information
// about the transfer, the sender and receiver accounts, and the amount transferred.
// @property {Transfer} Transfer - The Transfer property is of type Transfer and represents the
// transfer object that was created as a result of the transaction. It contains information such as the
// sender, recipient, and amount transferred.
// @property {Account} FromAccount - FromAccount is a property of the TransferTxResult struct that
// represents the account from which the transfer was made. It is of type Account, which likely
// contains information such as the account holder's name, account number, and balance.
// @property {Account} ToAccount - ToAccount is a property of the TransferTxResult struct and
// represents the account that received the transfer in a transaction. It is of type Account, which
// likely contains information such as the account holder's name, account number, and balance.
// @property {Entry} FromEntry - FromEntry is a property of the TransferTxResult struct that represents
// the entry (transaction) from which the transfer was made. It contains information such as the entry
// ID, the account ID from which the transfer was made, the amount transferred, and the time at which
// the transfer was made.
// @property {Entry} ToEntry - ToEntry is a property of the TransferTxResult struct that represents the
// entry created in the recipient's account as a result of the transfer transaction. An entry is a
// record of a financial transaction that includes information such as the amount transferred, the date
// and time of the transaction, and the accounts involved.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create from entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create to entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// get from account
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
