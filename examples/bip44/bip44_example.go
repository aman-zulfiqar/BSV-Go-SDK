package main

import (
	"fmt"
	"log"

	"github.com/muhammadamman/BSV-Go/pkg/bsv"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
)

func main() {
	fmt.Println("🚀 BSV BIP44 Dynamic Indexing Example")
	fmt.Println("=====================================")

	// Generate a test mnemonic
	mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated mnemonic: %s\n\n", mnemonicPhrase)

	// Create BSV instance
	configManager := config.NewManager()
	configManager.SetNetworkType(config.Testnet) // Use testnet for testing
	bsvInstance := bsv.NewBSV(configManager)

	// 1. Generate default wallet (account 0, change 0, address 0)
	fmt.Println("1. Default BIP44 Path (m/44'/1'/0'/0/0):")
	wallet1, err := bsvInstance.GenerateWallet(mnemonicPhrase)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Address: %s\n", wallet1.Address)
	fmt.Printf("   Private Key: %s\n\n", wallet1.PrivateKey)

	// 2. Generate wallet with custom path (account 0, change 0, address 1)
	fmt.Println("2. Custom BIP44 Path (m/44'/1'/0'/0/1):")
	wallet2, err := bsvInstance.GenerateWalletWithPath(mnemonicPhrase, 0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Address: %s\n", wallet2.Address)
	fmt.Printf("   Private Key: %s\n\n", wallet2.PrivateKey)

	// 3. Generate wallet for different account (account 1, change 0, address 0)
	fmt.Println("3. Different Account (m/44'/1'/1'/0/0):")
	wallet3, err := bsvInstance.GenerateWalletWithPath(mnemonicPhrase, 1, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Address: %s\n", wallet3.Address)
	fmt.Printf("   Private Key: %s\n\n", wallet3.PrivateKey)

	// 4. Generate change address (account 0, change 1, address 0)
	fmt.Println("4. Change Address (m/44'/1'/0'/1/0):")
	wallet4, err := bsvInstance.GenerateWalletWithPath(mnemonicPhrase, 0, 1, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Address: %s\n", wallet4.Address)
	fmt.Printf("   Private Key: %s\n\n", wallet4.PrivateKey)

	// 5. Generate multiple addresses for the same account
	fmt.Println("5. Multiple Addresses for Account 0:")
	for i := uint32(0); i < 5; i++ {
		wallet, err := bsvInstance.GenerateWalletWithPath(mnemonicPhrase, 0, 0, i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("   Address %d: %s\n", i, wallet.Address)
	}

	// 6. Get BIP44 path information
	fmt.Println("\n6. BIP44 Path Information:")
	defaultPath := bsvInstance.GetDefaultBIP44Path()
	fmt.Printf("   Default Path: m/%d'/%d'/%d'/%d/%d\n",
		defaultPath.Purpose,
		defaultPath.CoinType,
		defaultPath.Account,
		defaultPath.Change,
		defaultPath.AddressIndex)

	customPath := bsvInstance.GetBIP44Path(2, 1, 5)
	fmt.Printf("   Custom Path (account=2, change=1, address=5): m/%d'/%d'/%d'/%d/%d\n",
		customPath.Purpose,
		customPath.CoinType,
		customPath.Account,
		customPath.Change,
		customPath.AddressIndex)

	// 7. Demonstrate different coin types
	fmt.Println("\n7. Coin Types:")
	fmt.Printf("   Testnet Coin Type: %d\n", defaultPath.CoinType)

	// Switch to mainnet to show different coin type
	configManager.SetNetworkType(config.Mainnet)
	mainnetBSV := bsv.NewBSV(configManager)
	mainnetPath := mainnetBSV.GetDefaultBIP44Path()
	fmt.Printf("   Mainnet Coin Type: %d\n", mainnetPath.CoinType)

	fmt.Println("\n✅ BIP44 Dynamic Indexing Demo Complete!")
	fmt.Println("\n📚 Key Features Demonstrated:")
	fmt.Println("   • Default BIP44 path generation")
	fmt.Println("   • Custom account, change, and address indexing")
	fmt.Println("   • Multiple address generation from same mnemonic")
	fmt.Println("   • Different coin types for mainnet/testnet")
	fmt.Println("   • Proper HD wallet derivation")
}
