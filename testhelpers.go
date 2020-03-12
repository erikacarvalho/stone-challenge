package app

import "testing"

func AssertError(t *testing.T, got, want error) {
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body has a problem. got %q; want %q", got, want)
	}
}

func AssertHTTPStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("http status is incorrect. got %d; want %d", got, want)
	}
}
func AssertString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got response %q; want %q", got, want)
	}
}

func AssertUint64(t *testing.T, got, want uint64) {
	t.Helper()
	if want != got {
		t.Errorf("got %d; want %d", got, want)
	}
}

func StartingID(ID int) *uint64 {
	var ptr = uint64(ID)
	return &ptr
}
