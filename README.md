# BSV Go SDK - Production Ready

A comprehensive Bitcoin SV (BSV) SDK written in Go, providing dynamic configuration, secure mnemonic management, BIP44 HD wallets, and enhanced transaction utilities.

## 🚀 Features

- **🎯 Dynamic Configuration**: Runtime configuration management for networks, UTXOs, and transactions
- **🔑 BIP44 HD Wallets**: Dynamic indexing with custom account, change, and address paths
- **🔐 Mnemonic Generation**: BIP39-compliant 12/24-word mnemonic generation
- **🛡️ Shamir Secret Sharing**: 2/3 threshold sharding for secure key management
- **💰 Enhanced UTXO Management**: Native and non-native UTXO support with dynamic summing
- **📝 Transaction Building**: Native and non-native BSV transaction support
- **🌐 Network Support**: Both mainnet and testnet BSV networks with dynamic switching
- **⚡ Production Ready**: Comprehensive error handling, caching, and validation
- **🔄 Caching**: Configurable UTXO caching for performance optimization
- **🎛️ Runtime Updates**: Dynamic configuration updates without restart

## 📁 Project Structure

```
BSV-Go/
├── cmd/                   # Runnable example entrypoint (main.go)
├── pkg/                   # Public SDK surface
│   ├── bsv/               # SDK façade
│   │   ├── transaction/   # Transaction building/signing/broadcast
│   │   ├── utxo/          # UTXO retrieval, caching, selection
│   │   └── wallet/        # BIP39/BIP32/BIP44 derivation, addresses
│   ├── config/            # Runtime configuration manager
│   ├── mnemonic/          # Mnemonic generation/validation
│   ├── sharding/          # Demo sharding (see security note)
│   └── types/             # Shared DTOs and helpers
├── examples/              # Focused usage examples
│   ├── enhanced/          # Dynamic config + caching + UX demo
│   ├── bip44/             # HD wallet derivation paths
│   ├── sharding/          # Sharding flows
│   └── transaction/       # Transaction scenarios
├── tests/                 # Test suite
├── .github/workflows/     # CI (tidy/fmt/lint/test)
├── go.mod                 # Module definition
└── README.md
```

### Design & architecture
- Config-first design: all behavior is driven by `pkg/config.Manager` and read at call time (no restarts).
- Separation of concerns:
  - `pkg/bsv/wallet`: BIP39 seed → BIP32 xprv → BIP44 derivation → WIF and P2PKH address.
  - `pkg/bsv/utxo`: HTTP client, retries, caching, enhanced balance, UTXO selection heuristics.
  - `pkg/bsv/transaction`: builds wire transactions, estimates fees, signs inputs, broadcasts.
- Public API is a thin façade in `pkg/bsv` delegating to the above, returning strongly-typed `pkg/types`.

## 🛠️ Installation

```bash
# Initialize your Go module
go mod init your-project

# Add the BSV Go SDK dependency
go get github.com/muhammadamman/BSV-Go@latest
```

### Requirements
- Go 1.21+
- Internet access for explorer API calls in examples/tests (Whatsonchain-compatible)
- Optional: `golangci-lint` if you want to run `make lint`

## 🎯 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/muhammadamman/BSV-Go/pkg/bsv"
    "github.com/muhammadamman/BSV-Go/pkg/config"
    "github.com/muhammadamman/BSV-Go/pkg/mnemonic"
)

