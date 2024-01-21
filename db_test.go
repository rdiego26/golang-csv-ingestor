package main

import "testing"

func TestFilterParams(t *testing.T) {
	input := make(map[string]string)
	input["unknown"] = "unknown"
	input["email"] = "email"

	want := make(map[string]string)
	want["email"] = "email"

	got := filterParams(input)
	if len(got) == len(input) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
