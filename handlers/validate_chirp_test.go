package handlers

import (
	"testing"
)

func TestCensoredBannedWords(t *testing.T) {
	input := "What the sharbert? This is fornax crazy. Really, it's FORNAXing crazy kerfuffle"
	expected := "What the sharbert? This is **** crazy. Really, it's FORNAXing crazy ****"

	actual := censoredBannedWords(input)

	assertEqual(expected, actual, input, t)
}

func assertEqual(expected, actual, input any, t *testing.T) {
	if actual != expected {
		t.Errorf("\nInput:\n\t%v\nActual:\n\t%v\nExpected:\n\t%v", input, expected, actual)
	}
}