func main() {
    // 1. Generate random mnemonic
    mnemonicPhrase, err := mnemonic.Generate(mnemonic.Strength128) // 12 words
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Generated mnemonic: %s\n", mnemonicPhrase)
    
    // 2. Create BSV instance with dynamic configuration
    configManager := config.NewManager()
    configManager.SetNetworkType(config.Testnet)
    bsvInstance := bsv.NewBSV(configManager)
    
    // 3. Generate default wallet (BIP44: m/44'/1'/0'/0/0)
    wallet1, err := bsvInstance.GenerateWallet(mnemonicPhrase)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Default Address: %s\n", wallet1.Address)
    
    // 4. Generate custom wallet (BIP44: m/44'/1'/0'/0/1)
    wallet2, err := bsvInstance.GenerateWalletWithPath(mnemonicPhrase, 0, 0, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Custom Address: %s\n", wallet2.Address)
    
    // 5. Get enhanced balance (native + non-native UTXOs)
    balance, err := bsvInstance.GetEnhancedBalance(wallet1.Address)
    if err != nil {
        log.Printf("Balance check failed: %v", err)
    } else {
        fmt.Printf("Native Balance: %d satoshis\n", balance.Native.Total)
        fmt.Printf("Non-Native UTXOs: %d\n", balance.NonNative.UTXOCount)
    }
}
```

### Advanced Usage with Dynamic Configuration

```go
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
    // Create configuration manager
    configManager := config.NewManager()
    
    // Configure for production
    utxoConfig := &config.UTXOConfig{
        IncludeNative:    true,
        IncludeNonNative: true,
        MinConfirmations: 3,
        EnableCaching:    true,
        CacheExpiry:      300,
        MaxUTXOsPerQuery: 50,
    }
    configManager.UpdateUTXOConfig(utxoConfig)
    
    // Create BSV instance
    bsvInstance := bsv.NewBSV(configManager)
    
    // Switch to mainnet for production
    bsvInstance.SetNetworkType(config.Mainnet)
    
    // Generate wallet
    mnemonicPhrase, _ := mnemonic.Generate(mnemonic.Strength128)
    wallet, _ := bsvInstance.GenerateWallet(mnemonicPhrase)
    
    // Build and send transaction
    txParams := &types.TransactionParams{
        From:                wallet.Address,
        To:                  "recipient_address",
        Amount:              1000, // 1000 satoshis
        PrivateKey:          mnemonicPhrase,
        IncludeNativeUTXOs:  true,
        IncludeNonNativeUTXOs: false,
    }
    
    result, err := bsvInstance.SignAndSendTransaction(txParams)
    if err != nil {
        log.Printf("Transaction failed: %v", err)
    } else {
        fmt.Printf("Transaction sent: %s\n", result.TxID)
    }
}
```

## 🔑 BIP44 HD Wallet Support

The SDK provides full BIP44 HD wallet support with dynamic indexing:

```go
// Generate default wallet (m/44'/1'/0'/0/0 for testnet)
wallet1, err := bsvInstance.GenerateWallet(mnemonic)

// Generate custom wallet with specific path
wallet2, err := bsvInstance.GenerateWalletWithPath(mnemonic, 0, 0, 1) // m/44'/1'/0'/0/1

// Generate wallet for different account
wallet3, err := bsvInstance.GenerateWalletWithPath(mnemonic, 1, 0, 0) // m/44'/1'/1'/0/0

// Generate change address
changeAddr, err := bsvInstance.GenerateWalletWithPath(mnemonic, 0, 1, 0) // m/44'/1'/0'/1/0

// Get BIP44 path information
defaultPath := bsvInstance.GetDefaultBIP44Path()
fmt.Printf("Path: m/%d'/%d'/%d'/%d/%d\n", 
    defaultPath.Purpose, defaultPath.CoinType, defaultPath.Account, 
    defaultPath.Change, defaultPath.AddressIndex)
```

## 🔐 Security Features

- **🛡️ Shamir Secret Sharing**: 2/3 threshold for mnemonic reconstruction
- **✅ BIP39 Compliance**: Standard mnemonic generation and validation
- **🔑 BIP44 HD Wallets**: Hierarchical deterministic wallet generation
- **🎲 Secure Random**: Cryptographically secure random number generation
- **🚫 No Key Storage**: Keys generated on-demand from shards
- **🔒 Thread Safety**: All operations are thread-safe with proper locking

## 📊 Enhanced UTXO Management

```go
// Get enhanced balance (native + non-native UTXOs)
balance, err := bsvInstance.GetEnhancedBalance(address)
if err == nil {
    fmt.Printf("Native BSV: %d satoshis\n", balance.Native.Total)
    fmt.Printf("Native UTXOs: %d\n", balance.Native.UTXOCount)
    fmt.Printf("Non-Native UTXOs: %d\n", balance.NonNative.UTXOCount)
    
    // Check for specific tokens
    for tokenID, tokenBalance := range balance.NonNative.Tokens {
        fmt.Printf("Token %s: %d units\n", tokenID, tokenBalance.Total)
    }
}

// Get UTXOs with caching
utxos, err := bsvInstance.GetUTXOs(address)
```

## 🎛️ Dynamic Configuration

```go
// Update UTXO configuration
utxoConfig := &config.UTXOConfig{
    IncludeNative:    true,
    IncludeNonNative: true,
    MinConfirmations: 3,
    EnableCaching:    true,
    CacheExpiry:      300,
    MaxUTXOsPerQuery: 50,
}
bsvInstance.UpdateUTXOConfig(utxoConfig)

// Update transaction configuration
txConfig := &config.TransactionConfig{
    DefaultFeeRate:       5,
    MinFeeRate:           1,
    MaxFeeRate:           100,
    MaxTransactionSize:   1000000,
    DustLimit:            546,
    EnableRBF:            true,
}
bsvInstance.UpdateTransactionConfig(txConfig)

