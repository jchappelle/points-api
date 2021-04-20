package web

import (
	"bytes"
	"context"
	"encoding/json"
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
	withEnv(func(env serverEnv){
		for _, transaction := range test.Data {
			env.service.AddPoints(context.Background(), transaction)
		}

		r, err := http.NewRequest("GET", "/v1/payers", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := env.server.setupHandlers()
		handler.ServeHTTP(w, r)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status 200")

		accounts := []model.Account{}
		json.NewDecoder(resp.Body).Decode(&accounts)

		assert.Len(t, accounts, 3)
		assert.Equal(t, accounts[0].Payer, "DANNON")
		assert.Equal(t, accounts[0].Points, 1100)
	})
}

func TestAddTransaction(t *testing.T) {
	withEnv(func(env serverEnv){
		transaction := model.Transaction{
			Payer:     "DANNON",
			Points:    1000,
			Timestamp: test.ParseTime("2020-11-02T14:00:00Z"),
		}
		requestBody, err := json.Marshal(transaction)
		assert.NoError(t, err)

		r, err := http.NewRequest("POST", "/v1/points/add", bytes.NewBufferString(string(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		handler := env.server.setupHandlers()
		handler.ServeHTTP(w, r)

		resp := w.Result()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Should return status 204")
		assert.Len(t, env.db.Transactions, 1)
		assert.Equal(t, "DANNON", env.db.Transactions[0].Payer)
	})
}

func TestSpendPoints(t *testing.T) {
	withEnv(func(env serverEnv){
		for _, transaction := range test.Data {
			env.service.AddPoints(context.Background(), transaction)
		}

		spendPointsRequest := spendPointsRequest{
			Points: 5000,
		}
		requestBody, err := json.Marshal(spendPointsRequest)
		assert.NoError(t, err)

		r, err := http.NewRequest("POST", "/v1/points/spend", bytes.NewBufferString(string(requestBody)))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		handler := env.server.setupHandlers()
		handler.ServeHTTP(w, r)

		resp := w.Result()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return status 200")
		assert.Len(t, env.db.Transactions, 8)
		assert.Equal(t, "DANNON", env.db.Transactions[0].Payer)
		assert.Equal(t, 300, env.db.Transactions[0].Points)
		assert.Equal(t, "MILLER COORS", env.db.Transactions[7].Payer)
		assert.Equal(t, -4700, env.db.Transactions[7].Points)
	})
}

// serverEnv is a struct used to house test dependencies
type serverEnv struct {
	db *db.InMemoryDB
	service *services.PointService
	server *Server
}

// withEnv sets up common test dependencies and makes them available via a serverEnv to the given function
func withEnv(f func (env serverEnv)) {
	db := db.NewInMemoryDB()
	service := services.NewPointService(db)
	server := NewServer(service)

	env := serverEnv{
		db: db,
		service: service,
		server: server,
	}

	f(env)
}
