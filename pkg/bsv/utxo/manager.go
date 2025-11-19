package utxo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

// Manager handles UTXO management with dynamic configuration
type Manager struct {
	configManager *config.Manager
	httpClient    *http.Client
	maxRetries    int
	retryDelay    time.Duration
	cache         map[string]*CacheEntry
	cacheMutex    sync.RWMutex
}

// CacheEntry represents a cached UTXO entry
type CacheEntry struct {
	UTXOs     []types.UTXO
	Balance   *types.EnhancedBalanceInfo
	Timestamp time.Time
}

// EnhancedUTXOResponse represents the response from What's On Chain API
type EnhancedUTXOResponse struct {
	TxID          string `json:"txid"`
	Vout          uint32 `json:"vout"`
	Value         int64  `json:"value"`
	ScriptPubKey  string `json:"scriptPubKey"`
	Address       string `json:"address"`
	Confirmations int    `json:"confirmations"`
	Height        int    `json:"height"`
}

// EnhancedBalanceResponse represents balance response from API
type EnhancedBalanceResponse struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
}

// NewManager creates a new UTXO manager
func NewManager(configManager *config.Manager) *Manager {
	return &Manager{
		configManager: configManager,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		cache:      make(map[string]*CacheEntry),
	}
}

// GetUTXOs retrieves UTXOs for a given address with dynamic configuration
func (m *Manager) GetUTXOs(address string) ([]types.UTXO, error) {
	// Check cache first
	if cached := m.getFromCache(address); cached != nil {
		return cached.UTXOs, nil
	}

	networkConfig := m.configManager.GetNetworkConfig()
	utxoConfig := m.configManager.GetUTXOConfig()

	url := fmt.Sprintf("%s/address/%s/unspent", networkConfig.RPCURL, address)

	var utxoResponses []EnhancedUTXOResponse
	err := m.makeRequest(url, &utxoResponses)
	if err != nil {
		return nil, fmt.Errorf("failed to get UTXOs: %v", err)
	}

	// Convert to enhanced UTXO format
	var utxos []types.UTXO
	for _, resp := range utxoResponses {
		// Check minimum confirmations
		if resp.Confirmations < utxoConfig.MinConfirmations {
			continue
		}

		// Check max UTXOs per query
		if len(utxos) >= utxoConfig.MaxUTXOsPerQuery {
			break
		}

		utxo := types.UTXO{
			TxID:          resp.TxID,
			Vout:          resp.Vout,
			Value:         resp.Value,
			ScriptPubKey:  resp.ScriptPubKey,
			Address:       resp.Address,
			Confirmations: resp.Confirmations,
			Height:        resp.Height,
			IsNative:      true, // All UTXOs from this API are native BSV
			TokenID:       "",   // Empty for native UTXOs
			TokenAmount:   0,    // Zero for native UTXOs
		}

		utxos = append(utxos, utxo)
	}

	// Cache the results if caching is enabled
	if utxoConfig.EnableCaching {
		m.setCache(address, &CacheEntry{
			UTXOs:     utxos,
			Timestamp: time.Now(),
		})
	}

	return utxos, nil
}

