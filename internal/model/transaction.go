package model

import "time"

// Transaction represents an event that changes the state of the overall point balance for
// specific payers
type Transaction struct {
	Payer     string    `json:"payer"`
	Points    int       `json:"points"`
	Timestamp time.Time `json:"timestamp"`
}
