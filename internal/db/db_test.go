package db_test

import (
	"testing"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestAddTransaction(t *testing.T) {
	userID := "1"

	database := db.NewInMemoryDB()

	database.AddTransaction(userID, test.Data[0])
	database.AddTransaction(userID, test.Data[1])

	transactions := database.GetTransactions(userID)
	assert.Len(t, transactions, 2)

	assert.Equal(t, "UNILEVER", transactions[0].Payer)
	assert.Equal(t, 200, transactions[0].Points)
	assert.Equal(t, test.ParseTime("2020-10-31T11:00:00Z"), transactions[0].Timestamp)

	assert.Equal(t, "DANNON", transactions[1].Payer)
	assert.Equal(t, 1000, transactions[1].Points)
	assert.Equal(t, test.ParseTime("2020-11-02T14:00:00Z"), transactions[1].Timestamp)
}

func TestGetAccounts(t *testing.T) {

	t.Run("returns empty list when no transactions", func(t *testing.T) {
		database := db.NewInMemoryDB()
		accounts := database.GetAccounts("1")
		assert.NotNil(t, accounts)
		assert.Empty(t, accounts)
	})

	t.Run("returns accounts for the correct user", func(t *testing.T) {
		database := db.NewInMemoryDB()

		database.AddTransaction("1", test.Data[0])
		database.AddTransaction("1", test.Data[1])
		database.AddTransaction("2", test.Data[2])

		accounts := database.GetAccounts("1")
		assert.Len(t, accounts, 2)

		accounts = database.GetAccounts("2")
		assert.Len(t, accounts, 1)
	})

	t.Run("returns accounts with the correct point total", func(t *testing.T) {
		database := db.NewInMemoryDB()

		userID := "1"
		for _, tran := range test.Data {
			database.AddTransaction(userID, tran)
		}

		accounts := database.GetAccounts(userID)
		assert.Len(t, accounts, 3)
		assert.Equal(t, "DANNON", accounts[0].Payer)
		assert.Equal(t, 1100, accounts[0].Points)
		assert.Equal(t, "MILLER COORS", accounts[1].Payer)
		assert.Equal(t, 10000, accounts[1].Points)
		assert.Equal(t, "UNILEVER", accounts[2].Payer)
		assert.Equal(t, 200, accounts[2].Points)
	})

}
