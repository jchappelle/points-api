package model

// Account represents a payer and holds the payer's name and balance
type Account struct {
	Payer  string `json:"payer"`
	Points int    `json:"points"`
}
