package auth

import "testing"

func TestHashPassword(t *testing.T) {
	// 1. Hash not the same as original password
	input := "abc123password"
	actual, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	assertNotEqual(actual, input, input, t)

	// 2. Hash is not empty
	input = "abc123password"
	assertNotEqual(actual, "", input, t)

	// 3. Hash is salted (different hash for same passwords)
	input = "abc123"
	actual1, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	actual2, err := HashPassword(input)
	if err != nil {
		t.Errorf("%v", err)
	}
	assertNotEqual(actual1, actual2, input, t)
}

func TestCheckPassword(t *testing.T) {
	// Verify plaintext password matches hashed password
	plaintextPassword := "abc123password"
	hashedPassword, err := HashPassword(plaintextPassword)
	if err != nil {
		t.Errorf("%v", err)
	}

	assertEqual(CheckPasswordHash(plaintextPassword, hashedPassword), nil, nil, t)
}

func assertEqual(first, second, input any, t *testing.T) {
	if first != second {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, first, second)
	}
}

func assertNotEqual(first, second, input any, t *testing.T) {
	if first == second {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, first, second)
	}
}
