package mnemonic

import (
	"errors"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

// Strength constants for mnemonic generation
const (
	Strength128 = 128 // 12 words
	Strength256 = 256 // 24 words
)

// Manager handles mnemonic generation and validation
type Manager struct{}

// NewManager creates a new mnemonic manager
func NewManager() *Manager {
	return &Manager{}
}

// Generate creates a new mnemonic phrase with the specified strength
// strength: 128 for 12 words, 256 for 24 words
func (m *Manager) Generate(strength int) (string, error) {
	if strength != Strength128 && strength != Strength256 {
		return "", errors.New("invalid strength: must be 128 or 256")
	}

	// Generate entropy
	entropy, err := bip39.NewEntropy(strength)
	if err != nil {
		return "", err
	}

	// Generate mnemonic from entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// Validate checks if a mnemonic phrase is valid
func (m *Manager) Validate(mnemonic string) error {
	if !bip39.IsMnemonicValid(mnemonic) {
		return errors.New("invalid mnemonic phrase")
	}
	return nil
}

// GetWordCount returns the number of words in a mnemonic
func (m *Manager) GetWordCount(mnemonic string) int {
	words := strings.Fields(strings.TrimSpace(mnemonic))
	return len(words)
}

// Normalize cleans up a mnemonic phrase
func (m *Manager) Normalize(mnemonic string) string {
	// Trim whitespace and normalize spaces
	words := strings.Fields(strings.TrimSpace(mnemonic))
	return strings.Join(words, " ")
}

// Package-level functions for convenience

// Generate creates a new mnemonic phrase with the specified strength
func Generate(strength int) (string, error) {
	manager := NewManager()
	return manager.Generate(strength)
}

// Validate checks if a mnemonic phrase is valid
func Validate(mnemonic string) error {
	manager := NewManager()
	return manager.Validate(mnemonic)
}

// GetWordCount returns the number of words in a mnemonic
func GetWordCount(mnemonic string) int {
	manager := NewManager()
	return manager.GetWordCount(mnemonic)
}

// Normalize cleans up a mnemonic phrase
func Normalize(mnemonic string) string {
	manager := NewManager()
	return manager.Normalize(mnemonic)
}
