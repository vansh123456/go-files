package db

import (
	"context"
	"database/sql"
	"fmt"
)

// al functions to execute db queries and transactions
type Store struct {
	//extends struct func by composition
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// this executes a function within the db transaction
// store object hai Store pointer ka
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //begin txn,nil is the isolation level
	if err != nil {
		return err
	}

	q := New(tx) //new db.go wale main there is function
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v,rb:err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//this above was the simple logic to enable a new DB transaction,it is not exported as we dont want the app to call this directly

// now we implement how to transfer the money from one place to another
// this json one shows how FromAccountID will show up in the parsed json
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"` //this Transfer here is from the model.go Transfer
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// func (store *Store) TransferTx(ctx context.Context,arg TransferTxParams) (TransferTxResult,error) {
// 	var  result TransferTxResult

// 	err := store.execTX(ctx,func(q *Queries) error {
// 		var err error

// 		result.Transfer,err = q.//IMPLEMENT LIST OF OTHER CREATE ACCOUNT AND TRANSFER FIRST(ctx,CreateTransferParam {
// 			FromAccountID: arg.FromAccountID,
// 			ToAccountID: arg.ToAccountID,
// 			Amount: arg.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		result.FromEntry,err = q.CreateEntry(ctx,CreateEntryParams {
// 			AccountID: arg.FromAccountID,
// 			Amount: -arg.Amount,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return result,err
// }
