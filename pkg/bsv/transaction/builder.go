package transaction

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/muhammadamman/BSV-Go/pkg/bsv/utxo"
	"github.com/muhammadamman/BSV-Go/pkg/bsv/wallet"
	"github.com/muhammadamman/BSV-Go/pkg/config"
	"github.com/muhammadamman/BSV-Go/pkg/mnemonic"
	"github.com/muhammadamman/BSV-Go/pkg/types"
)

// Builder handles BSV transaction building with dynamic configuration
type Builder struct {
	configManager *config.Manager
	utxoManager   *utxo.Manager
	httpClient    *http.Client
}

// NewBuilder creates a new transaction builder
func NewBuilder(configManager *config.Manager) *Builder {
	return &Builder{
		configManager: configManager,
		utxoManager:   utxo.NewManager(configManager),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// BuildTransaction builds a BSV transaction with enhanced native/non-native support
func (b *Builder) BuildTransaction(params *types.TransactionParams) (*wire.MsgTx, error) {
	// Validate inputs
	if err := b.validateParams(params); err != nil {
		return nil, err
	}

	// Get sender address and keypair
	senderAddress, keyPair, err := b.getSenderInfo(params.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender info: %v", err)
	}

	// Validate sender address matches
	if senderAddress != params.From {
		return nil, fmt.Errorf("sender address mismatch: expected %s, got %s", params.From, senderAddress)
	}

	// Select UTXOs based on transaction type
	var selectedUTXOs []types.UTXO
	var fee int64

	txConfig := b.configManager.GetTransactionConfig()

	if len(params.TokenTransfers) > 0 {
		// Token transfer transaction
		selectedUTXOs, fee, err = b.selectUTXOsForTokenTransfer(params)
		if err != nil {
			return nil, fmt.Errorf("failed to select UTXOs for token transfer: %v", err)
		}
	} else {
		// Regular BSV transaction
		if params.FeeRate <= 0 {
			params.FeeRate = txConfig.DefaultFeeRate
		}
		selectedUTXOs, fee, err = b.utxoManager.SelectUTXOs(params.From, params.Amount, params.FeeRate)
		if err != nil {
			return nil, fmt.Errorf("failed to select UTXOs: %v", err)
		}
	}

	// Create new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs
	for _, utxo := range selectedUTXOs {
		txHash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, fmt.Errorf("invalid UTXO transaction hash: %v", err)
		}

		prevOut := wire.NewOutPoint(txHash, utxo.Vout)
		txIn := wire.NewTxIn(prevOut, nil, nil)
		tx.AddTxIn(txIn)
	}

	// Add outputs
	err = b.addOutputs(tx, params, selectedUTXOs, fee)
	if err != nil {
		return nil, fmt.Errorf("failed to add outputs: %v", err)
	}

	// Sign the transaction
	if err := b.signTransaction(tx, selectedUTXOs, keyPair); err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	return tx, nil
}

// SignAndSendTransaction builds, signs, and broadcasts a transaction
func (b *Builder) SignAndSendTransaction(params *types.TransactionParams) (*types.TransactionResult, error) {
	// Build the transaction
	tx, err := b.BuildTransaction(params)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction: %v", err)
	}

	// Serialize the transaction
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	// Get transaction ID
	txID := tx.TxHash().String()

	// Broadcast the transaction
	if err := b.broadcastTransaction(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %v", err)
	}

	// Calculate detailed transaction information
	result, err := b.calculateTransactionResult(tx, params, txID, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to calculate transaction result: %v", err)
	}

	return result, nil
}

// GetEnhancedBalance retrieves enhanced balance for an address
func (b *Builder) GetEnhancedBalance(address string) (*types.EnhancedBalanceInfo, error) {
	return b.utxoManager.GetEnhancedBalance(address)
}

// GetNativeBalance retrieves native balance for an address
func (b *Builder) GetNativeBalance(address string) (*types.NativeBalanceInfo, error) {
	return b.utxoManager.GetNativeBalance(address)
}

// GetNonNativeBalance retrieves non-native balance for an address
func (b *Builder) GetNonNativeBalance(address string) (*types.NonNativeBalanceInfo, error) {
	return b.utxoManager.GetNonNativeBalance(address)
}

// GetBalance retrieves balance for an address (backward compatibility)
func (b *Builder) GetBalance(address string) (int64, error) {
	return b.utxoManager.GetConfirmedBalance(address)
}

// GetUTXOs retrieves UTXOs for an address
func (b *Builder) GetUTXOs(address string) ([]types.UTXO, error) {
	return b.utxoManager.GetUTXOs(address)
}

// Helper methods

