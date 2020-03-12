package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	app "github.com/erikacarvalho/stone-challenge"
	"github.com/erikacarvalho/stone-challenge/store"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccounts(t *testing.T) {
	t.Run("should return list of all accounts on GET", func(t *testing.T) {
		account1 := app.Account{
			ID:        1,
			Name:      "Benício Clemente Shinoda",
			CPF:       "63000399003",
			Balance:   985845,
			CreatedAt: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		}
		account2 := app.Account{
			ID:        2,
			Name:      "Arlene Araújo Nogueira",
			CPF:       "08312653457",
			Balance:   2578265,
			CreatedAt: time.Date(2020, time.February, 9, 15, 0, 0, 0, time.UTC),
		}
		account3 := app.Account{
			ID:        3,
			Name:      "Bruna Carvalho Lemos",
			CPF:       "21715382609",
			Balance:   27380,
			CreatedAt: time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
		}

		accountStore := store.NewAccountStore(
			app.StartingID(3),
			account1, account2, account3,
		)
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []app.Account
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response. response: %q; error: '%v'", response.Body, err)
		}

		want := []app.Account{account1, account2, account3}

		for i := range want {
			if want[i] != got[i] {
				t.Errorf("got %v; want %v", got[i], want[i])
				t.FailNow()
			}
		}

		app.AssertHTTPStatus(t, response.Code, http.StatusOK)
	})

	t.Run("should return the accounts list as json", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(90))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(109))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodDelete, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		app.AssertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("should create account on POST", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(879))
		server := NewServer(accountStore, nil)
		accountRequest := CreateAccountRequest{
			Name:    "Arlene Araújo Nogueira",
			CPF:     "08312653457",
			Balance: 2578265,
		}

		jsonAcc, err := json.Marshal(accountRequest)
		if err != nil {
			t.Fatalf("could not marshal given account. error: %q", err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonAcc))
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":880}`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusCreated)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)

		//Integration test
		gotMaxID := accountStore.GetMaxID()
		wantMaxID := uint64(880)

		if gotMaxID != wantMaxID {
			t.Errorf("POST on /accounts is not writing on account store. got %d as maxID; want %d", gotMaxID, wantMaxID)
		}
	})

	t.Run("should return error if cpf is invalid", func(t *testing.T) {
		server := NewServer(nil, nil)

		accountRequest := CreateAccountRequest{
			Name:    "Maria das Neves",
			CPF:     "083",
			Balance: 1000,
		}

		jsonAcc, _ := json.Marshal(accountRequest)

		request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonAcc))
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `invalid cpf: it must have 11 numbers`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if name is invalid", func(t *testing.T) {
		server := NewServer(nil, nil)

		accountRequest := CreateAccountRequest{
			Name:    "",
			CPF:     "08389076580",
			Balance: 1000,
		}

		jsonAcc, _ := json.Marshal(accountRequest)

		request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonAcc))
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `invalid name: it cannot be empty`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return bad request when body is nil", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(5087))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodPost, "/accounts", nil)
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})
}

