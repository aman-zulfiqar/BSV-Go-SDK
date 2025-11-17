package config

import (
	"fmt"
	"sync"
)

// NetworkType represents the type of network
type NetworkType string

const (
	Mainnet NetworkType = "mainnet"
	Testnet NetworkType = "testnet"
	Custom  NetworkType = "custom"
)

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Name        string `json:"name"`        // Network name
	RPCURL      string `json:"rpcUrl"`      // RPC endpoint URL
	ExplorerURL string `json:"explorerUrl"` // Explorer URL
	IsTestnet   bool   `json:"isTestnet"`   // Whether this is testnet
	ChainID     string `json:"chainId"`     // Chain identifier
	CoinType    uint32 `json:"coinType"`    // BIP44 coin type
}

// UTXOConfig represents UTXO handling configuration
type UTXOConfig struct {
	IncludeNative    bool `json:"includeNative"`    // Include native BSV UTXOs
	IncludeNonNative bool `json:"includeNonNative"` // Include non-native token UTXOs
	MinConfirmations int  `json:"minConfirmations"` // Minimum confirmations required
	MaxUTXOsPerQuery int  `json:"maxUTXOsPerQuery"` // Maximum UTXOs per query
	EnableCaching    bool `json:"enableCaching"`    // Enable UTXO caching
	CacheExpiry      int  `json:"cacheExpiry"`      // Cache expiry in seconds
}

// TransactionConfig represents transaction configuration
type TransactionConfig struct {
	DefaultFeeRate        int64 `json:"defaultFeeRate"`        // Default fee rate in sat/vbyte
	MinFeeRate            int64 `json:"minFeeRate"`            // Minimum fee rate
	MaxFeeRate            int64 `json:"maxFeeRate"`            // Maximum fee rate
	DustLimit             int64 `json:"dustLimit"`             // Dust limit in satoshis
	MaxTransactionSize    int   `json:"maxTransactionSize"`    // Maximum transaction size in bytes
	EnableRBF             bool  `json:"enableRBF"`             // Enable Replace-By-Fee
	IncludeNativeUTXOs    bool  `json:"includeNativeUTXOs"`    // Include native BSV UTXOs in transactions
	IncludeNonNativeUTXOs bool  `json:"includeNonNativeUTXOs"` // Include non-native token UTXOs in transactions
}

// Config represents the complete configuration
type Config struct {
	Network     *NetworkConfig     `json:"network"`
	UTXO        *UTXOConfig        `json:"utxo"`
	Transaction *TransactionConfig `json:"transaction"`
}

// Manager handles dynamic configuration
type Manager struct {
	config *Config
	mutex  sync.RWMutex
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		config: getDefaultConfig(),
	}
}

// NewManagerWithConfig creates a new configuration manager with custom config
func NewManagerWithConfig(config *Config) *Manager {
	return &Manager{
		config: config,
	}
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a deep copy to prevent external modifications
	return m.deepCopyConfig()
}

// UpdateConfig updates the configuration
func (m *Manager) UpdateConfig(config *Config) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	m.config = m.deepCopyConfigFrom(config)
	return nil
}

// UpdateNetworkConfig updates only the network configuration
func (m *Manager) UpdateNetworkConfig(network *NetworkConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateNetworkConfig(network); err != nil {
		return fmt.Errorf("invalid network configuration: %v", err)
	}

	m.config.Network = m.deepCopyNetworkConfig(network)
	return nil
}

// UpdateUTXOConfig updates only the UTXO configuration
func (m *Manager) UpdateUTXOConfig(utxo *UTXOConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateUTXOConfig(utxo); err != nil {
		return fmt.Errorf("invalid UTXO configuration: %v", err)
	}

	m.config.UTXO = m.deepCopyUTXOConfig(utxo)
	return nil
}

// UpdateTransactionConfig updates only the transaction configuration
func (m *Manager) UpdateTransactionConfig(tx *TransactionConfig) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.validateTransactionConfig(tx); err != nil {
		return fmt.Errorf("invalid transaction configuration: %v", err)
	}

	m.config.Transaction = m.deepCopyTransactionConfig(tx)
	return nil
}

// GetNetworkConfig returns the current network configuration
func (m *Manager) GetNetworkConfig() *NetworkConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.deepCopyNetworkConfig(m.config.Network)
}

// GetUTXOConfig returns the current UTXO configuration
func (m *Manager) GetUTXOConfig() *UTXOConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.deepCopyUTXOConfig(m.config.UTXO)
}

// GetTransactionConfig returns the current transaction configuration
func (m *Manager) GetTransactionConfig() *TransactionConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.deepCopyTransactionConfig(m.config.Transaction)
}

// SetNetworkType sets the network type with predefined configurations
func (m *Manager) SetNetworkType(networkType NetworkType) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch networkType {
	case Mainnet:
		m.config.Network = getMainnetConfig()
	case Testnet:
		m.config.Network = getTestnetConfig()
	default:
		return fmt.Errorf("unsupported network type: %s", networkType)
	}

	return nil
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return getDefaultConfig()
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *Config {
	return &Config{
		Network:     getTestnetConfig(),
		UTXO:        getDefaultUTXOConfig(),
		Transaction: getDefaultTransactionConfig(),
	}
}

// getMainnetConfig returns mainnet configuration
func getMainnetConfig() *NetworkConfig {
	return &NetworkConfig{
		Name:        "BSV Mainnet",
		RPCURL:      "https://api.whatsonchain.com/v1/bsv/main",
		ExplorerURL: "https://whatsonchain.com",
		IsTestnet:   false,
		ChainID:     "mainnet",
		CoinType:    236, // BSV mainnet coin type
	}
}

