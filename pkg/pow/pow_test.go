package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHasher struct{}

func (m *mockHasher) Sum(data string) string {
	if data == "testChallenge000nonce" {
		return "00000abc"
	}
	return "abc12345"
}

func (s mockHasher) String() string {
	return "mock"
}

type mockGenerator struct{}

func (m *mockGenerator) Intn(n int) int {
	return 0
}

func TestGenerateChallenge(t *testing.T) {
	mockGen := &mockGenerator{}
	mockHasher := &mockHasher{}
	powInstance := New(5, mockHasher, mockGen)
	challenge := powInstance.GenerateChallenge()
	expected := "aaaaaaaaaaaaaaaa"

	assert.Equal(t, expected, challenge, "Expected challenge to match")
}

func TestVerify(t *testing.T) {
	mockGen := &mockGenerator{}
	mockHasher := &mockHasher{}

	tests := []struct {
		name       string
		challenge  string
		nonce      string
		difficulty int
		expected   bool
	}{
		{
			name:       "valid nonce with difficulty 5",
			challenge:  "testChallenge",
			nonce:      "000nonce",
			difficulty: 5,
			expected:   true,
		},
		{
			name:       "invalid nonce with difficulty 5",
			challenge:  "testChallenge",
			nonce:      "wrongNonce",
			difficulty: 5,
			expected:   false,
		},
		{
			name:       "valid nonce with difficulty 1",
			challenge:  "testChallenge",
			nonce:      "000nonce",
			difficulty: 1,
			expected:   true,
		},
		{
			name:       "invalid nonce with difficulty 1",
			challenge:  "testChallenge",
			nonce:      "wrongNonce",
			difficulty: 1,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			powInstance := New(tt.difficulty, mockHasher, mockGen)
			result := powInstance.Verify(tt.challenge, tt.nonce)
			assert.Equal(t, tt.expected, result, "Expected result to match for challenge %s and nonce %s", tt.challenge, tt.nonce)
		})
	}
}
