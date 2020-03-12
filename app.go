package app

import "time"

type Account struct {
	ID        uint64    `json:"id"` // This field is read-only
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   uint64    `json:"balance"` // Account balance in cents
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID                   uint64    `json:"id"` // This field is read-only
	AccountOriginID      uint64    `json:"account_origin_id"`
	AccountDestinationID uint64    `json:"account_destination_id"`
	Amount               uint64    `json:"amount"` // Transfer amount in cents
	CreatedAt            time.Time `json:"created_at"`
	Status               string    `json:"status"`
}