// getTestnetConfig returns testnet configuration
func getTestnetConfig() *NetworkConfig {
	return &NetworkConfig{
		Name:        "BSV Testnet",
		RPCURL:      "https://api.whatsonchain.com/v1/bsv/test",
		ExplorerURL: "https://test.whatsonchain.com",
		IsTestnet:   true,
		ChainID:     "testnet",
		CoinType:    1, // BSV testnet coin type
	}
}

// getDefaultUTXOConfig returns default UTXO configuration
func getDefaultUTXOConfig() *UTXOConfig {
	return &UTXOConfig{
		IncludeNative:    true,
		IncludeNonNative: true,
		MinConfirmations: 1,
		MaxUTXOsPerQuery: 100,
		EnableCaching:    true,
		CacheExpiry:      300, // 5 minutes
	}
}

// getDefaultTransactionConfig returns default transaction configuration
func getDefaultTransactionConfig() *TransactionConfig {
	return &TransactionConfig{
		DefaultFeeRate:        5,
		MinFeeRate:            1,
		MaxFeeRate:            1000,
		DustLimit:             546,
		MaxTransactionSize:    100000, // 100KB
		EnableRBF:             false,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
	}
}

// Validation methods
func (m *Manager) validateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	if err := m.validateNetworkConfig(config.Network); err != nil {
		return err
	}

	if err := m.validateUTXOConfig(config.UTXO); err != nil {
		return err
	}

	if err := m.validateTransactionConfig(config.Transaction); err != nil {
		return err
	}

	return nil
}

func (m *Manager) validateNetworkConfig(network *NetworkConfig) error {
	if network == nil {
		return fmt.Errorf("network configuration cannot be nil")
	}

	if network.Name == "" {
		return fmt.Errorf("network name is required")
	}

	if network.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}

	if network.ExplorerURL == "" {
		return fmt.Errorf("explorer URL is required")
	}

	if network.ChainID == "" {
		return fmt.Errorf("chain ID is required")
	}

	return nil
}

func (m *Manager) validateUTXOConfig(utxo *UTXOConfig) error {
	if utxo == nil {
		return fmt.Errorf("UTXO configuration cannot be nil")
	}

	if utxo.MinConfirmations < 0 {
		return fmt.Errorf("minimum confirmations cannot be negative")
	}

	if utxo.MaxUTXOsPerQuery <= 0 {
		return fmt.Errorf("maximum UTXOs per query must be positive")
	}

	if utxo.CacheExpiry < 0 {
		return fmt.Errorf("cache expiry cannot be negative")
	}

	return nil
}

func (m *Manager) validateTransactionConfig(tx *TransactionConfig) error {
	if tx == nil {
		return fmt.Errorf("transaction configuration cannot be nil")
	}

	if tx.DefaultFeeRate <= 0 {
		return fmt.Errorf("default fee rate must be positive")
	}

	if tx.MinFeeRate <= 0 {
		return fmt.Errorf("minimum fee rate must be positive")
	}

	if tx.MaxFeeRate <= 0 {
		return fmt.Errorf("maximum fee rate must be positive")
	}

	if tx.MinFeeRate > tx.MaxFeeRate {
		return fmt.Errorf("minimum fee rate cannot be greater than maximum fee rate")
	}

	if tx.DustLimit < 0 {
		return fmt.Errorf("dust limit cannot be negative")
	}

	if tx.MaxTransactionSize <= 0 {
		return fmt.Errorf("maximum transaction size must be positive")
	}

	return nil
}

// Deep copy methods
func (m *Manager) deepCopyConfig() *Config {
	return m.deepCopyConfigFrom(m.config)
}

func (m *Manager) deepCopyConfigFrom(config *Config) *Config {
	return &Config{
		Network:     m.deepCopyNetworkConfig(config.Network),
		UTXO:        m.deepCopyUTXOConfig(config.UTXO),
		Transaction: m.deepCopyTransactionConfig(config.Transaction),
	}
}

func (m *Manager) deepCopyNetworkConfig(network *NetworkConfig) *NetworkConfig {
	if network == nil {
		return nil
	}
	return &NetworkConfig{
		Name:        network.Name,
		RPCURL:      network.RPCURL,
		ExplorerURL: network.ExplorerURL,
		IsTestnet:   network.IsTestnet,
		ChainID:     network.ChainID,
		CoinType:    network.CoinType,
	}
}

func (m *Manager) deepCopyUTXOConfig(utxo *UTXOConfig) *UTXOConfig {
	if utxo == nil {
		return nil
	}
	return &UTXOConfig{
		IncludeNative:    utxo.IncludeNative,
		IncludeNonNative: utxo.IncludeNonNative,
		MinConfirmations: utxo.MinConfirmations,
		MaxUTXOsPerQuery: utxo.MaxUTXOsPerQuery,
		EnableCaching:    utxo.EnableCaching,
		CacheExpiry:      utxo.CacheExpiry,
	}
}

func (m *Manager) deepCopyTransactionConfig(tx *TransactionConfig) *TransactionConfig {
	if tx == nil {
		return nil
	}
	return &TransactionConfig{
		DefaultFeeRate:        tx.DefaultFeeRate,
		MinFeeRate:            tx.MinFeeRate,
		MaxFeeRate:            tx.MaxFeeRate,
		DustLimit:             tx.DustLimit,
		MaxTransactionSize:    tx.MaxTransactionSize,
		EnableRBF:             tx.EnableRBF,
		IncludeNativeUTXOs:    tx.IncludeNativeUTXOs,
		IncludeNonNativeUTXOs: tx.IncludeNonNativeUTXOs,
	}
}