func (b *Builder) validateParams(params *types.TransactionParams) error {
	txConfig := b.configManager.GetTransactionConfig()

	if params.From == "" {
		return fmt.Errorf("sender address is required")
	}
	if params.To == "" {
		return fmt.Errorf("recipient address is required")
	}
	if params.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if params.PrivateKey == "" {
		return fmt.Errorf("private key is required")
	}

	// Validate fee rate
	if params.FeeRate <= 0 {
		params.FeeRate = txConfig.DefaultFeeRate
	} else if params.FeeRate < txConfig.MinFeeRate {
		return fmt.Errorf("fee rate %d is below minimum %d", params.FeeRate, txConfig.MinFeeRate)
	} else if params.FeeRate > txConfig.MaxFeeRate {
		return fmt.Errorf("fee rate %d exceeds maximum %d", params.FeeRate, txConfig.MaxFeeRate)
	}

	// Validate token transfers
	for i, transfer := range params.TokenTransfers {
		if transfer.TokenID == "" {
			return fmt.Errorf("token transfer %d: token ID is required", i)
		}
		if transfer.To == "" {
			return fmt.Errorf("token transfer %d: recipient address is required", i)
		}
		if transfer.Amount <= 0 {
			return fmt.Errorf("token transfer %d: amount must be positive", i)
		}
	}

	return nil
}

func (b *Builder) getSenderInfo(privateKey string) (string, *wallet.KeyPair, error) {
	networkConfig := b.configManager.GetNetworkConfig()

	// Check if it's a mnemonic (12 or more words)
	words := strings.Fields(strings.TrimSpace(privateKey))
	if len(words) >= 12 {
		// It's a mnemonic - validate and generate wallet
		if err := mnemonic.Validate(privateKey); err != nil {
			return "", nil, fmt.Errorf("invalid mnemonic: %v", err)
		}

		walletResult, keyPair, err := wallet.GenerateWalletWithKeypair(privateKey, networkConfig.IsTestnet)
		if err != nil {
			return "", nil, fmt.Errorf("failed to generate wallet from mnemonic: %v", err)
		}

		return walletResult.Address, keyPair, nil
	} else {
		// It's a WIF private key
		var network *chaincfg.Params
		if networkConfig.IsTestnet {
			network = &chaincfg.TestNet3Params
		} else {
			network = &chaincfg.MainNetParams
		}

		wif, err := btcutil.DecodeWIF(privateKey)
		if err != nil {
			return "", nil, fmt.Errorf("invalid WIF private key: %v", err)
		}

		// Validate network
		if !wif.IsForNet(network) {
			return "", nil, fmt.Errorf("WIF private key is not for the correct network")
		}

		// Get address from public key
		pubKey := wif.PrivKey.PubKey()
		address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), network)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create address from public key: %v", err)
		}

		// Create keypair
		keyPair := &wallet.KeyPair{
			PrivateKey: wif.PrivKey,
			PublicKey:  pubKey,
			Network:    network,
		}

		return address.EncodeAddress(), keyPair, nil
	}
}

func (b *Builder) selectUTXOsForTokenTransfer(params *types.TransactionParams) ([]types.UTXO, int64, error) {
	// For now, we'll select UTXOs for the first token transfer
	// In a more sophisticated implementation, you might want to handle multiple token transfers
	if len(params.TokenTransfers) == 0 {
		return nil, 0, fmt.Errorf("no token transfers specified")
	}

	firstTransfer := params.TokenTransfers[0]
	return b.utxoManager.SelectUTXOsForTokenTransfer(params.From, firstTransfer.TokenID, firstTransfer.Amount, params.FeeRate)
}

func (b *Builder) addOutputs(tx *wire.MsgTx, params *types.TransactionParams, selectedUTXOs []types.UTXO, fee int64) error {
	networkConfig := b.configManager.GetNetworkConfig()

	var network *chaincfg.Params
	if networkConfig.IsTestnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}

	// Add recipient output for BSV
	recipientAddr, err := btcutil.DecodeAddress(params.To, network)
	if err != nil {
		return fmt.Errorf("invalid recipient address: %v", err)
	}

	recipientScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		return fmt.Errorf("failed to create recipient script: %v", err)
	}

	tx.AddTxOut(wire.NewTxOut(params.Amount, recipientScript))

	// Add token transfer outputs
	for _, transfer := range params.TokenTransfers {
		// For token transfers, we would typically add special outputs
		// This is a simplified implementation - in reality, you'd need to implement
		// the specific token protocol (e.g., SLP, STAS, etc.)

		// For now, we'll create a data output with token transfer information
		tokenData := fmt.Sprintf("TOKEN_TRANSFER:%s:%s:%d", transfer.TokenID, transfer.To, transfer.Amount)
		tokenDataHex := hex.EncodeToString([]byte(tokenData))

		// Create OP_RETURN script
		opReturnScript, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).
			AddData([]byte(tokenDataHex)).
			Script()
		if err != nil {
			return fmt.Errorf("failed to create token transfer script: %v", err)
		}

		tx.AddTxOut(wire.NewTxOut(0, opReturnScript)) // 0 value for OP_RETURN
	}

	// Add data outputs
	for _, dataOutput := range params.DataOutputs {
		data, err := hex.DecodeString(dataOutput.Data)
		if err != nil {
			return fmt.Errorf("invalid data output hex: %v", err)
		}

		opReturnScript, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_RETURN).
			AddData(data).
			Script()
		if err != nil {
			return fmt.Errorf("failed to create data output script: %v", err)
		}

		tx.AddTxOut(wire.NewTxOut(0, opReturnScript)) // 0 value for OP_RETURN
	}

	// Add change output if necessary
	change, hasChange := b.utxoManager.CalculateChange(selectedUTXOs, params.Amount, fee)
	if hasChange {
		senderAddr, err := btcutil.DecodeAddress(params.From, network)
		if err != nil {
			return fmt.Errorf("invalid sender address: %v", err)
		}

		changeScript, err := txscript.PayToAddrScript(senderAddr)
		if err != nil {
			return fmt.Errorf("failed to create change script: %v", err)
		}

		tx.AddTxOut(wire.NewTxOut(change, changeScript))
	}

	return nil
}

