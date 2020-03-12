package main

import (
	http2 "github.com/erikacarvalho/stone-challenge/http"
	"github.com/erikacarvalho/stone-challenge/store"
	"log"
	"net/http"
)

var (
	accountStoreStartingID  = uint64(0)
	transferStoreStartingID = uint64(0)
)

const address = ":3000"

func main() {
	log.Println("initializing server on", address)
	server := http2.NewServer(
		store.NewAccountStore(&accountStoreStartingID),
		store.NewTransferStore(&transferStoreStartingID),
	)
	log.Fatal(http.ListenAndServe(address, server))
}
