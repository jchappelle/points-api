package db

import (
	"context"
	"log"
	"sort"

	"fetchrewards.com/points-api/internal/model"
)

// InMemoryDB holds all Transactions in memory and provides various accessors for the data
type InMemoryDB struct {
	Transactions []model.Transaction
}

func NewInMemoryDB() *InMemoryDB {
	log.Println("Creating new in-memory database")
	return &InMemoryDB{
		Transactions: []model.Transaction{},
	}
}

func (db *InMemoryDB) GetTransactions(ctx context.Context) ([]model.Transaction, error) {
	result := db.Transactions
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result, nil
}

func (db *InMemoryDB) AddTransaction(ctx context.Context, transaction model.Transaction) error {
	db.Transactions = append(db.Transactions, transaction)
	return nil
}

func (db *InMemoryDB) GetAccounts(ctx context.Context) ([]model.Account, error) {
	accountMap, err := db.getAccountMap(ctx)
	if err != nil {
		return []model.Account{}, err
	}

	var result []model.Account
	for _, value := range accountMap {
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Payer < result[j].Payer
	})
	return result, nil
}

func (db *InMemoryDB) GetAccount(ctx context.Context, payer string) (model.Account, bool, error) {
	accountMap, err := db.getAccountMap(ctx)
	if err != nil {
		return model.Account{}, false, err
	}
	account, found := accountMap[payer]
	return account, found, nil
}

func (db *InMemoryDB) getAccountMap(ctx context.Context) (map[string]model.Account, error) {
	var accountMap = make(map[string]model.Account, 0)
	for _, tran := range db.Transactions {
		if account, ok := accountMap[tran.Payer]; ok {
			account.Points += tran.Points
			accountMap[tran.Payer] = account
		} else {
			accountMap[tran.Payer] = model.Account{
				Payer: tran.Payer,
				Points: tran.Points,
			}
		}
	}
	return accountMap, nil
}
