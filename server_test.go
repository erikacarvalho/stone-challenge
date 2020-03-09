package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccounts(t *testing.T) {
	t.Run("should return 200 on /accounts", func(t *testing.T) {
		store := &AccountStore{
			maxID:       19040,
			dataStorage: make(map[uint64]Account),
		}
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("should return list of all account on get", func(t *testing.T) {
		acc1 := Account{
			ID:        1,
			Name:      "Benício Clemente Shinoda",
			CPF:       "63000399003",
			Balance:   985845,
			CreatedAt: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		}
		acc2 := Account{
			ID:        2,
			Name:      "Arlene Araújo Nogueira",
			CPF:       "08312653457",
			Balance:   2578265,
			CreatedAt: time.Date(2020, time.February, 9, 15, 0, 0, 0, time.UTC),
		}
		acc3 := Account{
			ID:        3,
			Name:      "Bruna Carvalho Lemos",
			CPF:       "21715382609",
			Balance:   27380,
			CreatedAt: time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
		}
		accs := map[uint64]Account{
			1: acc1,
			2: acc2,
			3: acc3,
		}

		store := &AccountStore{
			maxID:       uint64(len(accs)),
			dataStorage: accs,
		}
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Account
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response. response: %q; error: '%v'", response.Body, err)
		}

		want := []Account{acc1, acc2, acc3}

		for i := range want {
			if want[i] != got[i] {
				t.Errorf("got %v; want %v", got[i], want[i])
				t.FailNow()
			}
		}
	})

	t.Run("should return the accounts list as JSON", func(t *testing.T) {
		store := &AccountStore{
			maxID:       90,
			dataStorage: make(map[uint64]Account),
		}
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Result().Header.Get("content-type")
		want := JsonContentType

		if got != want {
			t.Errorf("response did not have content-type of %v; got %v", want, got)
		}
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("http status is incorrect. got %d; want %d", got, want)
	}
}