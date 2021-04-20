package services_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/test"
)

func TestSpendPoints(t *testing.T) {
	ctx := context.Background()
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)

	for _, transaction := range test.Data {
		err := service.AddPoints(ctx, transaction)
		assert.NoError(t, err)
	}

	transactions, err := service.SpendPoints(ctx, 5000)
	assert.NoError(t, err)
	assert.Len(t, transactions, 3)

	tran := transactions[0]
	assert.Equal(t, "DANNON", tran.Payer)
	assert.Equal(t, -100, tran.Points)

	tran = transactions[1]
	assert.Equal(t, "UNILEVER", tran.Payer)
	assert.Equal(t, -200, tran.Points)

	tran = transactions[2]
	assert.Equal(t, "MILLER COORS", tran.Payer)
	assert.Equal(t, -4700, tran.Points)
}
