package store

import (
	app "github.com/erikacarvalho/stone-challenge"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	t.Run("should return autogenerated ID", func(t *testing.T) {
		store := NewAccountStore(app.StartingID(0))

		want := uint64(1)
		got, _ := store.CreateAccount("", "", 0)

		app.AssertUint64(t, got, want)
	})

	t.Run("should return last possible ID", func(t *testing.T) {
		store := NewAccountStore(app.StartingID(90))

		want := uint64(91)
		got, _ := store.CreateAccount("", "", 0)

		app.AssertUint64(t, got, want)
	})
}

func TestGetBalance(t *testing.T) {
	t.Run("should return account balance", func(t *testing.T) {
		store := NewAccountStore(app.StartingID(7690))

		newAccountID, _ := store.CreateAccount("", "", 10)
		accountBalance, _ := store.GetBalance(newAccountID)

		want := uint64(10)
		got := accountBalance

		app.AssertUint64(t, got, want)
	})

	t.Run("should return ErrAccountNotFound when there is no account with given ID", func(t *testing.T) {
		store := NewAccountStore(app.StartingID(7690))

		want := ErrAccountNotFound
		_, got := store.GetBalance(789)

		app.AssertError(t, got, want)
	})
}

func TestListAllAccounts(t *testing.T) {
	accounts := map[uint64]app.Account{
		1: {
			ID:        1,
			Name:      "Talita Barreto Coelho",
			CPF:       "96097705840",
			Balance:   7590000,
			CreatedAt: time.Now(),
		},
		2: {
			ID:        2,
			Name:      "Maurício Ximenes Brito",
			CPF:       "37320891697",
			Balance:   290000,
			CreatedAt: time.Now(),
		},
		3: {
			ID:        3,
			Name:      "Carolina Monteiro Hamada",
			CPF:       "54009199520",
			Balance:   15000,
			CreatedAt: time.Now(),
		},
	}

	store := &AccountStore{
		maxID:       app.StartingID(len(accounts)),
		dataStorage: accounts,
	}

	t.Run("should return slice with all created accounts", func(t *testing.T) {
		accountsList, _ := store.ListAllAccounts()

		for i, account := range accountsList {
			want := accounts[uint64(i+1)]
			got := account
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		}
	})

	t.Run("should return an ordered list of accounts", func(t *testing.T) {
		accountsList, _ := store.ListAllAccounts()

		for i, account := range accountsList {
			want := store.dataStorage[uint64(i+1)].ID
			got := account.ID
			if got != want {
				t.Errorf("couldn't return ordered account list. got %d, want %d", got, want)
			}
		}
	})

	t.Run("should return error if accountStore is empty", func(t *testing.T) {
		emptyStore := NewAccountStore(app.StartingID(0))

		_, got := emptyStore.ListAllAccounts()
		want := ErrNoRecords

		app.AssertError(t, got, want)
	})
}