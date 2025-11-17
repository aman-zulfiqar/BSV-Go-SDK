package types

import (
	"errors"
	"math/big"
)

// Common error definitions
var (
	ErrInvalidMnemonic    = errors.New("invalid mnemonic phrase")
	ErrInvalidShard       = errors.New("invalid shard format")
	ErrInsufficientShards = errors.New("insufficient shards to reconstruct mnemonic")
	ErrInvalidAddress     = errors.New("invalid BSV address")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrNetworkError       = errors.New("network error")
	ErrInvalidPrivateKey  = errors.New("invalid private key")
	ErrInvalidUTXO        = errors.New("invalid UTXO")
	ErrTransactionFailed  = errors.New("transaction failed")
)

// WalletResult represents a generated wallet
type WalletResult struct {
	Address    string `json:"address"`    // BSV address
	PrivateKey string `json:"privateKey"` // WIF private key
	PublicKey  string `json:"publicKey"`  // Public key in hex
}

// ShardingResult represents the result of mnemonic sharding
type ShardingResult struct {
	Shards      []string `json:"shards"`      // Array of shard strings
	Threshold   int      `json:"threshold"`   // Minimum shards needed (2)
	TotalShares int      `json:"totalShares"` // Total number of shards (3)
}

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID          string `json:"txid"`          // Transaction ID
	Vout          uint32 `json:"vout"`          // Output index
	Value         int64  `json:"value"`         // Value in satoshis
	ScriptPubKey  string `json:"scriptPubKey"`  // Script public key
	Address       string `json:"address"`       // Address (for convenience)
	Confirmations int    `json:"confirmations"` // Number of confirmations
	Height        int    `json:"height"`        // Block height
	IsNative      bool   `json:"isNative"`      // Whether this is native BSV UTXO
	TokenID       string `json:"tokenId"`       // Token ID for non-native UTXOs (empty for native)
	TokenAmount   int64  `json:"tokenAmount"`   // Token amount for non-native UTXOs
}

// TransactionParams represents parameters for building a transaction
type TransactionParams struct {
	From       string `json:"from"`       // Sender address
	To         string `json:"to"`         // Recipient address
	Amount     int64  `json:"amount"`     // Amount in satoshis
	FeeRate    int64  `json:"feeRate"`    // Fee rate in satoshis per vbyte (optional)
	PrivateKey string `json:"privateKey"` // Private key (WIF or mnemonic)
	// Enhanced parameters for native/non-native support
	IncludeNativeUTXOs    bool             `json:"includeNativeUTXOs"`    // Include native BSV UTXOs
	IncludeNonNativeUTXOs bool             `json:"includeNonNativeUTXOs"` // Include non-native token UTXOs
	TokenTransfers        []*TokenTransfer `json:"tokenTransfers"`        // Token transfers for non-native transactions
	DataOutputs           []*DataOutput    `json:"dataOutputs"`           // Data outputs (OP_RETURN)
}

// TokenTransfer represents a token transfer in a transaction
type TokenTransfer struct {
	TokenID string `json:"tokenId"` // Token identifier
	To      string `json:"to"`      // Recipient address
	Amount  int64  `json:"amount"`  // Token amount to transfer
}

// DataOutput represents a data output (OP_RETURN) in a transaction
type DataOutput struct {
	Data string `json:"data"` // Hex-encoded data to include in OP_RETURN
}

// TransactionResult represents the result of a signed transaction
type TransactionResult struct {
	SignedTx       string               `json:"signedTx"`       // Signed transaction in hex
	TxID           string               `json:"txId"`           // Transaction ID
	Fee            int64                `json:"fee"`            // Transaction fee in satoshis
	Change         int64                `json:"change"`         // Change amount in satoshis
	ExplorerURL    string               `json:"explorerUrl"`    // Explorer URL for the transaction
	InputsUsed     []*UTXO              `json:"inputsUsed"`     // UTXOs used as inputs
	OutputsCreated []*TransactionOutput `json:"outputsCreated"` // Outputs created
	TokenTransfers []*TokenTransfer     `json:"tokenTransfers"` // Token transfers executed
	DataOutputs    []*DataOutput        `json:"dataOutputs"`    // Data outputs included
}

// TransactionOutput represents an output in a transaction
type TransactionOutput struct {
	Address      string `json:"address"`      // Output address
	Amount       int64  `json:"amount"`       // Amount in satoshis
	ScriptPubKey string `json:"scriptPubKey"` // Script public key
	IsData       bool   `json:"isData"`       // Whether this is a data output
	Data         string `json:"data"`         // Data content (for data outputs)
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	Name        string `json:"name"`        // Network name
	RPCURL      string `json:"rpcUrl"`      // RPC endpoint URL
	ExplorerURL string `json:"explorerUrl"` // Explorer URL
	IsTestnet   bool   `json:"isTestnet"`   // Whether this is testnet
	ChainID     string `json:"chainId"`     // Chain identifier
}

// BalanceInfo represents balance information
type BalanceInfo struct {
	Confirmed   int64 `json:"confirmed"`   // Confirmed balance in satoshis
	Unconfirmed int64 `json:"unconfirmed"` // Unconfirmed balance in satoshis
	Total       int64 `json:"total"`       // Total balance in satoshis
}

// EnhancedBalanceInfo represents detailed balance information including native and non-native
type EnhancedBalanceInfo struct {
	Native    *NativeBalanceInfo    `json:"native"`    // Native BSV balance
	NonNative *NonNativeBalanceInfo `json:"nonNative"` // Non-native token balances
	Total     int64                 `json:"total"`     // Total native BSV balance
}

// NativeBalanceInfo represents native BSV balance information
type NativeBalanceInfo struct {
	Confirmed   int64 `json:"confirmed"`   // Confirmed native balance
	Unconfirmed int64 `json:"unconfirmed"` // Unconfirmed native balance
	Total       int64 `json:"total"`       // Total native balance
	UTXOCount   int   `json:"utxoCount"`   // Number of native UTXOs
}

// NonNativeBalanceInfo represents non-native token balance information
type NonNativeBalanceInfo struct {
	Tokens    map[string]*TokenBalance `json:"tokens"`    // Token balances by token ID
	UTXOCount int                      `json:"utxoCount"` // Number of non-native UTXOs
}

// TokenBalance represents balance for a specific token
type TokenBalance struct {
	TokenID     string `json:"tokenId"`     // Token identifier
	Confirmed   int64  `json:"confirmed"`   // Confirmed token balance
	Unconfirmed int64  `json:"unconfirmed"` // Unconfirmed token balance
	Total       int64  `json:"total"`       // Total token balance
	UTXOCount   int    `json:"utxoCount"`   // Number of UTXOs for this token
}

// Helper function to convert satoshis to BSV
func SatoshisToBSV(satoshis int64) *big.Float {
	// 1 BSV = 100,000,000 satoshis
	bsv := new(big.Float).SetInt64(satoshis)
	bsv.Quo(bsv, big.NewFloat(100000000))
	return bsv
}

// Helper function to convert BSV to satoshis
func BSVToSatoshis(bsv *big.Float) int64 {
	// 1 BSV = 100,000,000 satoshis
	bsv.Mul(bsv, big.NewFloat(100000000))
	satoshis, _ := bsv.Int64()
	return satoshis
}

// Helper function to format BSV amount
func FormatBSV(satoshis int64) string {
	bsv := SatoshisToBSV(satoshis)
	return bsv.Text('f', 8) // 8 decimal places
}
