package services

import (
	"context"
	"errors"
	"fetchrewards.com/points-api/internal/model"
	"sort"
	"time"
)

type pointsDB interface {
	AddTransaction(ctx context.Context, transaction model.Transaction) error
	GetAccounts(ctx context.Context) ([]model.Account, error)
	GetAccount(ctx context.Context, payer string) (model.Account, bool, error)
	GetTransactions(ctx context.Context) ([]model.Transaction, error)
}

var NotEnoughPointsErr = errors.New("not enough points")

type PointService struct {
	DB pointsDB
}
func NewPointService(db pointsDB) *PointService {
	return &PointService{
		DB: db,
	}
}

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
				Payer: tran.Payer,
				Points: balanceDiff,
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