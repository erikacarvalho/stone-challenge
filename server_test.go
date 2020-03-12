package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccounts(t *testing.T) {
	t.Run("should return list of all accounts on GET", func(t *testing.T) {
		account1 := Account{
			ID:        1,
			Name:      "Benício Clemente Shinoda",
			CPF:       "63000399003",
			Balance:   985845,
			CreatedAt: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		}
		account2 := Account{
			ID:        2,
			Name:      "Arlene Araújo Nogueira",
			CPF:       "08312653457",
			Balance:   2578265,
			CreatedAt: time.Date(2020, time.February, 9, 15, 0, 0, 0, time.UTC),
		}
		account3 := Account{
			ID:        3,
			Name:      "Bruna Carvalho Lemos",
			CPF:       "21715382609",
			Balance:   27380,
			CreatedAt: time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
		}
		accounts := map[uint64]Account{
			1: account1,
			2: account2,
			3: account3,
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Account
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response. response: %q; error: '%v'", response.Body, err)
		}

		want := []Account{account1, account2, account3}

		for i := range want {
			if want[i] != got[i] {
				t.Errorf("got %v; want %v", got[i], want[i])
				t.FailNow()
			}
		}

		assertHTTPStatus(t, response.Code, http.StatusOK)
	})

	t.Run("should return the accounts list as json", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(90))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(109))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodDelete, "/accounts", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})

	t.Run("should create account on POST", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(879))
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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusCreated)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)

		//Integration test
		gotMaxID := *accountStore.maxID
		wantMaxID := uint64(880)

		if gotMaxID != wantMaxID {
			t.Errorf("POST on /accounts is not writing on account store. got %d as maxID; want %d", gotMaxID, wantMaxID)
		}
	})

	t.Run("should return bad request when body is nil", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(5087))
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodPost, "/accounts", nil)
		request.Header.Set("content-type", JsonContentType)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})
}

func TestAccountsBalance(t *testing.T) {
	t.Run("should return balance by account ID on GET", func(t *testing.T) {
		account1 := Account{
			ID:        550,
			Name:      "Bruna Carvalho Lemos",
			CPF:       "21715382609",
			Balance:   27380,
			CreatedAt: time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
		}
		accounts := map[uint64]Account{
			550: account1,
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}
		server := NewServer(accountStore, nil)

		request, _ := http.NewRequest(http.MethodGet, "/accounts/550/balance", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":550,"balance":27380}`

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusOK)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should display error message if account ID is not found", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(15))
		server := NewServer(accountStore, nil)

		inexistentID := 97

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%v/balance", inexistentID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`account %v not found`, inexistentID)

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("should display error message if account ID is invalid", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(3007))
		server := NewServer(accountStore, nil)

		invalidID := "letters"

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%v/balance", invalidID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`account ID is invalid. ID given: %v`, invalidID)

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})
}

