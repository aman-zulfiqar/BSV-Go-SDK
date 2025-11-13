# BSV Custodial SDK - API Reference

## Overview

The BSV Custodial SDK provides a comprehensive Go library for Bitcoin SV (BSV) wallet management, mnemonic sharding, and transaction handling.

## Package Structure

```
pkg/
├── bsv/                    # Main BSV interface
│   ├── wallet/            # Wallet derivation helpers
│   ├── transaction/       # Transaction builder
│   └── utxo/              # UTXO management
├── mnemonic/              # Mnemonic generation and validation
├── sharding/              # Shamir Secret Sharing
├── types/                 # Common types and interfaces
└── config/                # Runtime configuration
```

## Core Types

### WalletResult
```go
type WalletResult struct {
    Address    string `json:"address"`    // BSV address
    PrivateKey string `json:"privateKey"` // WIF private key
    PublicKey  string `json:"publicKey"`  // Public key in hex
}
```

### ShardingResult
```go
type ShardingResult struct {
    Shards      []string `json:"shards"`      // Array of shard strings
    Threshold   int      `json:"threshold"`   // Minimum shards needed
    TotalShares int      `json:"totalShares"` // Total number of shards
}
```

### TransactionParams
```go
type TransactionParams struct {
    From       string `json:"from"`       // Sender address
    To         string `json:"to"`         // Recipient address
    Amount     int64  `json:"amount"`     // Amount in satoshis
    FeeRate    int64  `json:"feeRate"`    // Fee rate in sat/vbyte
    PrivateKey string `json:"privateKey"` // Private key (WIF or mnemonic)
}
```

### TransactionResult
```go
type TransactionResult struct {
    SignedTx    string `json:"signedTx"`    // Signed transaction in hex
    TxID        string `json:"txId"`        // Transaction ID
    Fee         int64  `json:"fee"`         // Transaction fee in satoshis
    Change      int64  `json:"change"`      // Change amount in satoshis
    ExplorerURL string `json:"explorerUrl"` // Explorer URL
}
```

## Mnemonic Management

### Generate Mnemonic
```go
// Generate 12-word mnemonic
mnemonic, err := mnemonic.Generate(mnemonic.Strength128)

// Generate 24-word mnemonic
mnemonic, err := mnemonic.Generate(mnemonic.Strength256)
```

### Validate Mnemonic
```go
err := mnemonic.Validate(mnemonicPhrase)
```

### Get Word Count
```go
wordCount := mnemonic.GetWordCount(mnemonicPhrase)
```

### Normalize Mnemonic
```go
normalized := mnemonic.Normalize(mnemonicPhrase)
```

## Sharding Operations

### Split Mnemonic into Shards
```go
// Create 2/3 threshold shards
result, err := sharding.SplitMnemonic(mnemonic, 2, 3)

// Create 3/5 threshold shards
result, err := sharding.SplitMnemonic(mnemonic, 3, 5)
```

### Combine Shards to Reconstruct Mnemonic
```go
reconstructed, err := sharding.CombineShards([]string{shard1, shard2})
```

### Validate Shard
```go
isValid := sharding.ValidateShard(shardString)
```

## BSV Wallet Operations

### Generate Wallet from Mnemonic
```go
wallet, err := bsv.GenerateWallet(mnemonic, isTestnet)
```

### Generate Wallet with Keypair
```go
wallet, keyPair, err := bsv.GenerateWalletWithKeypair(mnemonic, isTestnet)
```

### Validate Address
```go
err := bsv.ValidateAddress(address, isTestnet)
```

### Get Balance
```go
balance, err := bsv.GetBalance(address, isTestnet)
```

### Get UTXOs
```go
utxos, err := bsv.GetUTXOs(address, isTestnet)
```

## Transaction Operations

### Build Transaction
```go
params := &types.TransactionParams{
    From:       senderAddress,
    To:         recipientAddress,
    Amount:     100000, // 100000 satoshis
    FeeRate:    5,      // 5 sat/vbyte
    PrivateKey: mnemonic, // or WIF private key
}

result, err := bsv.SignAndSendTransaction(params, isTestnet)
```

### Send Transaction with WIF Private Key
```go
params := &types.TransactionParams{
    From:       senderAddress,
    To:         recipientAddress,
    Amount:     50000,
    FeeRate:    10,
    PrivateKey: "WIF_PRIVATE_KEY_HERE",
}

result, err := bsv.SignAndSendTransaction(params, isTestnet)
```

## Utility Functions

### Convert Satoshis to BSV
```go
bsv := types.SatoshisToBSV(100000000) // Returns *big.Float
```

### Convert BSV to Satoshis
```go
satoshis := types.BSVToSatoshis(big.NewFloat(1.0)) // Returns int64
```

### Format BSV Amount
```go
formatted := types.FormatBSV(100000000) // Returns "1.00000000"
```

## Error Handling

The SDK defines common errors for consistent error handling:

```go
var (
    ErrInvalidMnemonic     = errors.New("invalid mnemonic phrase")
    ErrInvalidShard        = errors.New("invalid shard format")
    ErrInsufficientShards  = errors.New("insufficient shards to reconstruct mnemonic")
    ErrInvalidAddress      = errors.New("invalid BSV address")
    ErrInsufficientFunds   = errors.New("insufficient funds")
    ErrInvalidAmount       = errors.New("invalid amount")
    ErrNetworkError        = errors.New("network error")
    ErrInvalidPrivateKey   = errors.New("invalid private key")
    ErrInvalidUTXO         = errors.New("invalid UTXO")
    ErrTransactionFailed   = errors.New("transaction failed")
)
```

## Network Configuration

### Testnet
- **Network**: BSV Testnet
- **Explorer**: https://test.whatsonchain.com
- **API**: https://api.whatsonchain.com/v1/bsv/test
- **Faucet**: https://faucet.bitcoincloud.net/

### Mainnet
- **Network**: BSV Mainnet
- **Explorer**: https://whatsonchain.com
- **API**: https://api.whatsonchain.com/v1/bsv/main

## Security Considerations

### Mnemonic Security
- Always validate mnemonics before use
- Store mnemonics securely
- Use strong entropy for generation

### Sharding Security
- Never store all shards in the same location
- Use encrypted storage for shards
- Validate shards before reconstruction
- Consider using different storage mediums

### Transaction Security
- Validate addresses before sending
- Double-check amounts before signing
- Use appropriate fee rates
- Test on testnet before mainnet use

## Best Practices

### Development
1. Always use testnet for development
2. Implement proper error handling
3. Validate all inputs
4. Use appropriate fee rates
5. Test thoroughly before production

### Production
1. Use mainnet with caution
2. Implement rate limiting
3. Monitor transaction fees
4. Use hardware wallets for large amounts
5. Implement proper backup strategies

## Examples

See the `examples/` directory for comprehensive usage examples:

- `basic_usage.go` - Basic SDK functionality
- `sharding_example.go` - Advanced sharding operations
- `transaction_example.go` - Transaction handling

## Testing

Run tests with:
```bash
make test
```

Individual test files:
- `tests/mnemonic_test.go` - Mnemonic operations
- `tests/sharding_test.go` - Sharding operations

## Building and Running

```bash
# Download dependencies
make deps

# Build the library
make build

# Run tests
make test

# Run examples
make run-basic
make run-sharding
make run-transaction

# Interactive example
make run-example
```