// GetEnhancedBalance retrieves enhanced balance information for an address
func (m *Manager) GetEnhancedBalance(address string) (*types.EnhancedBalanceInfo, error) {
	// Check cache first
	if cached := m.getFromCache(address); cached != nil && cached.Balance != nil {
		return cached.Balance, nil
	}

	networkConfig := m.configManager.GetNetworkConfig()
	utxoConfig := m.configManager.GetUTXOConfig()

	// Get native balance from API
	balanceURL := fmt.Sprintf("%s/address/%s/balance", networkConfig.RPCURL, address)
	var balanceResp EnhancedBalanceResponse
	err := m.makeRequest(balanceURL, &balanceResp)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	// Get UTXOs for detailed analysis
	utxos, err := m.GetUTXOs(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get UTXOs for balance calculation: %v", err)
	}

	// Separate native and non-native UTXOs
	var nativeUTXOs []types.UTXO
	var nonNativeUTXOs []types.UTXO

	for _, utxo := range utxos {
		if utxo.IsNative {
			nativeUTXOs = append(nativeUTXOs, utxo)
		} else {
			nonNativeUTXOs = append(nonNativeUTXOs, utxo)
		}
	}

	// Calculate native balance
	nativeBalance := &types.NativeBalanceInfo{
		Confirmed:   balanceResp.Confirmed,
		Unconfirmed: balanceResp.Unconfirmed,
		Total:       balanceResp.Confirmed + balanceResp.Unconfirmed,
		UTXOCount:   len(nativeUTXOs),
	}

	// Calculate non-native balances
	nonNativeBalance := &types.NonNativeBalanceInfo{
		Tokens:    make(map[string]*types.TokenBalance),
		UTXOCount: len(nonNativeUTXOs),
	}

	// Group non-native UTXOs by token ID
	for _, utxo := range nonNativeUTXOs {
		if utxo.TokenID != "" {
			if _, exists := nonNativeBalance.Tokens[utxo.TokenID]; !exists {
				nonNativeBalance.Tokens[utxo.TokenID] = &types.TokenBalance{
					TokenID:     utxo.TokenID,
					Confirmed:   0,
					Unconfirmed: 0,
					Total:       0,
					UTXOCount:   0,
				}
			}

			tokenBalance := nonNativeBalance.Tokens[utxo.TokenID]
			tokenBalance.UTXOCount++

			if utxo.Confirmations >= utxoConfig.MinConfirmations {
				tokenBalance.Confirmed += utxo.TokenAmount
			} else {
				tokenBalance.Unconfirmed += utxo.TokenAmount
			}
			tokenBalance.Total += utxo.TokenAmount
		}
	}

	enhancedBalance := &types.EnhancedBalanceInfo{
		Native:    nativeBalance,
		NonNative: nonNativeBalance,
		Total:     nativeBalance.Total,
	}

	// Cache the results if caching is enabled
	if utxoConfig.EnableCaching {
		if cached := m.getFromCache(address); cached != nil {
			cached.Balance = enhancedBalance
		} else {
			m.setCache(address, &CacheEntry{
				Balance:   enhancedBalance,
				Timestamp: time.Now(),
			})
		}
	}

	return enhancedBalance, nil
}

// GetNativeBalance returns only native BSV balance
func (m *Manager) GetNativeBalance(address string) (*types.NativeBalanceInfo, error) {
	enhancedBalance, err := m.GetEnhancedBalance(address)
	if err != nil {
		return nil, err
	}
	return enhancedBalance.Native, nil
}

// GetNonNativeBalance returns only non-native token balances
func (m *Manager) GetNonNativeBalance(address string) (*types.NonNativeBalanceInfo, error) {
	enhancedBalance, err := m.GetEnhancedBalance(address)
	if err != nil {
		return nil, err
	}
	return enhancedBalance.NonNative, nil
}

// GetConfirmedBalance returns only confirmed native balance (for backward compatibility)
func (m *Manager) GetConfirmedBalance(address string) (int64, error) {
	nativeBalance, err := m.GetNativeBalance(address)
	if err != nil {
		return 0, err
	}
	return nativeBalance.Confirmed, nil
}

