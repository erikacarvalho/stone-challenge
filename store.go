package app

import (
	"fmt"
	"sync/atomic"
	"time"
)

type AccountStore struct {
	maxID       uint64
	dataStorage map[uint64]Account //the map key is the account identifier
}

func (a *AccountStore) ListAll() []Account {
	accs := make([]Account, len(a.dataStorage))
	for _, v := range a.dataStorage {
		accs = append(accs, v)
	}
	return accs
}

func (a *AccountStore) GetBalance(ID uint64) (balance uint64, err error) {
	acc, ok := a.dataStorage[ID]
	if !ok {
		err = fmt.Errorf("account %q not found", ID)
		return
	}
	return acc.Balance, nil
}

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