package bsv

import (
	"fmt"

	"github.com/muhammadamman/BSV-Go/pkg/bsv/transaction"
	"github.com/muhammadamman/BSV-Go/pkg/bsv/wallet"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

// BSV provides the main interface for BSV operations with dynamic configuration
type BSV struct {
	configManager *config.Manager
	walletGen     *wallet.Generator
	txBuilder     *transaction.Builder
}

// NewBSV creates a new BSV instance
func NewBSV(configManager *config.Manager) *BSV {
	networkConfig := configManager.GetNetworkConfig()

	return &BSV{
		configManager: configManager,
		walletGen:     wallet.NewGenerator(networkConfig.IsTestnet),
		txBuilder:     transaction.NewBuilder(configManager),
	}
}

// NewBSVWithNetwork creates a new BSV instance with network type
func NewBSVWithNetwork(networkType config.NetworkType) (*BSV, error) {
	configManager := config.NewManager()

	if err := configManager.SetNetworkType(networkType); err != nil {
		return nil, fmt.Errorf("failed to set network type: %v", err)
	}

	return NewBSV(configManager), nil
}

// GenerateWallet creates a BSV wallet from a mnemonic phrase
func (b *BSV) GenerateWallet(mnemonicPhrase string) (*types.WalletResult, error) {
	return b.walletGen.GenerateWallet(mnemonicPhrase)
}

// GenerateWalletWithPath creates a BSV wallet using a specific BIP44 path
func (b *BSV) GenerateWalletWithPath(mnemonicPhrase string, account, change, addressIndex uint32) (*types.WalletResult, error) {
	path := b.walletGen.GetBIP44Path(account, change, addressIndex)
	return b.walletGen.GenerateWalletWithPath(mnemonicPhrase, path)
}

// GetBIP44Path returns a BIP44 path with custom indices
func (b *BSV) GetBIP44Path(account, change, addressIndex uint32) *wallet.BIP44Path {
	return b.walletGen.GetBIP44Path(account, change, addressIndex)
}

// GetDefaultBIP44Path returns the default BIP44 path
func (b *BSV) GetDefaultBIP44Path() *wallet.BIP44Path {
	return b.walletGen.GetDefaultBIP44Path()
}

// GenerateWalletWithKeypair creates a wallet with keypair for transaction signing
func (b *BSV) GenerateWalletWithKeypair(mnemonicPhrase string) (*types.WalletResult, *wallet.KeyPair, error) {
	return b.walletGen.GenerateWalletWithKeypair(mnemonicPhrase)
}

// GenerateRandomWallet creates a wallet with a random mnemonic
func (b *BSV) GenerateRandomWallet(strength int) (*types.WalletResult, string, error) {
	return b.walletGen.GenerateRandomWallet(strength)
}

// ValidateAddress validates a BSV address
func (b *BSV) ValidateAddress(address string) error {
	return b.walletGen.ValidateAddress(address)
}

// GetEnhancedBalance retrieves enhanced balance information for an address
func (b *BSV) GetEnhancedBalance(address string) (*types.EnhancedBalanceInfo, error) {
	return b.txBuilder.GetEnhancedBalance(address)
}

// GetNativeBalance retrieves native BSV balance for an address
func (b *BSV) GetNativeBalance(address string) (*types.NativeBalanceInfo, error) {
	return b.txBuilder.GetNativeBalance(address)
}

// GetNonNativeBalance retrieves non-native token balances for an address
func (b *BSV) GetNonNativeBalance(address string) (*types.NonNativeBalanceInfo, error) {
	return b.txBuilder.GetNonNativeBalance(address)
}

// GetBalance retrieves balance for an address (backward compatibility)
func (b *BSV) GetBalance(address string) (int64, error) {
	return b.txBuilder.GetBalance(address)
}

// GetUTXOs retrieves UTXOs for an address
func (b *BSV) GetUTXOs(address string) ([]types.UTXO, error) {
	return b.txBuilder.GetUTXOs(address)
}

// BuildTransaction builds a BSV transaction with enhanced support
func (b *BSV) BuildTransaction(params *types.TransactionParams) (*types.TransactionResult, error) {
	tx, err := b.txBuilder.BuildTransaction(params)
	if err != nil {
		return nil, err
	}

	// Convert to result format
	return &types.TransactionResult{
		SignedTx: fmt.Sprintf("%x", tx),
	}, nil
}

// SignAndSendTransaction builds, signs, and broadcasts a transaction
func (b *BSV) SignAndSendTransaction(params *types.TransactionParams) (*types.TransactionResult, error) {
	return b.txBuilder.SignAndSendTransaction(params)
}

// GetNetwork returns whether this is testnet
func (b *BSV) GetNetwork() bool {
	networkConfig := b.configManager.GetNetworkConfig()
	return networkConfig.IsTestnet
}

// GetNetworkConfig returns the current network configuration
func (b *BSV) GetNetworkConfig() *config.NetworkConfig {
	return b.configManager.GetNetworkConfig()
}

// GetUTXOConfig returns the current UTXO configuration
func (b *BSV) GetUTXOConfig() *config.UTXOConfig {
	return b.configManager.GetUTXOConfig()
}

// GetTransactionConfig returns the current transaction configuration
func (b *BSV) GetTransactionConfig() *config.TransactionConfig {
	return b.configManager.GetTransactionConfig()
}

// UpdateNetworkConfig updates the network configuration
func (b *BSV) UpdateNetworkConfig(network *config.NetworkConfig) error {
	err := b.configManager.UpdateNetworkConfig(network)
	if err != nil {
		return err
	}

	// Update wallet generator with new network
	networkConfig := b.configManager.GetNetworkConfig()
	b.walletGen = wallet.NewGenerator(networkConfig.IsTestnet)

	return nil
}

// UpdateUTXOConfig updates the UTXO configuration
func (b *BSV) UpdateUTXOConfig(utxo *config.UTXOConfig) error {
	return b.configManager.UpdateUTXOConfig(utxo)
}

// UpdateTransactionConfig updates the transaction configuration
func (b *BSV) UpdateTransactionConfig(tx *config.TransactionConfig) error {
	return b.configManager.UpdateTransactionConfig(tx)
}

// SetNetworkType sets the network type with predefined configurations
func (b *BSV) SetNetworkType(networkType config.NetworkType) error {
	err := b.configManager.SetNetworkType(networkType)
	if err != nil {
		return err
	}

	// Update wallet generator with new network
	networkConfig := b.configManager.GetNetworkConfig()
	b.walletGen = wallet.NewGenerator(networkConfig.IsTestnet)

	return nil
}

// ClearUTXOCache clears the UTXO cache
func (b *BSV) ClearUTXOCache() {
	b.txBuilder.ClearUTXOCache()
}

// ClearUTXOCacheForAddress clears UTXO cache for a specific address
func (b *BSV) ClearUTXOCacheForAddress(address string) {
	b.txBuilder.ClearUTXOCacheForAddress(address)
}

// Package-level enhanced functions for convenience

// NewBSVDefault creates a new BSV instance with default configuration
func NewBSVDefault() *BSV {
	configManager := config.NewManager()
	return NewBSV(configManager)
}

// GenerateWallet creates a BSV wallet from a mnemonic with enhanced support
func GenerateWalletEnhanced(mnemonicPhrase string, networkType config.NetworkType) (*types.WalletResult, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.GenerateWallet(mnemonicPhrase)
}

// GenerateWalletWithKeypair creates a wallet with keypair
func GenerateWalletWithKeypairEnhanced(mnemonicPhrase string, networkType config.NetworkType) (*types.WalletResult, *wallet.KeyPair, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, nil, err
	}
	return bsv.GenerateWalletWithKeypair(mnemonicPhrase)
}

