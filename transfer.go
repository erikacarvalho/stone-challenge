package app

import (
	"errors"
	"sort"
	"sync/atomic"
	"time"
)

const (
	StatusCreated       = 1
	StatusAuthorizing   = 2
	StatusNotAuthorized = 3
	StatusAuthorized    = 4
	StatusCancelled     = 5
	StatusConfirmed     = 6
)

var statusMessage = map[int]string{
	StatusCreated:       "Created",
	StatusAuthorizing:   "Authorizing",
	StatusNotAuthorized: "Not Authorized",
	StatusAuthorized:    "Authorized",
	StatusCancelled:     "Cancelled",
	StatusConfirmed:     "Confirmed",
}

var (
	ErrInsufficientBalance = errors.New("origin account balance is too low to allow this transfer")
	ErrSameID              = errors.New("origin and destination account ids are the same")
	ErrChargeBack          = errors.New("this transfer seems to be duplicated")
	ErrInvalidAmount       = errors.New("the amount entered is invalid")
	ErrNoTransfers         = errors.New("there are no transfers to be listed")
	ErrTransferNotFound    = errors.New("there is no transfer with this ID")
)

type TransferStore struct {
	maxID       *uint64
	dataStorage map[uint64]Transfer // The map key is the transfer identifier
}

// NewTransferStore generates a new TransferStore with a starting ID number and
// returns it.
func NewTransferStore(startingID *uint64) *TransferStore {
	ns := &TransferStore{
		maxID:       startingID,
		dataStorage: make(map[uint64]Transfer),
	}
	return ns
}

// CreateTransfer is a method that creates a transfer based on origin
// and destination account ids and an amount, and returns an incrementally
// generated ID. It also sets created time to Now and status to Created.
func (t *TransferStore) CreateTransfer(origin, destination, amount uint64) (id uint64, err error) {
	newID := atomic.AddUint64(t.maxID, 1)
	t.dataStorage[newID] = Transfer{
		ID:                   newID,
		AccountOriginID:      origin,
		AccountDestinationID: destination,
		Amount:               amount,
		CreatedAt:            time.Now(),
		Status:               toStatusMsg(StatusCreated),
	}
	return newID, nil
}

// authorizeTransfer checks if it is possible to perform the transfer
// based on the business rules, and returns error message depending on
// the outcome.
func (t *TransferStore) authorizeTransfer(origin, destination *Account, amount, id uint64) error {
	changeStatus(t, id, StatusAuthorizing)

	if origin.ID == destination.ID {
		changeStatus(t, id, StatusNotAuthorized)
		return ErrSameID
	}

	if amount == 0 {
		changeStatus(t, id, StatusNotAuthorized)
		return ErrInvalidAmount
	}

	if origin.Balance < amount {
		changeStatus(t, id, StatusNotAuthorized)
		return ErrInsufficientBalance
	}

	if t.isChargeBack(origin.ID, destination.ID, amount) {
		changeStatus(t, id, StatusNotAuthorized)
		return ErrChargeBack
	}
	changeStatus(t, id, StatusAuthorized)
	return nil
}

// isChargeBack indicates if there is a presumable processing error
// creating a duplicated transfer within a short period of time (set
// to 10 seconds by default). It helps avoiding being chargedbacked.
func (t *TransferStore) isChargeBack(origin, destination, amount uint64) bool {
	for _, transfer := range t.dataStorage {
		treshold := transfer.CreatedAt.Add(10 * time.Second)
		now := time.Now()
		if transfer.AccountOriginID == origin &&
			transfer.AccountDestinationID == destination &&
			transfer.Amount == amount &&
			transfer.Status == toStatusMsg(StatusConfirmed) &&
			now.Before(treshold) {
			return true
		}
	}
	return false
}

// Confirm sets the transfer status to confirmed.
func (t *TransferStore) Confirm(id uint64) {
	changeStatus(t, id, StatusConfirmed)
}

// Cancel sets the transfer status to cancelled.
func (t *TransferStore) Cancel(id uint64) {
	changeStatus(t, id, StatusCancelled)
}

// ListAllTransfers returns all transfers from the store sorted by ID,
// and an error if there are no transfers to be listed.
func (t *TransferStore) ListAllTransfers() ([]Transfer, error) {
	var transfers []Transfer
	for _, v := range t.dataStorage {
		transfers = append(transfers, v)
	}

	if len(transfers) == 0 {
		return nil, ErrNoTransfers
	}

	sort.Slice(transfers, func(i, j int) bool {
		return transfers[i].ID < transfers[j].ID
	})

	return transfers, nil
}

// GetTransfer returns a Transfer based on a given ID, and an error if
// no transfer with given ID is found.
func (t *TransferStore) GetTransfer(ID uint64) (Transfer, error) {
	transfer, ok := t.dataStorage[ID]
	if !ok {
		return Transfer{}, ErrTransferNotFound
	}
	return transfer, nil
}

func changeStatus(a *TransferStore, ID uint64, statusCode int) {
	transfer := a.dataStorage[ID]
	transfer.Status = toStatusMsg(statusCode)
	a.dataStorage[ID] = transfer
}

func toStatusMsg(code int) string {
	return statusMessage[code]
}
