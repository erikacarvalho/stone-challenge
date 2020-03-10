package app

import (
	"encoding/json"
	"log"
	"net/http"
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

type AccountServer struct {
	store AccountStore
	http.Handler
}

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

	router := http.NewServeMux()
	router.Handle("/accounts", http.HandlerFunc(p.accountsHandler))

	p.Handler = router

	return p
}