// ValidateAddress validates a BSV address
func ValidateAddressEnhanced(address string, networkType config.NetworkType) error {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return err
	}
	return bsv.ValidateAddress(address)
}

// GetEnhancedBalance retrieves enhanced balance for an address
func GetEnhancedBalance(address string, networkType config.NetworkType) (*types.EnhancedBalanceInfo, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.GetEnhancedBalance(address)
}

// GetNativeBalance retrieves native balance for an address
func GetNativeBalance(address string, networkType config.NetworkType) (*types.NativeBalanceInfo, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.GetNativeBalance(address)
}

// GetNonNativeBalance retrieves non-native balance for an address
func GetNonNativeBalance(address string, networkType config.NetworkType) (*types.NonNativeBalanceInfo, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.GetNonNativeBalance(address)
}

// GetBalance retrieves balance for an address (backward compatibility)
func GetBalanceEnhanced(address string, networkType config.NetworkType) (int64, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return 0, err
	}
	return bsv.GetBalance(address)
}

// GetUTXOs retrieves UTXOs for an address
func GetUTXOsEnhanced(address string, networkType config.NetworkType) ([]types.UTXO, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.GetUTXOs(address)
}

// SignAndSendTransaction creates and sends a transaction with enhanced support
func SignAndSendTransactionEnhanced(params *types.TransactionParams, networkType config.NetworkType) (*types.TransactionResult, error) {
	bsv, err := NewBSVWithNetwork(networkType)
	if err != nil {
		return nil, err
	}
	return bsv.SignAndSendTransaction(params)
}
