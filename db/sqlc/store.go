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

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	txName := ctx.Value(txKey)
	fmt.Println(txName, "create transfer")

	err := store.execTX(ctx, func(q *Queries) error {
		var err error
		fmt.Println("create transfer entry")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println("create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		fmt.Println("create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//update account ba;lances:

		//here we are using this if check to avoid deadlock
		// here we are doing ki firstly the smaller ID account will be made to perform the transaction first
		if arg.FromAccountID < arg.ToAccountID {
			fmt.Println("add account balance:update and transfer 1")
			//first wali be executed first then to_account next
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			//if the ID is reversed we do the updates in the reversed order :)
			fmt.Println("add account balance:update and transfer 2")
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAcountBalance(ctx, AddAcountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAcountBalance(ctx, AddAcountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
