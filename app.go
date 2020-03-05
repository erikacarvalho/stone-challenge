package app

import "time"

type Account struct {
	ID uint64
	Name string
	CPF string
	Balance uint64
	CreatedAt time.Time `json:"created_at"`
}

type AccountService interface {
	ListAll() []Account
	GetBalance(ID uint64) (uint64, error)
	Create(name, CPF string, balance uint64) uint64
}