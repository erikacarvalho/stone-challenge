package app

import (
	"testing"
	"time"
)

func TestCreateTransfer(t *testing.T) {
	store := NewTransferStore(startingID(67))

	origin := uint64(7)
	destination := uint64(30)
	amount := uint64(1500)

	newTransferID, err := store.CreateTransfer(origin, destination, amount)

	if err != nil {
		t.Errorf("error creating transfer. error: %q", err)
	}

	t.Run("should return new transfer ID based on origin and destination account ids and amount", func(t *testing.T) {
		wantID, gotID := uint64(68), newTransferID
		wantOriginID, gotOriginID := origin, store.dataStorage[newTransferID].AccountOriginID
		wantDestinationID, gotDestinationID := destination, store.dataStorage[newTransferID].AccountDestinationID
		wantAmount, gotAmount := amount, store.dataStorage[newTransferID].Amount

		assertUint64(t, gotID, wantID)
		assertUint64(t, gotOriginID, wantOriginID)
		assertUint64(t, gotDestinationID, wantDestinationID)
		assertUint64(t, gotAmount, wantAmount)
	})

	t.Run("should create new transfer with status code Created", func(t *testing.T) {
		want := toStatusMsg(StatusCreated)
		got := store.dataStorage[newTransferID].Status

		assertString(t, got, want)
	})
}

func TestConfirm(t *testing.T) {
	t.Run("should change transfer status to Confirmed", func(t *testing.T) {
		store := NewTransferStore(startingID(93))

		newID, _ := store.CreateTransfer(15, 21, 7800)

		store.Confirm(newID)

		want := toStatusMsg(StatusConfirmed)
		got := store.dataStorage[newID].Status

		assertString(t, got, want)
	})
}

func TestCancel(t *testing.T) {
	t.Run("should change transfer status to StatusCancelled", func(t *testing.T) {
		store := NewTransferStore(startingID(800))

		newID, _ := store.CreateTransfer(90, 2, 19000)

		store.Cancel(newID)

		want := toStatusMsg(StatusCancelled)
		got := store.dataStorage[newID].Status

		assertString(t, got, want)
	})
}

func TestListAllTransfers(t *testing.T) {
	t.Run("should return slice with all stored transfers", func(t *testing.T) {
		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      78,
				AccountDestinationID: 990,
				Amount:               15000,
				CreatedAt:            time.Date(2020, time.February, 15, 8, 0, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusConfirmed),
			},
			2: {
				ID:                   2,
				AccountOriginID:      501,
				AccountDestinationID: 97,
				Amount:               60000,
				CreatedAt:            time.Date(2020, time.February, 16, 10, 0, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusNotAuthorized),
			},
			3: {
				ID:                   3,
				AccountOriginID:      501,
				AccountDestinationID: 97,
				Amount:               5000,
				CreatedAt:            time.Date(2020, time.February, 16, 19, 0, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusConfirmed),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		transfersList, _ := store.ListAllTransfers()

		for i, transfer := range transfersList {
			want := transfers[uint64(i+1)]
			got := transfer
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		}
	})

	t.Run("should return ErrNoTransfers if there are no transfers", func(t *testing.T) {
		store := NewTransferStore(startingID(0))

		want := ErrNoTransfers
		_, got := store.ListAllTransfers()

		assertError(t, got, want)
	})
}

func TestGetTransfer(t *testing.T) {
	t.Run("should return Transfer for a given ID", func(t *testing.T) {
		transfers := map[uint64]Transfer{
			9: {
				ID:                   9,
				AccountOriginID:      55,
				AccountDestinationID: 411,
				Amount:               40000,
				CreatedAt:            time.Date(2020, time.February, 9, 10, 0, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusConfirmed),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		want := transfers[9]
		got, _ := store.GetTransfer(9)

		if got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})

	t.Run("should return ErrTransferNotFound when there is no account for given ID", func(t *testing.T) {
		store := NewTransferStore(startingID(7))

		want := ErrTransferNotFound
		_, got := store.GetTransfer(167)

		assertError(t, got, want)
	})
}

func TestChargeBack(t *testing.T) {
	originID := uint64(78)
	destinationID := uint64(990)
	amount := uint64(15000)

	t.Run("should indicate chargeback when threshold time is not over", func(t *testing.T) {
		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusConfirmed),
			},
			2: {
				ID:                   2,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusAuthorizing),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		want := true
		got := store.isChargeBack(originID, destinationID, amount)

		if got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})

	t.Run("should not indicate chargeback when threshold time is over", func(t *testing.T) {
		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now().Add(-11 * time.Second),
				Status:               toStatusMsg(StatusConfirmed),
			},
			2: {
				ID:                   2,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusAuthorizing),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		want := false
		got := store.isChargeBack(originID, destinationID, amount)

		if got != want {
			t.Errorf("got %v; want %v", got, want)
		}
	})

}

