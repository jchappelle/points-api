package model

import "time"

type Transaction struct {
	Payer           string    `json:"payer"`
	Points          int       `json:"points"`
	Timestamp       time.Time `json:"timestamp"`
}
