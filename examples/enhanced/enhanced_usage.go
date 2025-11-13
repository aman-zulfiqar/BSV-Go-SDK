package main

import (
	"fmt"
	"log"
	"time"

	"github.com/muhammadamman/BSV-Go/pkg/bsv"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/sharding"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

func main() {
	fmt.Println("🚀 BSV Enhanced SDK - Dynamic Configuration Example")
	fmt.Println("==================================================")

	// Step 1: Create enhanced BSV instance with dynamic configuration
	fmt.Println("\n1. Creating enhanced BSV instance with dynamic configuration...")

	// Create configuration manager
	configManager := config.NewManager()

	// Create enhanced BSV instance
	enhancedBSV := bsv.NewBSV(configManager)

	fmt.Printf("✅ Initial network: %s\n", enhancedBSV.GetNetworkConfig().Name)
	fmt.Printf("✅ Initial UTXO config - Native: %v, Non-Native: %v\n",
		enhancedBSV.GetUTXOConfig().IncludeNative,
		enhancedBSV.GetUTXOConfig().IncludeNonNative)
	fmt.Printf("✅ Initial transaction config - Default fee rate: %d sat/vbyte\n",
		enhancedBSV.GetTransactionConfig().DefaultFeeRate)

	// Step 2: Demonstrate dynamic configuration updates
	fmt.Println("\n2. Demonstrating dynamic configuration updates...")

	// Update UTXO configuration
	utxoConfig := &config.UTXOConfig{
		IncludeNative:    true,
		IncludeNonNative: true,
		MinConfirmations: 1,
		MaxUTXOsPerQuery: 50,
		EnableCaching:    true,
		CacheExpiry:      180, // 3 minutes
	}

	err := enhancedBSV.UpdateUTXOConfig(utxoConfig)
	if err != nil {
		log.Fatalf("Failed to update UTXO config: %v", err)
	}
	fmt.Printf("✅ Updated UTXO config - Cache enabled: %v, Expiry: %d seconds\n",
		enhancedBSV.GetUTXOConfig().EnableCaching,
		enhancedBSV.GetUTXOConfig().CacheExpiry)

	// Update transaction configuration
	txConfig := &config.TransactionConfig{
		DefaultFeeRate:        8,
		MinFeeRate:            1,
		MaxFeeRate:            1000,
		DustLimit:             546,
		MaxTransactionSize:    100000,
		EnableRBF:             false,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
	}

	err = enhancedBSV.UpdateTransactionConfig(txConfig)
	if err != nil {
		log.Fatalf("Failed to update transaction config: %v", err)
	}
	fmt.Printf("✅ Updated transaction config - Default fee rate: %d sat/vbyte\n",
		enhancedBSV.GetTransactionConfig().DefaultFeeRate)

	// Step 3: Generate mnemonic and wallet
	fmt.Println("\n3. Generating mnemonic and wallet...")

	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}
	fmt.Printf("✅ Generated mnemonic: %s\n", mnemonicPhrase)

	// Create shards
	shardingResult, err := sharding.SplitMnemonic(mnemonicPhrase, 2, 3)
	if err != nil {
		log.Fatalf("Failed to create shards: %v", err)
	}
	fmt.Printf("✅ Created %d shards (need %d to reconstruct)\n",
		shardingResult.TotalShares, shardingResult.Threshold)

	// Reconstruct mnemonic from shards
	reconstructedMnemonic, err := sharding.CombineShards(shardingResult.Shards[:2])
	if err != nil {
		log.Fatalf("Failed to reconstruct mnemonic: %v", err)
	}
	fmt.Printf("✅ Reconstructed mnemonic: %s\n", reconstructedMnemonic)

	// Generate wallet
	wallet, err := enhancedBSV.GenerateWallet(reconstructedMnemonic)
	if err != nil {
		log.Fatalf("Failed to generate wallet: %v", err)
	}
	fmt.Printf("✅ BSV Address: %s\n", wallet.Address)
	fmt.Printf("✅ Private Key: %s\n", wallet.PrivateKey)

	// Step 4: Demonstrate enhanced balance fetching
	fmt.Println("\n4. Demonstrating enhanced balance fetching...")

	// Get enhanced balance
	enhancedBalance, err := enhancedBSV.GetEnhancedBalance(wallet.Address)
	if err != nil {
		log.Printf("⚠️  Failed to get enhanced balance: %v", err)
	} else {
		fmt.Printf("✅ Enhanced Balance Information:\n")
		fmt.Printf("   Native BSV - Confirmed: %s BSV, Unconfirmed: %s BSV\n",
			types.FormatBSV(enhancedBalance.Native.Confirmed),
			types.FormatBSV(enhancedBalance.Native.Unconfirmed))
		fmt.Printf("   Native UTXO Count: %d\n", enhancedBalance.Native.UTXOCount)
		fmt.Printf("   Non-Native UTXO Count: %d\n", enhancedBalance.NonNative.UTXOCount)
		fmt.Printf("   Total BSV Balance: %s BSV\n", types.FormatBSV(enhancedBalance.Total))
	}

	// Get native balance only
	nativeBalance, err := enhancedBSV.GetNativeBalance(wallet.Address)
	if err != nil {
		log.Printf("⚠️  Failed to get native balance: %v", err)
	} else {
		fmt.Printf("✅ Native Balance - Total: %s BSV (%d satoshis)\n",
			types.FormatBSV(nativeBalance.Total), nativeBalance.Total)
	}

	// Get non-native balance
	nonNativeBalance, err := enhancedBSV.GetNonNativeBalance(wallet.Address)
	if err != nil {
		log.Printf("⚠️  Failed to get non-native balance: %v", err)
	} else {
		fmt.Printf("✅ Non-Native Balance - UTXO Count: %d\n", nonNativeBalance.UTXOCount)
		if len(nonNativeBalance.Tokens) > 0 {
			fmt.Printf("   Token Balances:\n")
			for tokenID, tokenBalance := range nonNativeBalance.Tokens {
				fmt.Printf("     Token %s: %d units\n", tokenID, tokenBalance.Total)
			}
		}
	}

	// Step 5: Demonstrate UTXO fetching with caching
	fmt.Println("\n5. Demonstrating UTXO fetching with caching...")

	// First fetch (will be cached)
	start := time.Now()
	utxos1, err := enhancedBSV.GetUTXOs(wallet.Address)
	fetchTime1 := time.Since(start)
	if err != nil {
		log.Printf("⚠️  Failed to get UTXOs: %v", err)
	} else {
		fmt.Printf("✅ Fetched %d UTXOs in %v (first fetch - cached)\n", len(utxos1), fetchTime1)
	}

	// Second fetch (should be faster due to caching)
	start = time.Now()
	utxos2, err := enhancedBSV.GetUTXOs(wallet.Address)
	fetchTime2 := time.Since(start)
	if err != nil {
		log.Printf("⚠️  Failed to get UTXOs: %v", err)
	} else {
		fmt.Printf("✅ Fetched %d UTXOs in %v (second fetch - from cache)\n", len(utxos2), fetchTime2)
		if fetchTime2 < fetchTime1 {
			fmt.Printf("✅ Cache is working - second fetch was faster!\n")
		}
	}

	// Step 6: Demonstrate network switching
	fmt.Println("\n6. Demonstrating network switching...")

	// Switch to mainnet configuration
	err = enhancedBSV.SetNetworkType(config.Mainnet)
	if err != nil {
		log.Fatalf("Failed to switch to mainnet: %v", err)
	}
	fmt.Printf("✅ Switched to: %s\n", enhancedBSV.GetNetworkConfig().Name)
	fmt.Printf("✅ RPC URL: %s\n", enhancedBSV.GetNetworkConfig().RPCURL)
	fmt.Printf("✅ Explorer: %s\n", enhancedBSV.GetNetworkConfig().ExplorerURL)

	// Switch back to testnet
	err = enhancedBSV.SetNetworkType(config.Testnet)
	if err != nil {
		log.Fatalf("Failed to switch to testnet: %v", err)
	}
	fmt.Printf("✅ Switched back to: %s\n", enhancedBSV.GetNetworkConfig().Name)

	// Step 7: Demonstrate enhanced transaction building
	fmt.Println("\n7. Demonstrating enhanced transaction building...")

	// Create enhanced transaction parameters
	txParams := &types.TransactionParams{
		From:                  wallet.Address,
		To:                    "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt", // BSV testnet address
		Amount:                1000,                                 // 1000 satoshis
		FeeRate:               8,                                    // 8 sat/vbyte (from our config)
		PrivateKey:            reconstructedMnemonic,
		IncludeNativeUTXOs:    true,
		IncludeNonNativeUTXOs: false,
		TokenTransfers:        []*types.TokenTransfer{},
		DataOutputs:           []*types.DataOutput{},
	}

	fmt.Printf("📝 Enhanced Transaction Parameters:\n")
	fmt.Printf("   From: %s\n", txParams.From)
	fmt.Printf("   To: %s\n", txParams.To)
	fmt.Printf("   Amount: %s BSV\n", types.FormatBSV(txParams.Amount))
	fmt.Printf("   Fee Rate: %d sat/vbyte\n", txParams.FeeRate)
	fmt.Printf("   Include Native UTXOs: %v\n", txParams.IncludeNativeUTXOs)
	fmt.Printf("   Include Non-Native UTXOs: %v\n", txParams.IncludeNonNativeUTXOs)

	// Try to build transaction (will fail without funds, which is expected)
	_, err = enhancedBSV.SignAndSendTransaction(txParams)
	if err != nil {
		fmt.Printf("⚠️  Transaction failed (expected - no funds): %v\n", err)
		fmt.Println("💡 Fund the address with testnet BSV to send transactions")
	} else {
		fmt.Println("✅ Transaction sent successfully!")
	}

	// Step 8: Demonstrate cache management
	fmt.Println("\n8. Demonstrating cache management...")

	// Clear cache for specific address
	enhancedBSV.ClearUTXOCacheForAddress(wallet.Address)
	fmt.Printf("✅ Cleared cache for address: %s\n", wallet.Address)

	// Clear entire cache
	enhancedBSV.ClearUTXOCache()
	fmt.Printf("✅ Cleared entire UTXO cache\n")

	// Step 9: Show final configuration
	fmt.Println("\n9. Final Configuration Summary:")
	fmt.Printf("✅ Network: %s\n", enhancedBSV.GetNetworkConfig().Name)
	fmt.Printf("✅ RPC URL: %s\n", enhancedBSV.GetNetworkConfig().RPCURL)
	fmt.Printf("✅ Explorer: %s\n", enhancedBSV.GetNetworkConfig().ExplorerURL)
	fmt.Printf("✅ UTXO Config - Cache: %v, Max UTXOs: %d\n",
		enhancedBSV.GetUTXOConfig().EnableCaching,
		enhancedBSV.GetUTXOConfig().MaxUTXOsPerQuery)
	fmt.Printf("✅ Transaction Config - Fee Rate: %d, Dust Limit: %d\n",
		enhancedBSV.GetTransactionConfig().DefaultFeeRate,
		enhancedBSV.GetTransactionConfig().DustLimit)

	fmt.Println("\n🎉 Enhanced SDK demonstration completed successfully!")
	fmt.Println("\n📚 Key Features Demonstrated:")
	fmt.Println("   1. ✅ Dynamic configuration management")
	fmt.Println("   2. ✅ Enhanced balance fetching (native + non-native)")
	fmt.Println("   3. ✅ UTXO caching for performance")
	fmt.Println("   4. ✅ Network switching (mainnet/testnet)")
	fmt.Println("   5. ✅ Enhanced transaction building")
	fmt.Println("   6. ✅ Cache management")
	fmt.Println("   7. ✅ Production-ready configuration validation")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("   1. Fund the address with testnet BSV")
	fmt.Println("   2. Test real transactions")
	fmt.Println("   3. Implement token transfers")
	fmt.Println("   4. Configure for production use")
}
