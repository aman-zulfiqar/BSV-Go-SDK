package main

import (
	"fmt"
	"log"

	"github.com/muhammadamman/BSV-Go/pkg/bsv"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/sharding"
)

func main() {
	fmt.Println("🔐 BSV Custodial SDK - Advanced Sharding Example")
	fmt.Println("================================================")

	// Step 1: Generate a random mnemonic
	fmt.Println("\n1. Generating random mnemonic...")
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength256) // 24 words for extra security
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}
	fmt.Printf("✅ Generated 24-word mnemonic: %s\n", mnemonicPhrase)

	// Step 2: Create shards with different thresholds
	fmt.Println("\n2. Creating shards with different configurations...")

	// 2/3 threshold (minimum for basic security)
	fmt.Println("   Creating 2/3 threshold shards...")
	shards23, err := sharding.SplitMnemonic(mnemonicPhrase, 2, 3)
	if err != nil {
		log.Fatalf("Failed to create 2/3 shards: %v", err)
	}
	fmt.Printf("   ✅ Created %d shards, need %d to reconstruct\n",
		shards23.TotalShares, shards23.Threshold)

	// 3/5 threshold (higher security)
	fmt.Println("   Creating 3/5 threshold shards...")
	shards35, err := sharding.SplitMnemonic(mnemonicPhrase, 3, 5)
	if err != nil {
		log.Fatalf("Failed to create 3/5 shards: %v", err)
	}
	fmt.Printf("   ✅ Created %d shards, need %d to reconstruct\n",
		shards35.TotalShares, shards35.Threshold)

	// Step 3: Demonstrate shard validation
	fmt.Println("\n3. Validating shards...")
	for i, shard := range shards23.Shards {
		isValid := sharding.ValidateShard(shard)
		fmt.Printf("   Shard %d: %s (valid: %t)\n", i+1, shard[:20]+"...", isValid)
	}

	// Step 4: Simulate different recovery scenarios
	fmt.Println("\n4. Simulating recovery scenarios...")

	// Scenario 1: Recover with first 2 shards (should work)
	fmt.Println("   Scenario 1: Recovering with shards 1 & 2...")
	recovered1, err := sharding.CombineShards([]string{shards23.Shards[0], shards23.Shards[1]})
	if err != nil {
		log.Printf("   ❌ Recovery failed: %v", err)
	} else {
		if recovered1 == mnemonicPhrase {
			fmt.Println("   ✅ Recovery successful!")
		} else {
			fmt.Println("   ❌ Recovery failed - mnemonic mismatch")
		}
	}

	// Scenario 2: Recover with shards 1 & 3 (should work)
	fmt.Println("   Scenario 2: Recovering with shards 1 & 3...")
	recovered2, err := sharding.CombineShards([]string{shards23.Shards[0], shards23.Shards[2]})
	if err != nil {
		log.Printf("   ❌ Recovery failed: %v", err)
	} else {
		if recovered2 == mnemonicPhrase {
			fmt.Println("   ✅ Recovery successful!")
		} else {
			fmt.Println("   ❌ Recovery failed - mnemonic mismatch")
		}
	}

	// Scenario 3: Try to recover with only 1 shard (should fail)
	fmt.Println("   Scenario 3: Trying to recover with only 1 shard...")
	_, err = sharding.CombineShards([]string{shards23.Shards[0]})
	if err != nil {
		fmt.Printf("   ✅ Recovery correctly failed (insufficient shards): %v\n", err)
	} else {
		fmt.Println("   ❌ Recovery should have failed but didn't!")
	}

	// Scenario 4: Recover 3/5 threshold with 3 shards
	fmt.Println("   Scenario 4: Recovering 3/5 threshold with shards 1, 2, & 3...")
	recovered4, err := sharding.CombineShards([]string{shards35.Shards[0], shards35.Shards[1], shards35.Shards[2]})
	if err != nil {
		log.Printf("   ❌ Recovery failed: %v", err)
	} else {
		if recovered4 == mnemonicPhrase {
			fmt.Println("   ✅ Recovery successful!")
		} else {
			fmt.Println("   ❌ Recovery failed - mnemonic mismatch")
		}
	}

	// Step 5: Demonstrate wallet generation from recovered mnemonic
	fmt.Println("\n5. Generating wallet from recovered mnemonic...")
	wallet, err := bsv.GenerateWalletEnhanced(recovered1, config.Testnet)
	if err != nil {
		log.Fatalf("Failed to generate wallet: %v", err)
	}
	fmt.Printf("✅ BSV Address: %s\n", wallet.Address)

	// Step 6: Show security implications
	fmt.Println("\n6. Security Analysis:")
	fmt.Println("   🔒 2/3 Threshold:")
	fmt.Println("      - Any 2 of 3 shards can reconstruct the mnemonic")
	fmt.Println("      - Good balance between security and convenience")
	fmt.Println("      - Recommended for most use cases")
	fmt.Println()
	fmt.Println("   🔒 3/5 Threshold:")
	fmt.Println("      - Any 3 of 5 shards can reconstruct the mnemonic")
	fmt.Println("      - Higher security, requires more shards")
	fmt.Println("      - Good for high-value wallets")
	fmt.Println()
	fmt.Println("   🛡️  Security Best Practices:")
	fmt.Println("      - Store shards in different physical locations")
	fmt.Println("      - Use encrypted storage for shards")
	fmt.Println("      - Never store all shards together")
	fmt.Println("      - Consider using different storage mediums")

	// Step 7: Demonstrate practical usage patterns
	fmt.Println("\n7. Practical Usage Patterns:")
	fmt.Println("   📱 Frontend (Shard A):")
	fmt.Printf("      - Store shard: %s\n", shards23.Shards[0][:20]+"...")
	fmt.Println("      - Handle user interactions")
	fmt.Println("      - Combine with backend shard for transactions")
	fmt.Println()
	fmt.Println("   🖥️  Backend (Shard B):")
	fmt.Printf("      - Store shard: %s\n", shards23.Shards[1][:20]+"...")
	fmt.Println("      - Handle transaction processing")
	fmt.Println("      - Combine with frontend shard for signing")
	fmt.Println()
	fmt.Println("   ☁️  Recovery Storage (Shard C):")
	fmt.Printf("      - Store shard: %s\n", shards23.Shards[2][:20]+"...")
	fmt.Println("      - Google Drive, cloud storage, etc.")
	fmt.Println("      - Used only for recovery scenarios")

	fmt.Println("\n🎉 Advanced sharding example completed!")
	fmt.Println("\n💡 Key Takeaways:")
	fmt.Println("   1. Sharding provides excellent security for mnemonic storage")
	fmt.Println("   2. 2/3 threshold is practical for most applications")
	fmt.Println("   3. Higher thresholds provide more security but less convenience")
	fmt.Println("   4. Never store all shards in the same location")
	fmt.Println("   5. Always validate shards before use")
}