func TestAccountsBalance(t *testing.T) {
	t.Run("should return balance by account ID on GET", func(t *testing.T) {
		account1 := app.Account{
			ID:        550,
			Name:      "Bruna Carvalho Lemos",
			CPF:       "21715382609",
			Balance:   27380,
			CreatedAt: time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
		}

		accountStore := store.NewAccountStore(
			app.StartingID(1),
			account1,
		)
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts/550/balance", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":550,"balance":27380}`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusOK)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should display error message if account ID is not found", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(15))
		server := NewServer(accountStore, nil)

		inexistentID := 97

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%v/balance", inexistentID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`account %v not found`, inexistentID)

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("should display error message if account ID is invalid", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(3007))
		server := NewServer(accountStore, nil)

		invalidID := "letters"

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%v/balance", invalidID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`account ID is invalid. ID given: %v`, invalidID)

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})
}

func TestTransfers(t *testing.T) {
	t.Run("should return empty list of transfers on GET", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(7))
		transferStore := store.NewTransferStore(app.StartingID(30))

		server := NewServer(accountStore, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `[]`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusOK)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return list of all transfers on GET", func(t *testing.T) {
		transfer1 := app.Transfer{
			ID:                   1,
			AccountOriginID:      60,
			AccountDestinationID: 190,
			Amount:               15000,
			CreatedAt:            time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
			Status:               store.ToStatusMsg(store.StatusConfirmed),
		}

		transfer2 := app.Transfer{
			ID:                   2,
			AccountOriginID:      190,
			AccountDestinationID: 97,
			Amount:               60000,
			CreatedAt:            time.Date(2020, time.February, 16, 10, 0, 0, 0, time.UTC),
			Status:               store.ToStatusMsg(store.StatusNotAuthorized),
		}

		transfer3 := app.Transfer{
			ID:                   3,
			AccountOriginID:      97,
			AccountDestinationID: 60,
			Amount:               5000,
			CreatedAt:            time.Date(2020, time.February, 16, 19, 0, 0, 0, time.UTC),
			Status:               store.ToStatusMsg(store.StatusConfirmed),
		}

		transfer4 := app.Transfer{
			ID:                   4,
			AccountOriginID:      60,
			AccountDestinationID: 190,
			Amount:               50000,
			CreatedAt:            time.Date(2020, time.February, 18, 14, 0, 0, 0, time.UTC),
			Status:               store.ToStatusMsg(store.StatusConfirmed),
		}

		transferStore := store.NewTransferStore(
			app.StartingID(4),
			transfer1, transfer2, transfer3, transfer4,
		)

		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []app.Transfer
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response. response: %q; error: '%v'", response.Body, err)
		}

		want := []app.Transfer{transfer1, transfer2, transfer3, transfer4}

		for i := range want {
			if want[i] != got[i] {
				t.Errorf("got %v; want %v", got[i], want[i])
				t.FailNow()
			}
		}

		app.AssertHTTPStatus(t, response.Code, http.StatusOK)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return bad request when body is nil on POST", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(5087))
		transferStore := store.NewTransferStore(app.StartingID(500))
		server := NewServer(accountStore, transferStore)

		request, _ := http.NewRequest(http.MethodPost, "/transfers", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
		app.AssertResponseBody(t, response.Body.String(), "invalid request")
	})

	t.Run("should successfully transfer amount from origin account to destination account on POST", func(t *testing.T) {
		account1 := app.Account{
			ID:        207,
			Name:      "Juliana da Cruz Clemente",
			CPF:       "63000399003",
			Balance:   70000,
			CreatedAt: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		}
		account2 := app.Account{
			ID:        986,
			Name:      "Marlene de Souza Dalponte",
			CPF:       "08312653457",
			Balance:   51000,
			CreatedAt: time.Date(2020, time.February, 9, 15, 0, 0, 0, time.UTC),
		}

		accountStore := store.NewAccountStore(
			app.StartingID(2),
			account1, account2,
		)

		transferStore := store.NewTransferStore(app.StartingID(870))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      207,
			AccountDestinationID: 986,
			Amount:               4000,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		wantTransferID := `{"id":871}`
		gotTransferID := response.Body.String()

		//Assert balance
		wantBalance1 := uint64(66000)
		wantBalance2 := uint64(55000)
		acc1, _ := accountStore.GetAccount(207)
		gotBalance1 := acc1.Balance
		acc2, _ := accountStore.GetAccount(986)
		gotBalance2 := acc2.Balance
		app.AssertUint64(t, gotBalance1, wantBalance1)
		app.AssertUint64(t, gotBalance2, wantBalance2)

		//Assert status
		wantTransferStatus := store.ToStatusMsg(store.StatusConfirmed)
		transfer, _ := transferStore.GetTransfer(871)
		gotTransferStatus := transfer.Status
		app.AssertString(t, gotTransferStatus, wantTransferStatus)

		//Assert http
		app.AssertResponseBody(t, gotTransferID, wantTransferID)
		app.AssertHTTPStatus(t, response.Code, http.StatusCreated)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return error if origin account ID is not found on POST", func(t *testing.T) {
		accountStore := store.NewAccountStore(app.StartingID(90))
		transferStore := store.NewTransferStore(app.StartingID(0))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 900,
			Amount:               100000,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `account 307 not found. error: "there is no account with this ID"`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if destination account ID is not found", func(t *testing.T) {
		account1 := app.Account{
			ID:        307,
			Name:      "Fernanda das Neves",
			CPF:       "78900167850",
			Balance:   205000,
			CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
		}

		accountStore := store.NewAccountStore(
			app.StartingID(1),
			account1,
		)
		transferStore := store.NewTransferStore(app.StartingID(0))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 900,
			Amount:               9700,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `account 900 not found. error: "there is no account with this ID"`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if origin and destination account are the same", func(t *testing.T) {
		account1 := app.Account{

			ID:        307,
			Name:      "Fernanda das Neves",
			CPF:       "78900167850",
			Balance:   205000,
			CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
		}
		accountStore := store.NewAccountStore(
			app.StartingID(1),
			account1,
		)
		transferStore := store.NewTransferStore(app.StartingID(0))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 307,
			Amount:               9700,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `error transferring from account [307] to account [307]: origin and destination account ids are the same`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if amount is invalid", func(t *testing.T) {
		account1 := app.Account{
			ID:        307,
			Name:      "Fernanda das Neves",
			CPF:       "78900167850",
			Balance:   205000,
			CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
		}
		account2 := app.Account{
			ID:        405,
			Name:      "Marcela das Neves",
			CPF:       "78900389020",
			Balance:   780000,
			CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
		}
		accountStore := store.NewAccountStore(
			app.StartingID(2),
			account1, account2,
		)

		transferStore := store.NewTransferStore(app.StartingID(0))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 405,
			Amount:               0,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `error transferring from account [307] to account [405]: the amount entered is invalid`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if origin account balance is smaller than amount to be transferred", func(t *testing.T) {
		account1 := app.Account{
			ID:        307,
			Name:      "Fernanda das Neves",
			CPF:       "78900167850",
			Balance:   205000,
			CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
		}
		account2 := app.Account{
			ID:        405,
			Name:      "Marcela das Neves",
			CPF:       "78900389020",
			Balance:   780000,
			CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
		}
		accountStore := store.NewAccountStore(
			app.StartingID(2),
			account1, account2,
		)

		transferStore := store.NewTransferStore(app.StartingID(0))

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 405,
			Amount:               300000,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `error transferring from account [307] to account [405]: origin account balance is too low to allow this transfer`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if transfer seems to be duplicated", func(t *testing.T) {
		account1 := app.Account{
			ID:        307,
			Name:      "Fernanda das Neves",
			CPF:       "78900167850",
			Balance:   205000,
			CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
		}
		account2 := app.Account{
			ID:        405,
			Name:      "Marcela das Neves",
			CPF:       "78900389020",
			Balance:   780000,
			CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
		}

		accountStore := store.NewAccountStore(
			app.StartingID(2),
			account1, account2,
		)

		transfer1 := app.Transfer{
			ID:                   1,
			AccountOriginID:      307,
			AccountDestinationID: 405,
			Amount:               15000,
			CreatedAt:            time.Now(),
			Status:               store.ToStatusMsg(store.StatusConfirmed),
		}

		transferStore := store.NewTransferStore(app.StartingID(1), transfer1)

		server := NewServer(accountStore, transferStore)

		transferRequest := CreateTransferRequest{
			AccountOriginID:      307,
			AccountDestinationID: 405,
			Amount:               15000,
		}

		jsonTransfer, err := json.Marshal(transferRequest)
		if err != nil {
			t.Fatalf("could not marshal given transfer. error: %q", jsonTransfer)
		}

		request, _ := http.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(jsonTransfer))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `error transferring from account [307] to account [405]: this transfer seems to be duplicated`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		transferStore := store.NewTransferStore(app.StartingID(109))
		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodDelete, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		app.AssertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestTransfersID(t *testing.T) {
	t.Run("should return transfer depending on given ID on GET", func(t *testing.T) {
		transfer1 := app.Transfer{
			ID:                   706,
			AccountOriginID:      45,
			AccountDestinationID: 78,
			Amount:               9500,
			CreatedAt:            time.Date(2020, time.March, 3, 11, 30, 0, 0, time.UTC),
			Status:               store.ToStatusMsg(store.StatusConfirmed),
		}
		transferStore := store.NewTransferStore(app.StartingID(1), transfer1)

		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers/706", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":706,"account_origin_id":45,"account_destination_id":78,"amount":9500,"created_at":"2020-03-03T11:30:00Z","status":"Confirmed"}`

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusOK)
		app.AssertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should display error message if transfer ID is not found", func(t *testing.T) {
		transferStore := store.NewTransferStore(app.StartingID(10))
		server := NewServer(nil, transferStore)

		inexistentID := 101

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfers/%v", inexistentID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`transfer %v not found`, inexistentID)

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("should display error message if transfer ID is invalid", func(t *testing.T) {
		transferStore := store.NewTransferStore(app.StartingID(30))
		server := NewServer(nil, transferStore)

		invalidID := "letters"

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfers/%v", invalidID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`transfer ID is invalid. ID given: %v`, invalidID)

		app.AssertResponseBody(t, got, want)
		app.AssertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		transferStore := store.NewTransferStore(app.StartingID(109))
		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodDelete, "/transfers/100", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		app.AssertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}
