package services

import (
	"context"
	"errors"
	"sort"
	"time"

	"fetchrewards.com/points-api/internal/model"
)

// pointsDB is an abstraction for the database layer dependencies used by this package
type pointsDB interface {
	AddTransaction(ctx context.Context, transaction model.Transaction) error
	GetAccounts(ctx context.Context) ([]model.Account, error)
	GetAccount(ctx context.Context, payer string) (model.Account, bool, error)
	GetTransactions(ctx context.Context) ([]model.Transaction, error)
}

var NotEnoughPointsErr = errors.New("not enough points")

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
// negative account balance, an error will be returned. An error can also be returned
// if the database layer returns an error.
func (s *PointService) AddPoints(ctx context.Context, transaction model.Transaction) error {
	if transaction.Points > 0 {
		return s.DB.AddTransaction(ctx, transaction)
	} else {
		totalPoints, err := s.getTotalPointsForPayer(ctx, transaction.Payer)
		if err != nil {
			return err
		}

		if totalPoints >= -transaction.Points {
			return s.DB.AddTransaction(ctx, transaction)
		} else {
			return NotEnoughPointsErr
		}
	}
}

// SpendPoints consumes points from transactions starting with the oldest transaction going
// forward and returns new transactions as a result of the operation. Returns an error if
// there are not enough points or if the db layer returns an error.
func (s *PointService) SpendPoints(ctx context.Context, points int) ([]model.Transaction, error) {
	transactions, err := s.DB.GetTransactions(ctx)
	if err != nil {
		return []model.Transaction{}, err
	}

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
		return []model.Transaction{}, NotEnoughPointsErr
	}

	var newTransactions []model.Transaction
	for _, val := range newTranMap {
		newTransactions = append(newTransactions, *val)
		s.DB.AddTransaction(ctx, *val)
	}

	sort.Slice(newTransactions, func(i, j int) bool {
		return newTransactions[i].Timestamp.Before(newTransactions[j].Timestamp)
	})

	return newTransactions, nil
}

// GetAccounts returns all payer accounts which includes the associated balances.
// Returns an error if the db layer returns an error.
func (s *PointService) GetAccounts(ctx context.Context) ([]model.Account, error) {
	return s.DB.GetAccounts(ctx)
}

func (s *PointService) getTotalPointsForPayer(ctx context.Context, payer string) (int, error) {
	pointSum := 0
	transactions, err := s.DB.GetTransactions(ctx)
	if err != nil {
		return 0, err
	}
	for _, tran := range transactions {
		if tran.Payer == payer {
			pointSum += tran.Points
		}
	}
	return pointSum, nil
}
