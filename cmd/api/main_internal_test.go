// Filename: cmd/api/main_test.go

package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	want := "Hello, UB!"
	got := printUB()

	if got != want {
		t.Errorf("expected: %q, got: %q", want, got)
	}
}