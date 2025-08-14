package main

import "testing"

func TestUnpackString_Basic(t *testing.T) {
	input := "a4bc2d5e"
	expected := "aaaabccddddde"

	result, err := unpackString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestUnpackString_WithoutDigits(t *testing.T) {
	input := "abcd"
	expected := "abcd"

	result, err := unpackString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestUnpackString_StartsWithDigit(t *testing.T) {
	input := "45"
	expectedErr := "string '45' contains pair of digits or does not contain non-digit letters. to avoid that use escape-sequence"

	_, err := unpackString(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestUnpackString_ContainsDigitPair(t *testing.T) {
	input := "abc45"
	expectedErr := "string 'abc45' contains pair of digits or does not contain non-digit letters. to avoid that use escape-sequence"

	_, err := unpackString(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestUnpackString_EmptyString(t *testing.T) {
	input := ""
	expected := ""

	result, err := unpackString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestUnpackString_EscapeBasic(t *testing.T) {
	input := `a\3c`
	expected := "a3c"

	result, err := unpackString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestUnpackString_EscapeWithRepeat(t *testing.T) {
	input := `a\33c`
	expected := "a333c"

	result, err := unpackString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
