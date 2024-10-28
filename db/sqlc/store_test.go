package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	Account1 := createRandomAccount(t)
	Account2 := createRandomAccount(t)

	//runs a concurrent transfer transactions

	n := 5
	amount := int64(10)

	//make use of channels which  is a goroutine-safe data structure,connects goroutines and serves as  a buffer for sending and receiving data(error)
	errs := make(chan error)               //starts an err
	results := make(chan TransferTxResult) //accepts channel

	for i := 0; i < n; i++ {
		//go routine
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: Account1.ID,
				ToAccountID:   Account2.ID,
				Amount:        amount,
			})
			//here we are sending  the result to the channel,err is being sent to errs channnel
			errs <- err
			results <- result
		}()
	}
	//check the channel results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check if transfer is working correctly
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, Account1.ID, transfer.FromAccountID)
		require.Equal(t, Account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, Account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, Account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO:CHECK ACC BALANCE
	}
}
