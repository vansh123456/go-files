package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vansh123456/simplebank/util"
)

func createRandomAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err) //t is an object error
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, arg.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
func TestCreateAccount(t *testing.T) {
	//create account
	createRandomAccount(t)
}
