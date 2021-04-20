package services

import (
	"errors"
	"sort"
	"time"

	"fetchrewards.com/points-api/internal/model"
)

// pointsDB is an abstraction for the database layer dependencies used by this package
type pointsDB interface {
	AddTransaction(userID string, transaction model.Transaction)
	GetAccounts(userID string) []model.Account
	GetAccount(userID string, payer string) (model.Account, bool)
	GetTransactions(userID string) []model.Transaction
}

var notEnoughPointsErr = errors.New("not enough points")

// PointService houses the business logic of the api. It delegates data manipulation tasks
// to a pointsDB interface.
type PointService struct {
	DB pointsDB
}

// NewPointService creates a new PointService with the given pointsDB
func NewPointService(db pointsDB) *PointService {
	return &PointService{
		DB: db,
	}
}

// AddPoints adds the given model.Transaction to the db. If the point value is negative
// then it must not take the payer's account balance lower than 0. If it results in a
// negative account balance, an error will be returned.
func (s *PointService) AddPoints(userID string, transaction model.Transaction) error {
	if transaction.Points > 0 {
		s.DB.AddTransaction(userID, transaction)
		return nil
	} else {
		totalPoints := s.getTotalPointsForPayer(userID, transaction.Payer)

		if totalPoints >= -transaction.Points {
			s.DB.AddTransaction(userID, transaction)
			return nil
		} else {
			return notEnoughPointsErr
		}
	}
}

// SpendPoints consumes points from transactions starting with the oldest transaction going
// forward and returns new transactions as a result of the operation. Returns an error if
// there are not enough points.
func (s *PointService) SpendPoints(userID string, points int) ([]model.Transaction, error) {
	transactions := s.DB.GetTransactions(userID)

	pointsRemaining := points
	newTranMap := make(map[string]*model.Transaction, 0)
	for i := 0; i < len(transactions) && pointsRemaining > 0; i++ {
		tran := transactions[i]

		balanceDiff := -tran.Points
		if tran.Points > pointsRemaining {
			balanceDiff = -pointsRemaining
		}

		pointsRemaining -= tran.Points
		if _, ok := newTranMap[tran.Payer]; ok {
			t := newTranMap[tran.Payer]
			t.Points = t.Points + balanceDiff
		} else {
			newTranMap[tran.Payer] = &model.Transaction{
				Payer:     tran.Payer,
				Points:    balanceDiff,
				Timestamp: time.Now(),
			}
		}
	}

	if pointsRemaining > 0 {
		return []model.Transaction{}, notEnoughPointsErr
	}

	var newTransactions []model.Transaction
	for _, val := range newTranMap {
		newTransactions = append(newTransactions, *val)
		s.DB.AddTransaction(userID, *val)
	}

	sort.Slice(newTransactions, func(i, j int) bool {
		return newTransactions[i].Timestamp.Before(newTransactions[j].Timestamp)
	})

	return newTransactions, nil
}

// GetAccounts returns all payer accounts which includes the associated balances.
func (s *PointService) GetAccounts(userID string) []model.Account {
	return s.DB.GetAccounts(userID)
}

func (s *PointService) getTotalPointsForPayer(userID string, payer string) int {
	pointSum := 0
	transactions := s.DB.GetTransactions(userID)
	for _, tran := range transactions {
		if tran.Payer == payer {
			pointSum += tran.Points
		}
	}
	return pointSum
}
