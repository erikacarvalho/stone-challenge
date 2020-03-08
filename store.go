package app

import (
	"fmt"
	"sort"
	"sync/atomic"
	"time"
)

type AccountStore struct {
	maxID       uint64
	dataStorage map[uint64]Account // The map key is the account identifier
}

// ListAll returns all accounts from the store sorted.
func (a *AccountStore) ListAll() []Account {
	var accs []Account
	for _, v := range a.dataStorage {
		accs = append(accs, v)
	}

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].ID < accs[j].ID
	})

	return accs
}

// GetBalance returns balance for account with given ID
// and an error if there is no such account.
func (a *AccountStore) GetBalance(ID uint64) (balance uint64, err error) {
	acc, ok := a.dataStorage[ID]
	if !ok {
		err = fmt.Errorf("account %q not found", ID)
		return
	}
	return acc.Balance, nil
}

// Create is a method that creates an account and returns its ID.
func (a *AccountStore) Create(name, CPF string, balance uint64) uint64 {
	newID := atomic.AddUint64(&a.maxID, 1)
	a.dataStorage[newID] = Account{
		ID:        newID,
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: time.Now(),
	}
	return newID
}