func TestAuthorize(t *testing.T) {
	t.Run("should return no error if transfer is authorized", func(t *testing.T) {
		originID := uint64(15)
		destinationID := uint64(87)
		amount := uint64(1000)

		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Date(2020, time.March, 10, 7, 0, 0, 0, time.UTC),
				Status:               toStatusMsg(StatusConfirmed),
			},
			2: {
				ID:                   2,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusCreated),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		origin := &Account{
			ID:      originID,
			Balance: 5000,
		}
		destination := &Account{
			ID:      destinationID,
			Balance: 9000,
		}
		got := store.authorizeTransfer(origin, destination, amount, 2)

		gotStatus := store.dataStorage[2].Status
		wantStatus := toStatusMsg(StatusAuthorized)

		assertError(t, got, nil)
		assertString(t, gotStatus, wantStatus)
	})

	t.Run("should return ErrInvalidAmount when amount to be transferred is zero", func(t *testing.T) {
		var amount uint64 = 0
		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      207,
				AccountDestinationID: 986,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusCreated),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		origin := &Account{
			ID:      207,
			Balance: 5000,
		}
		destination := &Account{
			ID:      986,
			Balance: 9000,
		}

		want := ErrInvalidAmount
		got := store.authorizeTransfer(origin, destination, amount, 1)

		gotStatus := store.dataStorage[1].Status
		wantStatus := toStatusMsg(StatusNotAuthorized)

		assertError(t, got, want)
		assertString(t, gotStatus, wantStatus)
	})

	t.Run("should return ErrSameID when origin and destination account ids are the same", func(t *testing.T) {
		amount := uint64(4000)

		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      207,
				AccountDestinationID: 207,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusCreated),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		acc := &Account{
			ID: 207,
		}

		want := ErrSameID
		got := store.authorizeTransfer(acc, acc, amount, 1)

		gotStatus := store.dataStorage[1].Status
		wantStatus := toStatusMsg(StatusNotAuthorized)

		assertError(t, got, want)
		assertString(t, gotStatus, wantStatus)
	})

	t.Run("should return ErrInsufficientBalance when origin account balance is insufficient", func(t *testing.T) {
		amount := uint64(5500)

		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      207,
				AccountDestinationID: 986,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusCreated),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		origin := &Account{
			ID:      207,
			Balance: 5000,
		}
		destination := &Account{
			ID:      986,
			Balance: 9000,
		}
		want := ErrInsufficientBalance
		got := store.authorizeTransfer(origin, destination, amount, 1)

		gotStatus := store.dataStorage[1].Status
		wantStatus := toStatusMsg(StatusNotAuthorized)

		assertError(t, got, want)
		assertString(t, gotStatus, wantStatus)
	})

	t.Run("should return ErrChargeBack when it seems to be a duplicated transfer", func(t *testing.T) {
		originID := uint64(15)
		destinationID := uint64(87)
		amount := uint64(1000)

		transfers := map[uint64]Transfer{
			1: {
				ID:                   1,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusConfirmed),
			},
			2: {
				ID:                   2,
				AccountOriginID:      originID,
				AccountDestinationID: destinationID,
				Amount:               amount,
				CreatedAt:            time.Now(),
				Status:               toStatusMsg(StatusCreated),
			},
		}

		store := &TransferStore{
			maxID:       startingID(len(transfers)),
			dataStorage: transfers,
		}

		origin := &Account{
			ID:      originID,
			Balance: 5000,
		}
		destination := &Account{
			ID:      destinationID,
			Balance: 9000,
		}

		want := ErrChargeBack
		got := store.authorizeTransfer(origin, destination, amount, 2)

		gotStatus := store.dataStorage[2].Status
		wantStatus := toStatusMsg(StatusNotAuthorized)

		assertError(t, got, want)
		assertString(t, gotStatus, wantStatus)
	})
}
