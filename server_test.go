package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccounts(t *testing.T) {
	t.Run("should return status ok on GET", func(t *testing.T) {
		store := newStore(19040)
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("should return list of all accounts on GET", func(t *testing.T) {
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
			maxID:       toPointer(len(accs)),
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

	t.Run("should return the accounts list as json", func(t *testing.T) {
		store := newStore(90)
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

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		store := newStore(109)
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodDelete, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("should create account on POST with correct data", func(t *testing.T) {
		store := newStore(879)
		server := NewAccountServer(store)
		newAcc := CreateAccountRequest{
			Name:      "Arlene Araújo Nogueira",
			CPF:       "08312653457",
			Balance:   2578265,
		}

		jsonAcc, err := json.Marshal(newAcc)
		if err != nil {
			t.Fatalf("could not marshal given account. error: %q", err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonAcc))
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":880}`

		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusCreated)

		//Integration test
		gotFromStore := *store.maxID
		wantFromStore := uint64(880)

		if gotFromStore != wantFromStore {
			t.Errorf("POST on /accounts is not writing on store. got %d as maxID; want %d", gotFromStore, wantFromStore)
		}
	})

	t.Run("should return bad request when body is nil", func(t *testing.T) {
		store := newStore(5087)
		server := NewAccountServer(store)

		request, _ := http.NewRequest(http.MethodPost, "/accounts", nil)
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body has a problem. got %q; want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("http status is incorrect. got %d; want %d", got, want)
	}
}

func toPointer(i int) *uint64 {
	var ptr = uint64(i)
	return &ptr
}

func newStore(i int) *AccountStore {
	ns := &AccountStore{
		maxID:       toPointer(i),
		dataStorage: make(map[uint64]Account),
	}
	return ns
}