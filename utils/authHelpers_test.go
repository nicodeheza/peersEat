package utils

import (
	"testing"
)

func TestGetRandomWords(t *testing.T) {

	a := &AuthHelpers{}

	result1, _ := a.GetRandomWords(3)
	result2, _ := a.GetRandomWords(3)

	if result1 == result2 {
		t.Error("same words")
	}
}

func TestGetRandomPassword(t *testing.T) {

	a := &AuthHelpers{}

	result1 := a.GetRandomPassword(20)
	result2 := a.GetRandomPassword(20)

	if result1 == result2 {
		t.Error("same password")
	}
}
