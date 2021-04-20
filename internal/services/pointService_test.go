package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/test"
)

func TestSpendPoints(t *testing.T) {
	userID := "1"
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)

	for _, transaction := range test.Data {
		err := service.AddPoints(userID, transaction)
		assert.NoError(t, err)
	}

	transactions, err := service.SpendPoints(userID, 5000)
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
