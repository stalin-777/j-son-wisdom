package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha1_Sum(t *testing.T) {
	hasher := Sha1{}
	input := "test"
	expected := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	result := hasher.Sum(input)
	assert.Equal(t, expected, result)
}

func TestSha256_Sum(t *testing.T) {
	hasher := Sha256{}
	input := "test"
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"

	result := hasher.Sum(input)
	assert.Equal(t, expected, result)
}

func TestGetAlgorithm(t *testing.T) {
	tests := []struct {
		algo     string
		expectOk bool
	}{
		{"sha1", true},
		{"sha256", true},
		{"unknown", false},
	}

	for _, tt := range tests {
		_, ok := GetAlgorithm(tt.algo)
		assert.Equal(t, tt.expectOk, ok, "Algorithm %s. Expected result %v, received %v", tt.algo, tt.expectOk, ok)
	}
}

type unknown struct{}

func (u unknown) Sum(data string) string {
	return ""
}

func TestRegisterAlgorithm(t *testing.T) {
	name := "unknown"

	RegisterAlgorithm(name, unknown{})

	_, ok := GetAlgorithm(name)
	assert.Equal(t, true, ok)
}
