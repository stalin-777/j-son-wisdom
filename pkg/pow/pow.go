package pow

import (
	"fmt"
	"strings"
	"time"

	"github.com/stalin-777/j-son-wisdom/pkg/hash"

	"golang.org/x/exp/rand"
)

// RandomGenerator defines an interface for generating random numbers.
type RandomGenerator interface {
	Intn(n int) int
}

type defaultGenerator struct{}

func (d *defaultGenerator) Intn(n int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return rand.Intn(n)
}

// PoW provides a framework for dealing with proof-of-work tasks.
type PoW struct {
	rg         RandomGenerator
	hasher     hash.Hasher
	difficulty int
}

// New creates a new PoW instance with the given difficulty and random number generator.
// If no generator is specified, the default generator is used.
func New(difficulty int, hasher hash.Hasher, rg RandomGenerator) *PoW {
	if rg == nil {
		rg = &defaultGenerator{}
	}

	return &PoW{
		rg:         rg,
		hasher:     hasher,
		difficulty: difficulty,
	}
}

// GeneratePoWChallenge generates a random 16-character string for proof of work challenge
func (p *PoW) GenerateChallenge() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)

	for i := range result {
		result[i] = charset[p.rg.Intn(len(charset))]
	}

	return string(result)
}

// VerifyPoW checks if the given nonce matches the proof of work challenge
func (p *PoW) Verify(challenge, nonce string) bool {
	hashStr := p.generateHash(challenge, nonce)
	return strings.HasPrefix(hashStr, strings.Repeat("0", p.difficulty))
}

func (p *PoW) generateHash(challenge, nonce string) string {
	guess := fmt.Sprintf("%s%s", challenge, nonce)
	hash := p.hasher.Sum(guess)
	return hash
}
