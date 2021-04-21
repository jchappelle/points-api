package services_test

import (
	"fetchrewards.com/points-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/test"
)

func TestSpendPoints(t *testing.T) {
	// Holder for inputs and expectations
	type spendTest struct {
		input []model.Transaction
		points int
		expected []model.Transaction
		errExpected bool
	}

	// Create test cases
	tests := map[string]spendTest{
		"Assignment example": {
			input: test.Data,
			points: 5000,
			expected: []model.Transaction{
				{Payer: "DANNON", Points: -100},
				{Payer: "UNILEVER", Points: -200},
				{Payer: "MILLER COORS", Points: -4700},
			},
		},
		"Single transaction": {
			input: []model.Transaction{
				{Payer: "DANNON", Points: 1000},
			},
			points: 999,
			expected: []model.Transaction{
				{Payer: "DANNON", Points: -999},
			},
		},
		"Negative points requested returns error": {
			input: test.Data,
			points: -1,
			expected: []model.Transaction{},
			errExpected: true,
		},
		"Insufficient points returns error": {
			input: []model.Transaction{
				{Payer: "DANNON", Points: 1000},
			},
			points: 1001,
			expected: []model.Transaction{},
			errExpected: true,
		},
	}


	// Execute tests
	for name, tc := range tests {
		t.Run(name, func(t *testing.T){
			userID := "1"
			db := db.NewInMemoryDB()
			service := services.NewPointService(db)

			for _, transaction := range tc.input {
				err := service.AddPoints(userID, transaction)
				assert.NoError(t, err)
			}

			actual, err := service.SpendPoints(userID, tc.points)
			if tc.errExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, len(actual), len(tc.expected))

			for i := range actual {
				assert.Equal(t, actual[i].Payer, tc.expected[i].Payer)
				assert.Equal(t, actual[i].Points, tc.expected[i].Points)
			}
		})
	}
}

func TestAddPoints(t *testing.T){
	userID := "1"
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)

	// Load test data
	for _, tran := range test.Data {
		err := service.AddPoints(userID, tran)
		assert.NoError(t, err)
	}

	// No error when negative points and balance is sufficient
	tran := model.Transaction{
		Payer: "MILLER COORS",
		Points: -4000,
	}
	err := service.AddPoints(userID, tran)
	assert.NoError(t, err)

	// Error when negative points and balance is insufficient
	tran = model.Transaction{
		Payer: "DANNON",
		Points: -4000,
	}
	err = service.AddPoints(userID, tran)
	assert.Error(t, err)
}
