package assert

import "testing"

func Equal[T comparable](t testing.TB, expected, actual T) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected: %v; got: %v", expected, actual)
	}
}
