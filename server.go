package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const JsonContentType = "application/json"

type CreateAccountRequest struct {
	Name    string `json:"name"`
	CPF     string `json:"cpf"`
	Balance uint64 `json:"balance"`
}

type CreateAccountResponse struct {
	ID uint64 `json:"id"`
}

type GetBalanceResponse struct {
	ID uint64 `json:"id"`
	Balance uint64 `json:"balance"`
}

type AccountServer struct {
	store AccountStore
	http.Handler
}

// accountsHandler redirects '/accounts' endpoint requests to their
// proper Handler depending on the HTTP method.
func (a *AccountServer) accountsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.list(w)
	case http.MethodPost:
		a.add(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// balanceHandler responds with balance to a given account ID.
func (a *AccountServer) balanceHandler(w http.ResponseWriter, r *http.Request) {
	varsMap := mux.Vars(r)
	idStr, ok := varsMap["account_id"]
	if !ok {
		log.Println("it was impossible to obtain the ID from the path")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("account ID is invalid. ID given: %v", idStr)
		log.Println(errMsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errMsg))
		return
	}

	balance, err := a.store.GetBalance(ID)
	if err == ErrAccountNotFound {
		errMsg := fmt.Sprintf("account %v not found", ID)
		log.Println(errMsg)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errMsg))
		return
	}
	jsonBytes, err := json.Marshal(GetBalanceResponse{
		ID:      ID,
		Balance: balance,
	})
	if err != nil {
		log.Printf("error marshaling balance: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", JsonContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// list returns the list of all accounts.
func (a *AccountServer) list(w http.ResponseWriter) {
	getList, err := a.store.ListAll()

	if err == ErrNoRecords {
		w.Header().Set("content-type", JsonContentType)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	err = json.NewEncoder(w).Encode(getList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", JsonContentType)
	w.WriteHeader(http.StatusOK)
}

// add creates a new account based on a CreateAccountRequest and returns
// its ID.
func (a *AccountServer) add(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		log.Println("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	creationRequest := CreateAccountRequest{}

	err := json.NewDecoder(r.Body).Decode(&creationRequest)
	if err != nil {
		log.Printf("error decoding body to CreateAccountRequest: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request"))
		return
	}

	newAccID := a.store.Create(creationRequest.Name, creationRequest.CPF, creationRequest.Balance)
	jsonBytes, err := json.Marshal(CreateAccountResponse{ID: newAccID})
	if err != nil {
		log.Printf("error marshaling new account ID: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", JsonContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func NewAccountServer(store *AccountStore) *AccountServer {
	p := &AccountServer{store: *store}

	router := mux.NewRouter()

	router.HandleFunc("/accounts", p.accountsHandler)
	router.HandleFunc("/accounts/{account_id}/balance", p.balanceHandler)

	p.Handler = router

	return p
}