package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// Hasher defines the interface for hashing algorithms
type Hasher interface {
	Sum(data string) string
}

// Sha1 represents the SHA-1 hashing algorithm
type Sha1 struct{}

// Sum calculates the SHA-1 hash of the input data
func (s Sha1) Sum(data string) string {
	hash := sha1.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}

// Sha256 represents the SHA-256 hashing algorithm
type Sha256 struct{}

// Sum calculates the SHA-256 hash of the input data
func (s Sha256) Sum(data string) string {
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}

var algorithms = map[string]Hasher{
	"sha1":   Sha1{},
	"sha256": Sha256{},
}

// GetAlgorithm retrieves a hasher based on the algorithm name
func GetAlgorithm(name string) (Hasher, bool) {
	algorithm, ok := algorithms[name]
	return algorithm, ok
}

// RegisterAlgorithm registers a new hash algorithm
func RegisterAlgorithm(name string, algo Hasher) {
	algorithms[name] = algo
}