func TestTransfers(t *testing.T) {
	t.Run("should return empty list of transfers on GET", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(7))
		transferStore := NewTransferStore(startingID(30))

		server := NewServer(accountStore, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `[]`

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusOK)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return list of all transfers on GET", func(t *testing.T) {
		transfer1 := Transfer{
			ID:                   1,
			AccountOriginID:      60,
			AccountDestinationID: 190,
			Amount:               15000,
			CreatedAt:            time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
			Status:               toStatusMsg(StatusConfirmed),
		}

		transfer2 := Transfer{
			ID:                   2,
			AccountOriginID:      190,
			AccountDestinationID: 97,
			Amount:               60000,
			CreatedAt:            time.Date(2020, time.February, 16, 10, 0, 0, 0, time.UTC),
			Status:               toStatusMsg(StatusNotAuthorized),
		}

		transfer3 := Transfer{
			ID:                   3,
			AccountOriginID:      97,
			AccountDestinationID: 60,
			Amount:               5000,
			CreatedAt:            time.Date(2020, time.February, 16, 19, 0, 0, 0, time.UTC),
			Status:               toStatusMsg(StatusConfirmed),
		}

		transfer4 := Transfer{
			ID:                   4,
			AccountOriginID:      60,
			AccountDestinationID: 190,
			Amount:               50000,
			CreatedAt:            time.Date(2020, time.February, 18, 14, 0, 0, 0, time.UTC),
			Status:               toStatusMsg(StatusConfirmed),
		}

		transfers := map[uint64]Transfer{
			1: transfer1,
			2: transfer2,
			3: transfer3,
			4: transfer4,
		}

		transferStore := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Transfer
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("unable to parse response. response: %q; error: '%v'", response.Body, err)
		}

		want := []Transfer{transfer1, transfer2, transfer3, transfer4}

		for i := range want {
			if want[i] != got[i] {
				t.Errorf("got %v; want %v", got[i], want[i])
				t.FailNow()
			}
		}

		assertHTTPStatus(t, response.Code, http.StatusOK)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return bad request when body is nil on POST", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(5087))
		transferStore := NewTransferStore(startingID(500))
		server := NewServer(accountStore, transferStore)

		request, _ := http.NewRequest(http.MethodPost, "/transfers", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
		assertResponseBody(t, response.Body.String(), "invalid request")
	})

	t.Run("should successfully transfer amount from origin account to destination account on POST", func(t *testing.T) {
		account1 := Account{
			ID:        207,
			Name:      "Juliana da Cruz Clemente",
			CPF:       "63000399003",
			Balance:   70000,
			CreatedAt: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		}
		account2 := Account{
			ID:        986,
			Name:      "Marlene de Souza Dalponte",
			CPF:       "08312653457",
			Balance:   51000,
			CreatedAt: time.Date(2020, time.February, 9, 15, 0, 0, 0, time.UTC),
		}

		accounts := map[uint64]Account{
			207: account1,
			986: account2,
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}

		transferStore := NewTransferStore(startingID(870))

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

		//assert balance
		wantBalance1 := uint64(66000)
		wantBalance2 := uint64(55000)
		gotBalance1 := accountStore.dataStorage[207].Balance
		gotBalance2 := accountStore.dataStorage[986].Balance
		assertUint64(t, gotBalance1, wantBalance1)
		assertUint64(t, gotBalance2, wantBalance2)

		//assert status
		wantTransferStatus := toStatusMsg(StatusConfirmed)
		gotTransferStatus := transferStore.dataStorage[871].Status
		assertString(t, gotTransferStatus, wantTransferStatus)

		//assert http
		assertResponseBody(t, gotTransferID, wantTransferID)
		assertHTTPStatus(t, response.Code, http.StatusCreated)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should return error if origin account ID is not found on POST", func(t *testing.T) {
		accountStore := NewAccountStore(startingID(90))
		transferStore := NewTransferStore(startingID(0))

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if destination account ID is not found", func(t *testing.T) {
		accounts := map[uint64]Account{
			307: {
				ID:        307,
				Name:      "Fernanda das Neves",
				CPF:       "78900167850",
				Balance:   205000,
				CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
			},
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}
		transferStore := NewTransferStore(startingID(0))

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if origin and destination account are the same", func(t *testing.T) {
		accounts := map[uint64]Account{
			307: {
				ID:        307,
				Name:      "Fernanda das Neves",
				CPF:       "78900167850",
				Balance:   205000,
				CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
			},
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}
		transferStore := NewTransferStore(startingID(0))

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if amount is invalid", func(t *testing.T) {
		accounts := map[uint64]Account{
			307: {
				ID:        307,
				Name:      "Fernanda das Neves",
				CPF:       "78900167850",
				Balance:   205000,
				CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
			},
			405: {
				ID:        405,
				Name:      "Marcela das Neves",
				CPF:       "78900389020",
				Balance:   780000,
				CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
			},
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}
		transferStore := NewTransferStore(startingID(0))

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if origin account balance is smaller than amount to be transferred", func(t *testing.T) {
		accounts := map[uint64]Account{
			307: {
				ID:        307,
				Name:      "Fernanda das Neves",
				CPF:       "78900167850",
				Balance:   205000,
				CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
			},
			405: {
				ID:        405,
				Name:      "Marcela das Neves",
				CPF:       "78900389020",
				Balance:   780000,
				CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
			},
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}

		transferStore := NewTransferStore(startingID(0))

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return error if transfer seems to be duplicated", func(t *testing.T) {
		accounts := map[uint64]Account{
			307: {
				ID:        307,
				Name:      "Fernanda das Neves",
				CPF:       "78900167850",
				Balance:   205000,
				CreatedAt: time.Date(2020, time.March, 05, 13, 0, 0, 0, time.UTC),
			},
			405: {
				ID:        405,
				Name:      "Marcela das Neves",
				CPF:       "78900389020",
				Balance:   780000,
				CreatedAt: time.Date(2020, time.March, 07, 11, 0, 0, 0, time.UTC),
			},
		}

		accountStore := &AccountStore{
			maxID:       startingID(len(accounts)),
			dataStorage: accounts,
		}

		transfer1 := Transfer{
			ID:                   1,
			AccountOriginID:      307,
			AccountDestinationID: 405,
			Amount:               15000,
			CreatedAt:            time.Now(),
			Status:               toStatusMsg(StatusConfirmed),
		}

		transfers := map[uint64]Transfer{
			1: transfer1,
		}

		transferStore := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

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

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		transferStore := NewTransferStore(startingID(109))
		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodDelete, "/transfers", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestTransfersID(t *testing.T) {
	t.Run("should return transfer depending on given ID on GET", func(t *testing.T) {
		transfers := map[uint64]Transfer{
			706: {
				ID:                   706,
				AccountOriginID:      45,
				AccountDestinationID: 78,
				Amount:               9500,
				CreatedAt:            time.Date(2020, time.March, 3, 11, 30, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusConfirmed),
			},
		}

		transferStore := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodGet, "/transfers/706", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := `{"id":706,"account_origin_id":45,"account_destination_id":78,"amount":9500,"created_at":"2020-03-03T11:30:00Z","status":"Confirmed"}`

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusOK)
		assertString(t, response.Result().Header.Get("content-type"), JsonContentType)
	})

	t.Run("should display error message if transfer ID is not found", func(t *testing.T) {
		transferStore := NewTransferStore(startingID(10))
		server := NewServer(nil, transferStore)

		inexistentID := 101

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfers/%v", inexistentID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`transfer %v not found`, inexistentID)

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("should display error message if transfer ID is invalid", func(t *testing.T) {
		transferStore := NewTransferStore(startingID(30))
		server := NewServer(nil, transferStore)

		invalidID := "letters"

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/transfers/%v", invalidID), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := fmt.Sprintf(`transfer ID is invalid. ID given: %v`, invalidID)

		assertResponseBody(t, got, want)
		assertHTTPStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("should return method not allowed to methods other than GET and POST", func(t *testing.T) {
		transferStore := NewTransferStore(startingID(109))
		server := NewServer(nil, transferStore)

		request, _ := http.NewRequest(http.MethodDelete, "/transfers/100", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertHTTPStatus(t, response.Code, http.StatusMethodNotAllowed)
	})
}
