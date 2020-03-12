package store

import (
	"errors"
	app "github.com/erikacarvalho/stone-challenge"
	"sort"
	"sync/atomic"
	"time"
)

var (
	ErrNoRecords       = errors.New("there are no accounts to be listed")
	ErrAccountNotFound = errors.New("there is no account with this ID")
)

type AccountStore struct {
	maxID       *uint64
	dataStorage map[uint64]app.Account // The map key is the account identifier
}

func NewAccountStore(startingID *uint64, accounts ...app.Account) *AccountStore {
	storage := make(map[uint64]app.Account)
	for _, account := range accounts {
		storage[account.ID] = account
	}
	ns := &AccountStore{
		maxID:       startingID,
		dataStorage: storage,
	}
	return ns
}

func (a *AccountStore) GetMaxID() uint64 {
	return *a.maxID
}

// CreateAccount is a method that creates an account and returns its ID.
func (a *AccountStore) CreateAccount(name, CPF string, balance uint64) (ID uint64, err error) {
	newID := atomic.AddUint64(a.maxID, 1)
	a.dataStorage[newID] = app.Account{
		ID:        newID,
		Name:      name,
		CPF:       CPF,
		Balance:   balance,
		CreatedAt: time.Now(),
	}
	return newID, nil
}

// ListAllAccounts returns all accounts from the account store sorted.
func (a *AccountStore) ListAllAccounts() ([]app.Account, error) {
	var accs []app.Account
	for _, v := range a.dataStorage {
		accs = append(accs, v)
	}

	if len(accs) == 0 {
		return nil, ErrNoRecords
	}

	sort.Slice(accs, func(i, j int) bool {
		return accs[i].ID < accs[j].ID
	})

	return accs, nil
}

// GetBalance returns balance for account with given ID
// and an error if there is no such account.
func (a *AccountStore) GetBalance(ID uint64) (balance uint64, err error) {
	acc, ok := a.dataStorage[ID]
	if !ok {
		return 0, ErrAccountNotFound
	}
	return acc.Balance, nil
}

func (a *AccountStore) GetAccount(ID uint64) (app.Account, error) {
	acc, ok := a.dataStorage[ID]
	if !ok {
		return app.Account{}, ErrAccountNotFound
	}
	return acc, nil
}

func (a *AccountStore) SetAccount(account app.Account) {
	a.dataStorage[account.ID] = account
}