// Switch networks dynamically
bsvInstance.SetNetworkType(config.Mainnet) // Switch to mainnet
bsvInstance.SetNetworkType(config.Testnet) // Switch to testnet
```

### Configuration defaults
The defaults in `pkg/config/manager.go`:

```go
// Network (default: Testnet)
Name: "BSV Testnet"
RPCURL: "https://api.whatsonchain.com/v1/bsv/test"
ExplorerURL: "https://test.whatsonchain.com"
IsTestnet: true
CoinType: 1

// UTXO
IncludeNative:    true
IncludeNonNative: true
MinConfirmations: 1
MaxUTXOsPerQuery: 100
EnableCaching:    true
CacheExpiry:      300 // seconds

// Transaction
DefaultFeeRate:        5
MinFeeRate:            1
MaxFeeRate:            1000
DustLimit:             546
MaxTransactionSize:    100000
EnableRBF:             false
IncludeNativeUTXOs:    true
IncludeNonNativeUTXOs: false
```

## 📝 Transaction Building

```go
// Native BSV transaction
txParams := &types.TransactionParams{
    From:                senderAddress,
    To:                  recipientAddress,
    Amount:              1000, // 1000 satoshis
    PrivateKey:          mnemonic,
    IncludeNativeUTXOs:  true,
    IncludeNonNativeUTXOs: false,
}

// Token transfer transaction
txParams.TokenTransfers = []*types.TokenTransfer{
    {
        TokenID: "token_id_here",
        To:      recipientAddress,
        Amount:  100,
    },
}

// Data output transaction
txParams.DataOutputs = []*types.DataOutput{
    {
        Data: "48656c6c6f", // "Hello" in hex
    },
}

// Build and send transaction
result, err := bsvInstance.SignAndSendTransaction(txParams)
if err != nil {
    log.Printf("Transaction failed: %v", err)
} else {
    fmt.Printf("Transaction ID: %s\n", result.TxID)
    fmt.Printf("Fee: %d satoshis\n", result.Fee)
    fmt.Printf("Explorer: %s\n", result.ExplorerURL)
}
```

## 🛡️ Shamir Secret Sharing

Important: The current sharding implementation in `pkg/sharding` is a demo for examples/tests and is not a production-grade threshold scheme. Do not rely on it for real secret sharing until replaced with a proper SSS implementation.

```go
// Generate mnemonic
mnemonic, err := mnemonic.Generate(mnemonic.Strength128)

// Split into shards (2/3 threshold)
shards, err := sharding.SplitMnemonic(mnemonic, 2, 3)
fmt.Printf("Created %d shards, need %d to reconstruct\n", 
    shards.TotalShares, shards.Threshold)

// Reconstruct from any 2 shards
reconstructed, err := sharding.CombineShards(shards.Shards[:2])
if reconstructed == mnemonic {
    fmt.Println("✅ Reconstruction successful!")
}

// Validate shards
for i, shard := range shards.Shards {
    isValid := sharding.ValidateShard(shard)
    fmt.Printf("Shard %d valid: %v\n", i+1, isValid)
}
```

## 🌐 Network Support

### Testnet Configuration
```go
configManager := config.NewManager()
configManager.SetNetworkType(config.Testnet)
bsvInstance := bsv.NewBSV(configManager)

// Coin type: 1
// RPC URL: https://api.whatsonchain.com/v1/bsv/test
// Explorer: https://test.whatsonchain.com
```

### Mainnet Configuration
```go
configManager := config.NewManager()
configManager.SetNetworkType(config.Mainnet)
bsvInstance := bsv.NewBSV(configManager)

// Coin type: 236
// RPC URL: https://api.whatsonchain.com/v1/bsv/main
// Explorer: https://whatsonchain.com
```

## 📚 Examples

The SDK includes comprehensive examples:

```bash
# Run the main example
go run cmd/main.go

# Run BIP44 HD wallet examples
go run examples/bip44/bip44_example.go

# Run enhanced features example
go run examples/enhanced/enhanced_usage.go

# Run sharding example
go run examples/sharding/sharding_example.go
```

## 🧪 Testing

```bash
# Run all tests
go test ./tests/ -v

# Run specific test
go test ./tests/ -run TestEnhancedWalletGeneration -v

# Run with coverage
go test ./tests/ -cover
```

## 🔁 CI
- GitHub Actions workflow runs on push/PR:
  - `go mod tidy` verification
  - `go fmt` verification
  - `golangci-lint` (if enabled in your environment)
  - `go test ./... -v`
- See `.github/workflows/ci.yml`.

## 📖 API Reference

### Core BSV Interface

```go
// Create BSV instance
bsvInstance := bsv.NewBSV(configManager)
bsvInstance := bsv.NewBSVWithNetwork(config.Testnet)

