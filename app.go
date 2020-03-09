package app

import "time"

type Account struct {
	ID uint64 `json:"id"`
	Name string `json:"name"`
	CPF string `json:"cpf"`
	Balance uint64 `json:"balance"` // Account balance in cents
	CreatedAt time.Time `json:"created_at"`
}

type AccountService interface {
	ListAll() ([]Account, error)
	GetBalance(ID uint64) (uint64, error)
	Create(name, CPF string, balance uint64) uint64
}