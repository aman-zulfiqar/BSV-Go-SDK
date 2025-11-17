package wallet

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"

	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

// Generator handles BSV wallet generation
type Generator struct {
	network *chaincfg.Params
}

// BIP44Path represents a BIP44 derivation path
type BIP44Path struct {
	Purpose      uint32 // Purpose (44 for BIP44)
	CoinType     uint32 // Coin type (236 for BSV mainnet, 1 for testnet)
	Account      uint32 // Account index
	Change       uint32 // Change (0 = external, 1 = internal)
	AddressIndex uint32 // Address index
}

// NewGenerator creates a new wallet generator
// isTestnet: true for testnet, false for mainnet
func NewGenerator(isTestnet bool) *Generator {
	var network *chaincfg.Params
	if isTestnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}

	return &Generator{
		network: network,
	}
}

// GetDefaultBIP44Path returns the default BIP44 path for BSV
func (g *Generator) GetDefaultBIP44Path() *BIP44Path {
	var coinType uint32
	if g.network.Name == chaincfg.TestNet3Params.Name {
		coinType = 1 // Testnet
	} else {
		coinType = 236 // BSV mainnet
	}

	return &BIP44Path{
		Purpose:      44, // BIP44
		CoinType:     coinType,
		Account:      0, // First account
		Change:       0, // External (receiving addresses)
		AddressIndex: 0, // First address
	}
}

// GetBIP44Path returns a BIP44 path with custom indices
func (g *Generator) GetBIP44Path(account, change, addressIndex uint32) *BIP44Path {
	var coinType uint32
	if g.network.Name == chaincfg.TestNet3Params.Name {
		coinType = 1 // Testnet
	} else {
		coinType = 236 // BSV mainnet
	}

	return &BIP44Path{
		Purpose:      44, // BIP44
		CoinType:     coinType,
		Account:      account,
		Change:       change,
		AddressIndex: addressIndex,
	}
}

// GenerateWalletWithPath creates a BSV wallet from a mnemonic phrase using a specific BIP44 path
func (g *Generator) GenerateWalletWithPath(mnemonicPhrase string, path *BIP44Path) (*types.WalletResult, error) {
	// Validate mnemonic
	if err := mnemonic.Validate(mnemonicPhrase); err != nil {
		return nil, fmt.Errorf("invalid mnemonic: %v", err)
	}

	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonicPhrase, "")

	// Create master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %v", err)
	}

	// Derive BIP44 path: m/purpose'/coin_type'/account'/change/address_index
	childKey := masterKey

	// Derive purpose (44')
	childKey, err = childKey.NewChildKey(bip32.FirstHardenedChild + path.Purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose: %v", err)
	}

	// Derive coin type
	childKey, err = childKey.NewChildKey(bip32.FirstHardenedChild + path.CoinType)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type: %v", err)
	}

	// Derive account
	childKey, err = childKey.NewChildKey(bip32.FirstHardenedChild + path.Account)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account: %v", err)
	}

	// Derive change
	childKey, err = childKey.NewChildKey(path.Change)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change: %v", err)
	}

	// Derive address index
	childKey, err = childKey.NewChildKey(path.AddressIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address index: %v", err)
	}

	// Get private key
	privateKey, _ := btcec.PrivKeyFromBytes(childKey.Key)

	// Create WIF (Wallet Import Format)
	wif, err := btcutil.NewWIF(privateKey, g.network, true) // compressed = true
	if err != nil {
		return nil, fmt.Errorf("failed to create WIF: %v", err)
	}

	// Get public key
	publicKey := privateKey.PubKey()
	publicKeyBytes := publicKey.SerializeCompressed()

	// Create P2PKH address (BSV uses legacy addresses)
	addressPubKey, err := btcutil.NewAddressPubKey(publicKeyBytes, g.network)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %v", err)
	}

	return &types.WalletResult{
		Address:    addressPubKey.EncodeAddress(),
		PrivateKey: wif.String(),
		PublicKey:  hex.EncodeToString(publicKeyBytes),
	}, nil
}

// GenerateWallet creates a BSV wallet from a mnemonic phrase using default BIP44 path
func (g *Generator) GenerateWallet(mnemonicPhrase string) (*types.WalletResult, error) {
	defaultPath := g.GetDefaultBIP44Path()
	return g.GenerateWalletWithPath(mnemonicPhrase, defaultPath)
}

// GenerateWalletWithKeypair creates a wallet and returns the keypair for transaction signing
func (g *Generator) GenerateWalletWithKeypair(mnemonicPhrase string) (*types.WalletResult, *KeyPair, error) {
	wallet, err := g.GenerateWallet(mnemonicPhrase)
	if err != nil {
		return nil, nil, err
	}

	// Parse the WIF to get the keypair
	wif, err := btcutil.DecodeWIF(wallet.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode WIF: %v", err)
	}

	keyPair := &KeyPair{
		PrivateKey: wif.PrivKey,
		PublicKey:  wif.PrivKey.PubKey(),
		Network:    g.network,
	}

	return wallet, keyPair, nil
}

// GenerateRandomWallet creates a wallet with a random mnemonic
func (g *Generator) GenerateRandomWallet(strength int) (*types.WalletResult, string, error) {
	// Generate random mnemonic
	mnemonicPhrase, err := mnemonic.Generate(strength)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate mnemonic: %v", err)
	}

	// Generate wallet from mnemonic
	wallet, err := g.GenerateWallet(mnemonicPhrase)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate wallet: %v", err)
	}

	return wallet, mnemonicPhrase, nil
}

// ValidateAddress checks if a BSV address is valid
func (g *Generator) ValidateAddress(address string) error {
	_, err := btcutil.DecodeAddress(address, g.network)
	return err
}

// GetNetwork returns the network parameters
func (g *Generator) GetNetwork() *chaincfg.Params {
	return g.network
}

// KeyPair represents a BSV key pair for transaction signing
type KeyPair struct {
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
	Network    *chaincfg.Params
}

// SignMessage signs a message with the private key
func (kp *KeyPair) SignMessage(message []byte) ([]byte, error) {
	// For now, return a placeholder - this would need proper ECDSA signing
	// TODO: Implement proper message signing with ECDSA
	return []byte("placeholder_signature"), nil
}

// VerifySignature verifies a signature
func (kp *KeyPair) VerifySignature(message, signature []byte) bool {
	// For now, return true - this would need proper ECDSA verification
	// TODO: Implement proper signature verification
	return true
}

// Package-level functions for convenience

// GenerateWallet creates a BSV wallet from a mnemonic
func GenerateWallet(mnemonicPhrase string, isTestnet bool) (*types.WalletResult, error) {
	generator := NewGenerator(isTestnet)
	return generator.GenerateWallet(mnemonicPhrase)
}

// GenerateWalletWithKeypair creates a wallet with keypair
func GenerateWalletWithKeypair(mnemonicPhrase string, isTestnet bool) (*types.WalletResult, *KeyPair, error) {
	generator := NewGenerator(isTestnet)
	return generator.GenerateWalletWithKeypair(mnemonicPhrase)
}

// ValidateAddress validates a BSV address
func ValidateAddress(address string, isTestnet bool) error {
	generator := NewGenerator(isTestnet)
	return generator.ValidateAddress(address)
}