// Wallet generation
wallet, err := bsvInstance.GenerateWallet(mnemonic)
wallet, err := bsvInstance.GenerateWalletWithPath(mnemonic, account, change, addressIndex)
wallet, keypair, err := bsvInstance.GenerateWalletWithKeypair(mnemonic)

// BIP44 path management
defaultPath := bsvInstance.GetDefaultBIP44Path()
customPath := bsvInstance.GetBIP44Path(account, change, addressIndex)

// Balance and UTXO management
balance, err := bsvInstance.GetEnhancedBalance(address)
nativeBalance, err := bsvInstance.GetNativeBalance(address)
nonNativeBalance, err := bsvInstance.GetNonNativeBalance(address)
utxos, err := bsvInstance.GetUTXOs(address)

// Transaction operations
result, err := bsvInstance.BuildTransaction(txParams)
result, err := bsvInstance.SignAndSendTransaction(txParams)

// Configuration management
bsvInstance.SetNetworkType(config.Mainnet)
bsvInstance.UpdateUTXOConfig(utxoConfig)
bsvInstance.UpdateTransactionConfig(txConfig)
networkConfig := bsvInstance.GetNetworkConfig()
utxoConfig := bsvInstance.GetUTXOConfig()
txConfig := bsvInstance.GetTransactionConfig()

// Cache management
bsvInstance.ClearUTXOCache()
bsvInstance.ClearUTXOCacheForAddress(address)

// Address validation
err := bsvInstance.ValidateAddress(address)
```

### Configuration Types

```go
// Network configuration
type NetworkConfig struct {
    Name      string // "BSV Mainnet" or "BSV Testnet"
    IsTestnet bool
    CoinType  uint32 // 236 for mainnet, 1 for testnet
    RPCURL    string
    Explorer  string
}

// UTXO configuration
type UTXOConfig struct {
    IncludeNative      bool
    IncludeNonNative   bool
    MinConfirmations   int
    EnableCaching      bool
    CacheExpiry        int
    MaxUTXOsPerQuery   int
}

// Transaction configuration
type TransactionConfig struct {
    DefaultFeeRate     int
    MinFeeRate         int
    MaxFeeRate         int
    MaxTransactionSize int
    DustLimit          int
    EnableRBF          bool
}

// BIP44 path
type BIP44Path struct {
    Purpose      uint32
    CoinType     uint32
    Account      uint32
    Change       uint32
    AddressIndex uint32
}
```

### Data Types

```go
// Wallet result
type WalletResult struct {
    Address    string
    PrivateKey string
    PublicKey  string
}

// Enhanced balance information
type EnhancedBalanceInfo struct {
    Native    NativeBalanceInfo
    NonNative NonNativeBalanceInfo
}

type NativeBalanceInfo struct {
    Confirmed   int64
    Unconfirmed int64
    Total       int64
    UTXOCount   int
}

type NonNativeBalanceInfo struct {
    UTXOCount int
    Tokens    map[string]*TokenBalance
}

// UTXO information
type UTXO struct {
    TxID       string
    Vout       int
    Value      int64
    Address    string
    IsNative   bool
    TokenID    string
    TokenAmount int64
    Height     int
}

// Transaction parameters
type TransactionParams struct {
    From                 string
    To                   string
    Amount               int64
    PrivateKey           string
    IncludeNativeUTXOs   bool
    IncludeNonNativeUTXOs bool
    TokenTransfers       []*TokenTransfer
    DataOutputs          []*DataOutput
}

// Transaction result
type TransactionResult struct {
    TxID         string
    Fee          int64
    ExplorerURL  string
    InputsUsed   []*UTXO
    OutputsCreated []*TransactionOutput
    TokenTransfers []*TokenTransfer
    DataOutputs    []*DataOutput
}
```

## 🔧 Configuration Options

### UTXO Configuration
- `IncludeNative`: Include native BSV UTXOs
- `IncludeNonNative`: Include non-native token UTXOs
- `MinConfirmations`: Minimum confirmations required
- `EnableCaching`: Enable UTXO caching
- `CacheExpiry`: Cache expiration time in seconds
- `MaxUTXOsPerQuery`: Maximum UTXOs per API query

### Transaction Configuration
- `DefaultFeeRate`: Default fee rate in sat/vbyte
- `MinFeeRate`: Minimum fee rate in sat/vbyte
- `MaxFeeRate`: Maximum fee rate in sat/vbyte
- `MaxTransactionSize`: Maximum transaction size in bytes
- `DustLimit`: Dust limit in satoshis
- `EnableRBF`: Enable Replace-By-Fee

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.
