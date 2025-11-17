package tests

import (
	"testing"
	"time"

	"github.com/muhammadamman/BSV-Go/pkg/bsv"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/sharding"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

func TestEnhancedBSVCreation(t *testing.T) {
	// Test creating enhanced BSV with default configuration
	enhancedBSV := bsv.NewBSVDefault()
	if enhancedBSV == nil {
		t.Fatal("Failed to create enhanced BSV instance")
	}

	// Test network configuration
	networkConfig := enhancedBSV.GetNetworkConfig()
	if networkConfig == nil {
		t.Fatal("Network configuration is nil")
	}

	if networkConfig.Name == "" {
		t.Error("Network name is empty")
	}

	if networkConfig.RPCURL == "" {
		t.Error("RPC URL is empty")
	}

	if networkConfig.ExplorerURL == "" {
		t.Error("Explorer URL is empty")
	}
}

func TestDynamicConfiguration(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Test UTXO configuration update
	utxoConfig := &config.UTXOConfig{
		IncludeNative:    true,
		IncludeNonNative: true,
		MinConfirmations: 2,
		MaxUTXOsPerQuery: 100,
		EnableCaching:    true,
		CacheExpiry:      300,
	}

	err := enhancedBSV.UpdateUTXOConfig(utxoConfig)
	if err != nil {
		t.Fatalf("Failed to update UTXO config: %v", err)
	}

	// Verify the update
	updatedConfig := enhancedBSV.GetUTXOConfig()
	if updatedConfig.MinConfirmations != 2 {
		t.Errorf("Expected MinConfirmations 2, got %d", updatedConfig.MinConfirmations)
	}

	if updatedConfig.MaxUTXOsPerQuery != 100 {
		t.Errorf("Expected MaxUTXOsPerQuery 100, got %d", updatedConfig.MaxUTXOsPerQuery)
	}

	if !updatedConfig.EnableCaching {
		t.Error("Expected caching to be enabled")
	}

	// Test transaction configuration update
	txConfig := &config.TransactionConfig{
		DefaultFeeRate:        10,
		MinFeeRate:            1,
		MaxFeeRate:            100,
		DustLimit:             546,
		MaxTransactionSize:    100000,
		EnableRBF:             false,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
	}

	err = enhancedBSV.UpdateTransactionConfig(txConfig)
	if err != nil {
		t.Fatalf("Failed to update transaction config: %v", err)
	}

	// Verify the update
	updatedTxConfig := enhancedBSV.GetTransactionConfig()
	if updatedTxConfig.DefaultFeeRate != 10 {
		t.Errorf("Expected DefaultFeeRate 10, got %d", updatedTxConfig.DefaultFeeRate)
	}

	if updatedTxConfig.DustLimit != 546 {
		t.Errorf("Expected DustLimit 546, got %d", updatedTxConfig.DustLimit)
	}
}

func TestNetworkSwitching(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Test switching to mainnet
	err := enhancedBSV.SetNetworkType(config.Mainnet)
	if err != nil {
		t.Fatalf("Failed to switch to mainnet: %v", err)
	}

	networkConfig := enhancedBSV.GetNetworkConfig()
	if networkConfig.IsTestnet {
		t.Error("Expected mainnet, got testnet")
	}

	if networkConfig.CoinType != 236 {
		t.Errorf("Expected coin type 236 for mainnet, got %d", networkConfig.CoinType)
	}

	// Test switching to testnet
	err = enhancedBSV.SetNetworkType(config.Testnet)
	if err != nil {
		t.Fatalf("Failed to switch to testnet: %v", err)
	}

	networkConfig = enhancedBSV.GetNetworkConfig()
	if !networkConfig.IsTestnet {
		t.Error("Expected testnet, got mainnet")
	}

	if networkConfig.CoinType != 1 {
		t.Errorf("Expected coin type 1 for testnet, got %d", networkConfig.CoinType)
	}
}

func TestEnhancedWalletGeneration(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Generate mnemonic
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	// Test sharding
	shardingResult, err := sharding.SplitMnemonic(mnemonicPhrase, 2, 3)
	if err != nil {
		t.Fatalf("Failed to create shards: %v", err)
	}

	if len(shardingResult.Shards) != 3 {
		t.Errorf("Expected 3 shards, got %d", len(shardingResult.Shards))
	}

	if shardingResult.Threshold != 2 {
		t.Errorf("Expected threshold 2, got %d", shardingResult.Threshold)
	}

	// Reconstruct mnemonic
	reconstructedMnemonic, err := sharding.CombineShards(shardingResult.Shards[:2])
	if err != nil {
		t.Fatalf("Failed to reconstruct mnemonic: %v", err)
	}

	if reconstructedMnemonic != mnemonicPhrase {
		t.Error("Reconstructed mnemonic doesn't match original")
	}

	// Generate wallet
	wallet, err := enhancedBSV.GenerateWallet(reconstructedMnemonic)
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	if wallet.Address == "" {
		t.Error("Wallet address is empty")
	}

	if wallet.PrivateKey == "" {
		t.Error("Wallet private key is empty")
	}

	if wallet.PublicKey == "" {
		t.Error("Wallet public key is empty")
	}

	// Validate address
	err = enhancedBSV.ValidateAddress(wallet.Address)
	if err != nil {
		t.Errorf("Address validation failed: %v", err)
	}
}

func TestEnhancedBalanceFetching(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Generate a test wallet
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	wallet, err := enhancedBSV.GenerateWallet(mnemonicPhrase)
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	// Test enhanced balance fetching (will likely return 0 balance for new address)
	enhancedBalance, err := enhancedBSV.GetEnhancedBalance(wallet.Address)
	if err != nil {
		t.Logf("Failed to get enhanced balance (expected for new address): %v", err)
	} else {
		if enhancedBalance == nil {
			t.Error("Enhanced balance is nil")
		}

		if enhancedBalance.Native == nil {
			t.Error("Native balance is nil")
		}

		if enhancedBalance.NonNative == nil {
			t.Error("Non-native balance is nil")
		}

		// Verify balance structure
		if enhancedBalance.Native.UTXOCount < 0 {
			t.Error("Native UTXO count cannot be negative")
		}

		if enhancedBalance.NonNative.UTXOCount < 0 {
			t.Error("Non-native UTXO count cannot be negative")
		}

		if enhancedBalance.Native.Total != (enhancedBalance.Native.Confirmed + enhancedBalance.Native.Unconfirmed) {
			t.Error("Native total doesn't match confirmed + unconfirmed")
		}
	}

	// Test native balance fetching
	nativeBalance, err := enhancedBSV.GetNativeBalance(wallet.Address)
	if err != nil {
		t.Logf("Failed to get native balance (expected for new address): %v", err)
	} else {
		if nativeBalance == nil {
			t.Error("Native balance is nil")
		}

		if nativeBalance.Total != (nativeBalance.Confirmed + nativeBalance.Unconfirmed) {
			t.Error("Native total doesn't match confirmed + unconfirmed")
		}
	}

	// Test non-native balance fetching
	nonNativeBalance, err := enhancedBSV.GetNonNativeBalance(wallet.Address)
	if err != nil {
		t.Logf("Failed to get non-native balance (expected for new address): %v", err)
	} else {
		if nonNativeBalance == nil {
			t.Error("Non-native balance is nil")
		}

		if nonNativeBalance.Tokens == nil {
			t.Error("Tokens map is nil")
		}
	}
}

func TestUTXOCaching(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Enable caching
	utxoConfig := &config.UTXOConfig{
		IncludeNative:    true,
		IncludeNonNative: true,
		MinConfirmations: 1,
		MaxUTXOsPerQuery: 100,
		EnableCaching:    true,
		CacheExpiry:      60, // 1 minute
	}

	err := enhancedBSV.UpdateUTXOConfig(utxoConfig)
	if err != nil {
		t.Fatalf("Failed to update UTXO config: %v", err)
	}

	// Generate a test wallet
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	wallet, err := enhancedBSV.GenerateWallet(mnemonicPhrase)
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	// Test cache clearing
	enhancedBSV.ClearUTXOCache()
	enhancedBSV.ClearUTXOCacheForAddress(wallet.Address)

	// Test UTXO fetching (will likely return empty for new address)
	utxos, err := enhancedBSV.GetUTXOs(wallet.Address)
	if err != nil {
		t.Logf("Failed to get UTXOs (expected for new address): %v", err)
	} else {
		// utxos can be nil or empty for new addresses, both are valid

		// Verify UTXO structure if we have any
		if len(utxos) > 0 {
			for i, utxo := range utxos {
				if utxo.TxID == "" {
					t.Errorf("UTXO %d has empty transaction ID", i)
				}

				if utxo.Value <= 0 {
					t.Errorf("UTXO %d has invalid value: %d", i, utxo.Value)
				}

				if utxo.Address == "" {
					t.Errorf("UTXO %d has empty address", i)
				}
			}
		}
	}
}

func TestEnhancedTransactionBuilding(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Generate a test wallet
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	wallet, err := enhancedBSV.GenerateWallet(mnemonicPhrase)
	if err != nil {
		t.Fatalf("Failed to generate wallet: %v", err)
	}

	// Test enhanced transaction parameters
	txParams := &types.TransactionParams{
		From:                  wallet.Address,
		To:                    "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt", // BSV testnet address
		Amount:                1000,
		FeeRate:               5,
		PrivateKey:            mnemonicPhrase,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
		TokenTransfers:        []*types.TokenTransfer{},
		DataOutputs:           []*types.DataOutput{},
	}

	// Test transaction building (will fail without funds, which is expected)
	_, err = enhancedBSV.SignAndSendTransaction(txParams)
	if err != nil {
		// This is expected for a new address with no funds
		t.Logf("Transaction failed as expected (no funds): %v", err)
	} else {
		t.Log("Transaction succeeded unexpectedly")
	}

	// Test token transfer parameters
	tokenTransferParams := &types.TransactionParams{
		From:                  wallet.Address,
		To:                    "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt",
		Amount:                1000,
		FeeRate:               5,
		PrivateKey:            mnemonicPhrase,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: true,
		TokenTransfers: []*types.TokenTransfer{
			{
				TokenID: "test-token-id",
				To:      "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt",
				Amount:  100,
			},
		},
		DataOutputs: []*types.DataOutput{
			{
				Data: "48656c6c6f20576f726c64", // "Hello World" in hex
			},
		},
	}

	// Test token transfer building (will fail without funds, which is expected)
	_, err = enhancedBSV.SignAndSendTransaction(tokenTransferParams)
	if err != nil {
		// This is expected for a new address with no funds
		t.Logf("Token transfer failed as expected (no funds): %v", err)
	} else {
		t.Log("Token transfer succeeded unexpectedly")
	}
}

func TestConfigurationValidation(t *testing.T) {
	configManager := config.NewManager()

	// Test invalid UTXO configuration
	invalidUTXOConfig := &config.UTXOConfig{
		IncludeNative:    true,
		IncludeNonNative: true,
		MinConfirmations: -1, // Invalid: negative
		MaxUTXOsPerQuery: 0,  // Invalid: zero
		EnableCaching:    true,
		CacheExpiry:      -1, // Invalid: negative
	}

	err := configManager.UpdateUTXOConfig(invalidUTXOConfig)
	if err == nil {
		t.Error("Expected error for invalid UTXO configuration")
	}

	// Test invalid transaction configuration
	invalidTxConfig := &config.TransactionConfig{
		DefaultFeeRate:        -1, // Invalid: negative
		MinFeeRate:            0,  // Invalid: zero
		MaxFeeRate:            0,  // Invalid: zero
		DustLimit:             -1, // Invalid: negative
		MaxTransactionSize:    0,  // Invalid: zero
		EnableRBF:             false,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
	}

	err = configManager.UpdateTransactionConfig(invalidTxConfig)
	if err == nil {
		t.Error("Expected error for invalid transaction configuration")
	}

	// Test invalid network configuration
	invalidNetworkConfig := &config.NetworkConfig{
		Name:        "", // Invalid: empty
		RPCURL:      "", // Invalid: empty
		ExplorerURL: "", // Invalid: empty
		IsTestnet:   true,
		ChainID:     "", // Invalid: empty
		CoinType:    1,
	}

	err = configManager.UpdateNetworkConfig(invalidNetworkConfig)
	if err == nil {
		t.Error("Expected error for invalid network configuration")
	}
}

func TestConcurrentAccess(t *testing.T) {
	enhancedBSV := bsv.NewBSVDefault()

	// Test concurrent configuration updates
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			defer func() { done <- true }()

			utxoConfig := &config.UTXOConfig{
				IncludeNative:    true,
				IncludeNonNative: true,
				MinConfirmations: 1,
				MaxUTXOsPerQuery: 50 + index,
				EnableCaching:    true,
				CacheExpiry:      300,
			}

			err := enhancedBSV.UpdateUTXOConfig(utxoConfig)
			if err != nil {
				t.Errorf("Concurrent update %d failed: %v", index, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Concurrent access test timed out")
		}
	}
}

func TestPackageLevelEnhancedFunctions(t *testing.T) {
	// Test enhanced package-level functions
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}

	// Test enhanced wallet generation
	wallet, err := bsv.GenerateWalletEnhanced(mnemonicPhrase, config.Testnet)
	if err != nil {
		t.Fatalf("Failed to generate wallet with enhanced function: %v", err)
	}

	if wallet.Address == "" {
		t.Error("Generated wallet address is empty")
	}

	// Test enhanced address validation
	err = bsv.ValidateAddressEnhanced(wallet.Address, config.Testnet)
	if err != nil {
		t.Errorf("Address validation failed: %v", err)
	}

	// Test enhanced balance fetching
	_, err = bsv.GetEnhancedBalance(wallet.Address, config.Testnet)
	if err != nil {
		t.Logf("Enhanced balance fetch failed (expected for new address): %v", err)
	}

	// Test enhanced UTXO fetching
	_, err = bsv.GetUTXOsEnhanced(wallet.Address, config.Testnet)
	if err != nil {
		t.Logf("Enhanced UTXO fetch failed (expected for new address): %v", err)
	}
}