// SelectUTXOs selects UTXOs for a transaction with enhanced filtering
func (m *Manager) SelectUTXOs(address string, amount, feeRate int64) ([]types.UTXO, int64, error) {
	txConfig := m.configManager.GetTransactionConfig()

	// Get all UTXOs
	allUTXOs, err := m.GetUTXOs(address)
	if err != nil {
		return nil, 0, err
	}

	if len(allUTXOs) == 0 {
		return nil, 0, fmt.Errorf("no UTXOs available for address: %s", address)
	}

	// Filter UTXOs based on configuration
	var availableUTXOs []types.UTXO
	for _, utxo := range allUTXOs {
		// Check if we should include native UTXOs
		if utxo.IsNative && txConfig.IncludeNativeUTXOs {
			availableUTXOs = append(availableUTXOs, utxo)
		}
		// Check if we should include non-native UTXOs
		if !utxo.IsNative && txConfig.IncludeNonNativeUTXOs {
			availableUTXOs = append(availableUTXOs, utxo)
		}
	}

	if len(availableUTXOs) == 0 {
		return nil, 0, fmt.Errorf("no suitable UTXOs available based on configuration")
	}

	// Sort UTXOs by value (largest first for efficiency)
	sortedUTXOs := m.sortUTXOsByValue(availableUTXOs)

	var selectedUTXOs []types.UTXO
	var totalValue int64
	var estimatedFee int64

	// Estimate transaction size (simplified)
	// Input: ~148 bytes, Output: ~34 bytes, Change: ~34 bytes
	estimatedSize := 10 + len(sortedUTXOs)*148 + 34 + 34 // Base size + inputs + outputs
	estimatedFee = int64(estimatedSize) * feeRate

	// Validate fee rate
	if feeRate < txConfig.MinFeeRate {
		feeRate = txConfig.DefaultFeeRate
	}
	if feeRate > txConfig.MaxFeeRate {
		return nil, 0, fmt.Errorf("fee rate %d exceeds maximum allowed %d", feeRate, txConfig.MaxFeeRate)
	}

	// Select UTXOs until we have enough funds
	for _, utxo := range sortedUTXOs {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalValue += utxo.Value

		// Recalculate fee with current number of inputs
		currentSize := 10 + len(selectedUTXOs)*148 + 34 + 34
		currentFee := int64(currentSize) * feeRate

		if totalValue >= amount+currentFee {
			return selectedUTXOs, currentFee, nil
		}
	}

	return nil, 0, fmt.Errorf("insufficient funds: need %d satoshis, have %d satoshis",
		amount+estimatedFee, totalValue)
}

// SelectUTXOsForTokenTransfer selects UTXOs for token transfers
func (m *Manager) SelectUTXOsForTokenTransfer(address string, tokenID string, amount int64, feeRate int64) ([]types.UTXO, int64, error) {
	// Get all UTXOs
	allUTXOs, err := m.GetUTXOs(address)
	if err != nil {
		return nil, 0, err
	}

	// Filter for specific token UTXOs
	var tokenUTXOs []types.UTXO
	var nativeUTXOs []types.UTXO

	for _, utxo := range allUTXOs {
		if utxo.IsNative {
			nativeUTXOs = append(nativeUTXOs, utxo)
		} else if utxo.TokenID == tokenID {
			tokenUTXOs = append(tokenUTXOs, utxo)
		}
	}

	// Check if we have enough token UTXOs
	var totalTokenAmount int64
	for _, utxo := range tokenUTXOs {
		totalTokenAmount += utxo.TokenAmount
	}

	if totalTokenAmount < amount {
		return nil, 0, fmt.Errorf("insufficient token balance: need %d, have %d", amount, totalTokenAmount)
	}

	// Select token UTXOs for the transfer
	var selectedTokenUTXOs []types.UTXO
	var selectedTokenAmount int64

	// Sort token UTXOs by token amount
	sortedTokenUTXOs := m.sortUTXOsByTokenAmount(tokenUTXOs)

	for _, utxo := range sortedTokenUTXOs {
		selectedTokenUTXOs = append(selectedTokenUTXOs, utxo)
		selectedTokenAmount += utxo.TokenAmount

		if selectedTokenAmount >= amount {
			break
		}
	}

	// We also need native UTXOs for fees
	// Estimate fee for transaction with token UTXOs
	estimatedSize := 10 + len(selectedTokenUTXOs)*148 + 34 + 34 // Base + token inputs + outputs
	estimatedFee := int64(estimatedSize) * feeRate

	// Select native UTXOs for fees
	var selectedNativeUTXOs []types.UTXO
	var totalNativeValue int64

	sortedNativeUTXOs := m.sortUTXOsByValue(nativeUTXOs)

	for _, utxo := range sortedNativeUTXOs {
		selectedNativeUTXOs = append(selectedNativeUTXOs, utxo)
		totalNativeValue += utxo.Value

		if totalNativeValue >= estimatedFee {
			break
		}
	}

	if totalNativeValue < estimatedFee {
		return nil, 0, fmt.Errorf("insufficient native balance for fees: need %d, have %d", estimatedFee, totalNativeValue)
	}

	// Combine all selected UTXOs
	var allSelectedUTXOs []types.UTXO
	allSelectedUTXOs = append(allSelectedUTXOs, selectedTokenUTXOs...)
	allSelectedUTXOs = append(allSelectedUTXOs, selectedNativeUTXOs...)

	return allSelectedUTXOs, estimatedFee, nil
}

