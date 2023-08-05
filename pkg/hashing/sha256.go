package hashing

import (
	"crypto/sha256"
	"fmt"
)

type SHA256Hasher struct {
	Salt string
}

func NewSHA256Hasher(salt string) *SHA256Hasher {
	return &SHA256Hasher{Salt: salt}
}

func (h *SHA256Hasher) Hash(password string) string {

	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.Salt)))
}