func (b *Builder) signTransaction(tx *wire.MsgTx, utxos []types.UTXO, keyPair *wallet.KeyPair) error {
	networkConfig := b.configManager.GetNetworkConfig()

	var network *chaincfg.Params
	if networkConfig.IsTestnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}

	for i, utxo := range utxos {
		// Create the script to sign
		senderAddr, err := btcutil.DecodeAddress(utxo.Address, network)
		if err != nil {
			return fmt.Errorf("failed to decode address: %v", err)
		}

		script, err := txscript.PayToAddrScript(senderAddr)
		if err != nil {
			return fmt.Errorf("failed to create script: %v", err)
		}

		// Create signature script
		sigScript, err := txscript.SignatureScript(tx, i, script, txscript.SigHashAll, keyPair.PrivateKey, true)
		if err != nil {
			return fmt.Errorf("failed to create signature script: %v", err)
		}

		tx.TxIn[i].SignatureScript = sigScript
	}

	return nil
}

func (b *Builder) broadcastTransaction(txBytes []byte) error {
	networkConfig := b.configManager.GetNetworkConfig()
	url := networkConfig.RPCURL + "/tx/raw"

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(txBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("User-Agent", "BSV-Enhanced-SDK/1.0.0")

	// Send request
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("broadcast failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ClearUTXOCache clears the UTXO cache
func (b *Builder) ClearUTXOCache() {
	b.utxoManager.ClearCache()
}

// ClearUTXOCacheForAddress clears UTXO cache for a specific address
func (b *Builder) ClearUTXOCacheForAddress(address string) {
	b.utxoManager.ClearCacheForAddress(address)
}

func (b *Builder) calculateTransactionResult(tx *wire.MsgTx, params *types.TransactionParams, txID string, txBytes []byte) (*types.TransactionResult, error) {
	networkConfig := b.configManager.GetNetworkConfig()

	// Get UTXOs used as inputs
	var inputsUsed []*types.UTXO
	selectedUTXOs, _, err := b.utxoManager.SelectUTXOs(params.From, params.Amount, params.FeeRate)
	if err == nil {
		for _, utxo := range selectedUTXOs {
			inputsUsed = append(inputsUsed, &utxo)
		}
	}

	// Calculate outputs created
	var outputsCreated []*types.TransactionOutput
	for _, txOut := range tx.TxOut {
		output := &types.TransactionOutput{
			Amount:       txOut.Value,
			ScriptPubKey: hex.EncodeToString(txOut.PkScript),
			IsData:       false, // Would need to check if it's OP_RETURN
		}
		outputsCreated = append(outputsCreated, output)
	}

	// Calculate fee and change
	var totalInput int64
	for _, utxo := range selectedUTXOs {
		totalInput += utxo.Value
	}

	fee := totalInput - params.Amount
	var change int64
	if len(tx.TxOut) > 1 {
		change = tx.TxOut[1].Value
		fee = totalInput - params.Amount - change
	}

	// Create explorer URL
	explorerURL := fmt.Sprintf("%s/tx/%s", networkConfig.ExplorerURL, txID)

	return &types.TransactionResult{
		SignedTx:       hex.EncodeToString(txBytes),
		TxID:           txID,
		Fee:            fee,
		Change:         change,
		ExplorerURL:    explorerURL,
		InputsUsed:     inputsUsed,
		OutputsCreated: outputsCreated,
		TokenTransfers: params.TokenTransfers,
		DataOutputs:    params.DataOutputs,
	}, nil
}
