package app

import "testing"

func assertError(t *testing.T, got, want error) {
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body has a problem. got %q; want %q", got, want)
	}
}

func assertHTTPStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("http status is incorrect. got %d; want %d", got, want)
	}
}
func assertString(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got response %q; want %q", got, want)
	}
}

func assertUint64(t *testing.T, got, want uint64) {
	t.Helper()
	if want != got {
		t.Errorf("got %d; want %d", got, want)
	}
}

func startingID(ID int) *uint64 {
	var ptr = uint64(ID)
	return &ptr
}
