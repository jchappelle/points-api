package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"fetchrewards.com/points-api/internal/db"
	"fetchrewards.com/points-api/internal/model"
	"fetchrewards.com/points-api/internal/services"
	"fetchrewards.com/points-api/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestGetPayers(t *testing.T) {
	withEnv(t, func(env serverEnv) {
		userID := "1"
		for _, transaction := range test.Data {
			err := env.service.AddPoints(userID, transaction)
			if err != nil {
				t.Fatal(err)
			}
		}
		resp := env.PerformRequest("GET", fmt.Sprintf("/v1/users/%s/payers", userID), nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status 200")

		accounts := []model.Account{}
		err := json.NewDecoder(resp.Body).Decode(&accounts)
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, accounts, 3)
		assert.Equal(t, accounts[0].Payer, "DANNON")
		assert.Equal(t, accounts[0].Points, 1100)
	})
}

func TestAddTransaction(t *testing.T) {
	withEnv(t, func(env serverEnv) {
		userID := "1"

		resp := env.PerformRequest(
			"POST",
			fmt.Sprintf("/v1/users/%s/points/add", userID),
			model.Transaction{
				Payer:     "DANNON",
				Points:    1000,
				Timestamp: test.ParseTime("2020-11-02T14:00:00Z"),
			})
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Should return status 204")

		transactions := env.db.GetTransactions(userID)
		assert.Len(t, transactions, 1)
		assert.Equal(t, "DANNON", transactions[0].Payer)
	})
}

func TestSpendPoints(t *testing.T) {
	withEnv(t, func(env serverEnv) {
		userID := "1"
		for _, transaction := range test.Data {
			err := env.service.AddPoints(userID, transaction)
			assert.NoError(t, err)
		}

		resp := env.PerformRequest(
			"POST",
			fmt.Sprintf("/v1/users/%s/points/spend", userID),
			spendPointsRequest{
				Points: 5000,
			})
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status 200")

		// Check response for new transactions
		transactions := []model.Transaction{}
		err := json.NewDecoder(resp.Body).Decode(&transactions)
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, transactions, 3)
		assert.Equal(t, "DANNON", transactions[0].Payer)
		assert.Equal(t, -100, transactions[0].Points)
		assert.Equal(t, "UNILEVER", transactions[1].Payer)
		assert.Equal(t, -200, transactions[1].Points)
		assert.Equal(t, "MILLER COORS", transactions[2].Payer)
		assert.Equal(t, -4700, transactions[2].Points)

		// Check DB has the new transactions
		transactions = env.db.GetTransactions(userID)
		assert.Len(t, transactions, 8)
		assert.Equal(t, "DANNON", transactions[5].Payer)
		assert.Equal(t, -100, transactions[5].Points)
		assert.Equal(t, "UNILEVER", transactions[6].Payer)
		assert.Equal(t, -200, transactions[6].Points)
		assert.Equal(t, "MILLER COORS", transactions[7].Payer)
		assert.Equal(t, -4700, transactions[7].Points)
	})
}

// serverEnv is a struct used to house test dependencies
type serverEnv struct {
	db      *db.InMemoryDB
	service *services.PointService
	server  *Server
	t *testing.T
}

// withEnv sets up common test dependencies and helper methods and makes them available via a serverEnv
// to the given function
func withEnv(t *testing.T, f func(env serverEnv)) {
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)
	server := NewServer(service)

	env := serverEnv{
		db:      db,
		service: service,
		server:  server,
		t: t,
	}

	f(env)
}

func (e *serverEnv) PerformRequest(method, url string, payloadObj interface{}) *http.Response {
	body, err := json.Marshal(payloadObj)
	if err != nil {
		e.t.Fatal(err)
	}

	r, err := http.NewRequest(method, url, bytes.NewBufferString(string(body)))
	if err != nil {
		e.t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := e.server.setupHandlers()
	handler.ServeHTTP(w, r)

	return w.Result()
}