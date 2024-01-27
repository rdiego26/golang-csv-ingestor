package server

import "testing"

func TestParseMilliseconds(t *testing.T) {
	input := "1527867048000"
	want := 1527867048000

	got, err := parseMilliseconds(input)
	if err != nil {
		t.Errorf("got %q, wanted %q", got, want)
	}

	wrongInput := "A"
	result, err := parseMilliseconds(wrongInput)
	if result != 0 {
		t.Errorf("result %q, wanted %q", result, 0)
	}

	if err == nil {
		t.Errorf("expects err non nil but received %v", err)
	}
}
