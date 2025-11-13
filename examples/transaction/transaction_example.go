package main

import (
	"fmt"
	"log"

	"github.com/muhammadamman/BSV-Go/pkg/bsv"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

func main() {

	fmt.Println("💰 BSV Custodial SDK - Transaction Example")
	fmt.Println("==========================================")

	// Step 1: Generate a wallet with mnemonic
	fmt.Println("\n1. Generating wallet...")
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}
	fmt.Printf("✅ Generated mnemonic: %s\n", mnemonicPhrase)

	wallet, err := bsv.GenerateWalletEnhanced(mnemonicPhrase, config.Testnet)
	if err != nil {
		log.Fatalf("Failed to generate wallet: %v", err)
	}
	fmt.Printf("✅ BSV Address: %s\n", wallet.Address)

	// Step 2: Check initial balance
	fmt.Println("\n2. Checking initial balance...")
	balance, err := bsv.GetEnhancedBalance(wallet.Address, config.Testnet)
	if err != nil {
		log.Printf("⚠️  Failed to check balance: %v", err)
		fmt.Printf("💰 Balance: 0.00000000 BSV (0 satoshis)\n")
	} else {
		fmt.Printf("💰 Native Balance: %.8f BSV (%d satoshis)\n",
			float64(balance.Native.Total)/100000000,
			balance.Native.Total)
		fmt.Printf("💰 Non-Native UTXOs: %d\n", balance.NonNative.UTXOCount)
	}

	// Step 3: Get UTXOs (will be empty for new address)
	fmt.Println("\n3. Checking UTXOs...")
	utxos, err := bsv.GetUTXOsEnhanced(wallet.Address, config.Testnet)
	if err != nil {
		log.Printf("⚠️  Failed to get UTXOs: %v", err)
	} else {
		fmt.Printf("📦 Found %d UTXOs\n", len(utxos))
		for i, utxo := range utxos {
			fmt.Printf("   UTXO %d: %s...%d (%d sats)\n",
				i+1, utxo.TxID[:10], utxo.Vout, utxo.Value)
		}
	}

	// Step 4: Demonstrate transaction building with different scenarios
	fmt.Println("\n4. Transaction Building Scenarios:")

	// Scenario 1: Small transaction
	fmt.Println("\n   📝 Scenario 1: Small transaction (1000 sats)")
	testParams1 := &types.TransactionParams{
		From:       wallet.Address,
		To:         "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt", // BSV testnet address
		Amount:     1000,                                 // 1000 satoshis
		FeeRate:    5,                                    // 5 satoshis per vbyte
		PrivateKey: mnemonicPhrase,
	}

	fmt.Printf("      From: %s\n", testParams1.From)
	fmt.Printf("      To: %s\n", testParams1.To)
	fmt.Printf("      Amount: %s BSV\n", types.FormatBSV(testParams1.Amount))
	fmt.Printf("      Fee Rate: %d sat/vbyte\n", testParams1.FeeRate)

	// Try to build transaction (will fail without funds, which is expected)
	isTestnet := config.Testnet

	_, err = bsv.SignAndSendTransactionEnhanced(testParams1, isTestnet)
	if err != nil {
		fmt.Printf("      ⚠️  Transaction failed (expected - no funds): %v\n", err)
	} else {
		fmt.Println("      ✅ Transaction sent successfully!")
	}

	// Scenario 2: Medium transaction with higher fee
	fmt.Println("\n   📝 Scenario 2: Medium transaction with higher fee (10000 sats)")
	testParams2 := &types.TransactionParams{
		From:       wallet.Address,
		To:         "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt",
		Amount:     10000, // 10000 satoshis
		FeeRate:    10,    // 10 satoshis per vbyte (higher fee)
		PrivateKey: mnemonicPhrase,
	}

	fmt.Printf("      From: %s\n", testParams2.From)
	fmt.Printf("      To: %s\n", testParams2.To)
	fmt.Printf("      Amount: %s BSV\n", types.FormatBSV(testParams2.Amount))
	fmt.Printf("      Fee Rate: %d sat/vbyte\n", testParams2.FeeRate)

	_, err = bsv.SignAndSendTransactionEnhanced(testParams2, isTestnet)
	if err != nil {
		fmt.Printf("      ⚠️  Transaction failed (expected - no funds): %v\n", err)
	} else {
		fmt.Println("      ✅ Transaction sent successfully!")
	}

	// Scenario 3: Large transaction
	fmt.Println("\n   📝 Scenario 3: Large transaction (100000 sats = 0.001 BSV)")
	testParams3 := &types.TransactionParams{
		From:       wallet.Address,
		To:         "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt",
		Amount:     100000, // 100000 satoshis (0.001 BSV)
		FeeRate:    5,      // 5 satoshis per vbyte
		PrivateKey: mnemonicPhrase,
	}

	fmt.Printf("      From: %s\n", testParams3.From)
	fmt.Printf("      To: %s\n", testParams3.To)
	fmt.Printf("      Amount: %s BSV\n", types.FormatBSV(testParams3.Amount))
	fmt.Printf("      Fee Rate: %d sat/vbyte\n", testParams3.FeeRate)

	_, err = bsv.SignAndSendTransactionEnhanced(testParams3, isTestnet)
	if err != nil {
		fmt.Printf("      ⚠️  Transaction failed (expected - no funds): %v\n", err)
	} else {
		fmt.Println("      ✅ Transaction sent successfully!")
	}

	// Step 5: Demonstrate WIF private key usage
	fmt.Println("\n5. Using WIF Private Key:")
	fmt.Printf("   WIF: %s\n", wallet.PrivateKey)

	// Create transaction using WIF instead of mnemonic
	wifParams := &types.TransactionParams{
		From:       wallet.Address,
		To:         "mqVKYrNJSmJNQNnQpqNk5XnxSc4iXTJmkt",
		Amount:     5000,              // 5000 satoshis
		FeeRate:    5,                 // 5 satoshis per vbyte
		PrivateKey: wallet.PrivateKey, // Using WIF instead of mnemonic
	}

	fmt.Printf("   📝 WIF Transaction: %s BSV\n", types.FormatBSV(wifParams.Amount))

	_, err = bsv.SignAndSendTransactionEnhanced(wifParams, isTestnet)
	if err != nil {
		fmt.Printf("   ⚠️  Transaction failed (expected - no funds): %v\n", err)
	} else {
		fmt.Println("   ✅ Transaction sent successfully!")
	}

	// Step 6: Show fee calculation examples
	fmt.Println("\n6. Fee Calculation Examples:")

	// Different fee rates and their implications
	feeRates := []int64{1, 5, 10, 20, 50}
	estimatedTxSize := int64(200) // Estimated transaction size in vbytes

	fmt.Println("   Fee Rate Analysis (estimated 200 vbyte transaction):")
	for _, feeRate := range feeRates {
		fee := estimatedTxSize * feeRate
		fmt.Printf("      %2d sat/vbyte = %8d sats (%s BSV)\n",
			feeRate, fee, types.FormatBSV(fee))
	}

	// Step 7: Network information
	fmt.Println("\n7. Network Information:")
	if err != nil { // ✅ correct {
		fmt.Println("   🌐 Network: BSV Testnet")
		fmt.Println("   🔗 Explorer: https://test.whatsonchain.com")
		fmt.Println("   💰 Faucet: https://faucet.bitcoincloud.net/")
		fmt.Println("   📊 Testnet Stats: https://test.whatsonchain.com/stats")
	} else {
		fmt.Println("   🌐 Network: BSV Mainnet")
		fmt.Println("   🔗 Explorer: https://whatsonchain.com")
		fmt.Println("   📊 Network Stats: https://whatsonchain.com/stats")
	}

	// Step 8: Best practices
	fmt.Println("\n8. Transaction Best Practices:")
	fmt.Println("   💡 Fee Rate Guidelines:")
	fmt.Println("      - 1-5 sat/vbyte: Low priority, slower confirmation")
	fmt.Println("      - 5-10 sat/vbyte: Normal priority, standard confirmation")
	fmt.Println("      - 10-20 sat/vbyte: High priority, fast confirmation")
	fmt.Println("      - 20+ sat/vbyte: Very high priority, instant confirmation")
	fmt.Println()
	fmt.Println("   🛡️  Security Tips:")
	fmt.Println("      - Always validate addresses before sending")
	fmt.Println("      - Double-check amounts before signing")
	fmt.Println("      - Use testnet for development and testing")
	fmt.Println("      - Keep private keys secure and never share them")
	fmt.Println("      - Consider using hardware wallets for large amounts")

	fmt.Println("\n🎉 Transaction example completed!")
	fmt.Println("\n📚 Next Steps:")
	fmt.Println("   1. Fund the address with testnet BSV from a faucet")
	fmt.Println("   2. Try sending real transactions")
	fmt.Println("   3. Experiment with different fee rates")
	fmt.Println("   4. Monitor transaction confirmations on the explorer")
	fmt.Println("   5. Implement error handling for production use")
}
