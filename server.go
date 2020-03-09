package app

import (
	"encoding/json"
	"net/http"
)

const JsonContentType = "application/json"

type AccountServer struct {
	store AccountStore
	http.Handler
}

func(a *AccountServer) accountsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.list(w)
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
	}
	w.Header().Set("content-type", JsonContentType)
	w.WriteHeader(http.StatusOK)
}

func NewAccountServer(store *AccountStore) *AccountServer {
	p := &AccountServer{store: *store}

	router := http.NewServeMux()
	router.Handle("/accounts", http.HandlerFunc(p.accountsHandler))

	p.Handler = router

	return p
}