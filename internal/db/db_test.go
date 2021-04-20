package db_test

import (
	"testing"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestAddTransaction(t *testing.T) {
	userID := "1"

	db := db.NewInMemoryDB()

	db.AddTransaction(userID, test.Data[0])
	db.AddTransaction(userID, test.Data[1])

	transactions := db.GetTransactions(userID)
	assert.Len(t, transactions, 2)

	assert.Equal(t, "UNILEVER", transactions[0].Payer)
	assert.Equal(t, 200, transactions[0].Points)
	assert.Equal(t, test.ParseTime("2020-10-31T11:00:00Z"), transactions[0].Timestamp)

	assert.Equal(t, "DANNON", transactions[1].Payer)
	assert.Equal(t, 1000, transactions[1].Points)
	assert.Equal(t, test.ParseTime("2020-11-02T14:00:00Z"), transactions[1].Timestamp)
}
