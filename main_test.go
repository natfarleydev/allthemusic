package main

import "testing"

func TestIncrString(t *testing.T) {
	x := "aaa"
	y, _ := incrString(x)
	if y != "baa" {
		t.Errorf("y is not equal to baa, is %q", y)
	}
}

func TestIncrStringWillRollover(t *testing.T) {
	x := "zaa"
	y, _ := incrString(x)
	expectedY := "aba"
	if y != expectedY {
		t.Errorf("y is not equal to %q, is %q", expectedY, y)
	}
}

func TestIncrStringWillAddNewChar(t *testing.T) {
	x := "zzz"
	y, err := incrString(x)
	if err != nil {
		t.Errorf("Unexpected error: %q", err)
	}
	expectedY := "aaaa"
	if y != expectedY {
		t.Errorf("y is not equal to %q, is %q", expectedY, y)
	}
}

func TestIncrStringErrorsOnInvalidChar(t *testing.T) {
	x := "111"
	y, err := incrString(x)
	if err == nil {
		t.Errorf("Expected error, got: %q, %q", y, err)
	}
}

func BenchmarkIncrString4chars(b *testing.B) {
	for n := 0; n < b.N; n++ {
		incrString("aaaa")
	}
}
func BenchmarkIncrString10chars(b *testing.B) {
	for n := 0; n < b.N; n++ {
		incrString("aaaaaaaaaa")
	}
}
