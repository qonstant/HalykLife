package db

import (
	"Simple-Bank/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	balance1 := account1.Balance
	balance2 := account2.Balance

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, account1.ID)
	require.NotZero(t, account2.ID)
	require.Positive(t, account1.Balance)
	require.NotZero(t, account2.Balance)

	require.Equal(t, account2.Balance-balance2, balance1-account1.Balance)

	return transfer
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer2.ID, transfer1.ID)
	require.Equal(t, transfer2.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer2.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer2.Amount, transfer1.Amount)
	require.Equal(t, transfer2.CreatedAt, transfer1.CreatedAt)

	require.NotZero(t, transfer2.ID)
	require.NotZero(t, transfer2.CreatedAt)

	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestCreateTransfer(t *testing.T) {
	transfer := createRandomTransfer(t)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

	require.Empty(t, transfer2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 15; i++ {
		createRandomTransfer(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 10,
	}

	transfers, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