// CalculateChange calculates the change amount after selecting UTXOs
func (m *Manager) CalculateChange(selectedUTXOs []types.UTXO, amount, fee int64) (int64, bool) {
	txConfig := m.configManager.GetTransactionConfig()

	var totalInput int64
	for _, utxo := range selectedUTXOs {
		totalInput += utxo.Value
	}

	change := totalInput - amount - fee

	// Only return change if it's above dust limit
	if change > txConfig.DustLimit {
		return change, true
	}

	return 0, false
}

// ClearCache clears the UTXO cache
func (m *Manager) ClearCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.cache = make(map[string]*CacheEntry)
}

// ClearCacheForAddress clears cache for a specific address
func (m *Manager) ClearCacheForAddress(address string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	delete(m.cache, address)
}

// Helper methods

func (m *Manager) getFromCache(address string) *CacheEntry {
	utxoConfig := m.configManager.GetUTXOConfig()
	if !utxoConfig.EnableCaching {
		return nil
	}

	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	entry, exists := m.cache[address]
	if !exists {
		return nil
	}

	// Check if cache entry has expired
	if time.Since(entry.Timestamp) > time.Duration(utxoConfig.CacheExpiry)*time.Second {
		return nil
	}

	return entry
}

func (m *Manager) setCache(address string, entry *CacheEntry) {
	utxoConfig := m.configManager.GetUTXOConfig()
	if !utxoConfig.EnableCaching {
		return
	}

	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	m.cache[address] = entry
}

func (m *Manager) sortUTXOsByValue(utxos []types.UTXO) []types.UTXO {
	sorted := make([]types.UTXO, len(utxos))
	copy(sorted, utxos)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Value < sorted[j+1].Value {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func (m *Manager) sortUTXOsByTokenAmount(utxos []types.UTXO) []types.UTXO {
	sorted := make([]types.UTXO, len(utxos))
	copy(sorted, utxos)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].TokenAmount < sorted[j+1].TokenAmount {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func (m *Manager) makeRequest(url string, result interface{}) error {
	var lastErr error

	for attempt := 1; attempt <= m.maxRetries; attempt++ {
		resp, err := m.httpClient.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %v", err)
			if attempt < m.maxRetries {
				time.Sleep(m.retryDelay * time.Duration(attempt))
				continue
			}
			break
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
			if attempt < m.maxRetries {
				time.Sleep(m.retryDelay * time.Duration(attempt))
				continue
			}
			break
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %v", err)
			if attempt < m.maxRetries {
				time.Sleep(m.retryDelay * time.Duration(attempt))
				continue
			}
			break
		}

		err = json.Unmarshal(body, result)
		if err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %v", err)
			if attempt < m.maxRetries {
				time.Sleep(m.retryDelay * time.Duration(attempt))
				continue
			}
			break
		}

		return nil // Success
	}

	return lastErr
}

// SetRetryConfig sets the retry configuration
func (m *Manager) SetRetryConfig(maxRetries int, retryDelay time.Duration) {
	m.maxRetries = maxRetries
	m.retryDelay = retryDelay
}
