package db_test

import (
	"context"
	"testing"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestAddTransaction(t *testing.T) {
	ctx := context.Background()
	db := db.NewInMemoryDB()

	err := db.AddTransaction(ctx, test.Data[0])
	assert.NoError(t, err)

	err = db.AddTransaction(ctx, test.Data[1])
	assert.NoError(t, err)

	assert.Len(t, db.Transactions, 2)
	assert.Equal(t, "DANNON", db.Transactions[0].Payer)
	assert.Equal(t, 1000, db.Transactions[0].Points)
	assert.Equal(t, test.ParseTime("2020-11-02T14:00:00Z"), db.Transactions[0].Timestamp)
}
