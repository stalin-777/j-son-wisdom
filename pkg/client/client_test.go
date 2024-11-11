package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHasher struct{}

func (m *mockHasher) Sum(input string) string {
	return "0fixed result"
}

func TestSolvePoW(t *testing.T) {
	challenge := "test"
	difficulty := 1
	hasher := mockHasher{}
	result := solvePoW(challenge, &hasher, difficulty)
	assert.Equal(t, "0", result)
}
