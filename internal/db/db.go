package db

import (
	"log"
	"sort"

	"fetchrewards.com/points-api/internal/model"
)

// InMemoryDB holds all Transactions in memory and provides various accessors for the data
type InMemoryDB struct {
	UserTransactions map[string][]model.Transaction
}

// NewInMemoryDB returns a new InMemoryDB and initializes an empty slice of model.Transactions.
// Because this is an in-memory implementation, state is not maintained between app restarts.
func NewInMemoryDB() *InMemoryDB {
	log.Println("Creating new in-memory database")

	userTransactions := make(map[string][]model.Transaction, 0)
	return &InMemoryDB{
		UserTransactions: userTransactions,
	}
}

// GetTransactions returns all the model.Transaction records in time ascending order for the user
func (db *InMemoryDB) GetTransactions(userID string) []model.Transaction {
	if transactions, ok := db.UserTransactions[userID]; ok {
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Timestamp.Before(transactions[j].Timestamp)
		})
		return transactions
	} else {
		return []model.Transaction{}
	}
}

// AddTransaction adds the given model.Transaction for this user. This function
// makes no assumptions about business logic. For example it does no validation
// that the sum of transactions should not be negative
func (db *InMemoryDB) AddTransaction(userID string, transaction model.Transaction) {
	if transactions, ok := db.UserTransactions[userID]; ok {
		db.UserTransactions[userID] = append(transactions, transaction)
	} else {
		db.UserTransactions[userID] = []model.Transaction{transaction}
	}
}

// GetAccounts returns all model.Accounts, or payers, across all transactions for this user
func (db *InMemoryDB) GetAccounts(userID string) []model.Account {
	accountMap := db.getAccountMap(userID)

	result := make([]model.Account, 0)
	for _, value := range accountMap {
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Payer < result[j].Payer
	})
	return result
}

// GetAccount returns the model.Account associated with the payer for this user
func (db *InMemoryDB) GetAccount(userID, payer string) (model.Account, bool) {
	accountMap := db.getAccountMap(userID)
	account, found := accountMap[payer]
	return account, found
}

func (db *InMemoryDB) getAccountMap(userID string) map[string]model.Account {
	var accountMap = make(map[string]model.Account, 0)
	transactions := db.GetTransactions(userID)
	for _, tran := range transactions {
		if account, ok := accountMap[tran.Payer]; ok {
			account.Points += tran.Points
			accountMap[tran.Payer] = account
		} else {
			accountMap[tran.Payer] = model.Account{
				Payer:  tran.Payer,
				Points: tran.Points,
			}
		}
	}
	return accountMap
}
