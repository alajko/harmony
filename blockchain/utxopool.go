package blockchain

import (
	"encoding/hex"
)

// UTXOPool stores transactions and balance associated with each address.
type UTXOPool struct {
	// Mapping from address to a map of transaction id to a map of the index of output
	// array in that transaction to that balance.
	utxo map[string]map[string]int
}

// VerifyTransactions verifies if a list of transactions valid.
func (utxopool *UTXOPool) VerifyTransactions(transactions []*Transaction) bool {
	spentTXOs := make(map[string]map[string]bool)
	if utxopool != nil {
		for _, tx := range transactions {
			inTotal := 0
			// Calculate the sum of TxInput
			for _, in := range tx.TxInput {
				inTxID := hex.EncodeToString(in.TxID)
				// Check if the transaction with the addres is spent or not.
				if val, ok := spentTXOs[in.Address][inTxID]; ok {
					if val {
						return false
					}
				}
				// Mark the transactions with the addres spent.
				if _, ok := spentTXOs[in.Address]; ok {
					spentTXOs[in.Address][inTxID] = true
				} else {
					spentTXOs[in.Address] = make(map[string]bool)
					spentTXOs[in.Address][inTxID] = true
				}

				// Sum the balance up to the inTotal.
				if val, ok := utxopool.utxo[in.Address][inTxID]; ok {
					inTotal += val
				} else {
					return false
				}
			}

			outTotal := 0
			// Calculate the sum of TxOutput
			for _, out := range tx.TxOutput {
				outTotal += out.Value
			}
			if inTotal != outTotal {
				return false
			}
		}
	}
	return true
}

// VerifyAndUpdate verifies a list of transactions and update utxoPool.
func (utxopool *UTXOPool) VerifyAndUpdate(transactions []*Transaction) bool {
	if utxopool.VerifyTransactions(transactions) {
		utxopool.Update(transactions)
		return true
	}
	return false
}

// Update utxo balances with a list of new transactions.
func (utxopool *UTXOPool) Update(transactions []*Transaction) {
	if utxopool != nil {
		for _, tx := range transactions {
			curTxID := hex.EncodeToString(tx.ID)

			// Remove
			for _, in := range tx.TxInput {
				inTxID := hex.EncodeToString(in.TxID)
				delete(utxopool.utxo[in.Address], inTxID)
			}

			// Update
			for _, out := range tx.TxOutput {
				if _, ok := utxopool.utxo[out.Address]; ok {
					utxopool.utxo[out.Address][curTxID] = out.Value
				} else {
					utxopool.utxo[out.Address] = make(map[string]int)
					utxopool.utxo[out.Address][curTxID] = out.Value
				}
			}
		}
	}
}

// CreateUTXOPoolFromTransaction a utxo pool from a genesis transaction.
func CreateUTXOPoolFromTransaction(tx *Transaction) *UTXOPool {
	var utxoPool UTXOPool
	txID := hex.EncodeToString(tx.ID)
	utxoPool.utxo = make(map[string]map[string]int)
	for _, out := range tx.TxOutput {
		utxoPool.utxo[out.Address] = make(map[string]int)
		utxoPool.utxo[out.Address][txID] = DefaultCoinbaseValue
	}
	return &utxoPool
}

// CreateUTXOPoolFromGenesisBlockChain a utxo pool from a genesis blockchain.
func CreateUTXOPoolFromGenesisBlockChain(bc *Blockchain) *UTXOPool {
	tx := bc.blocks[0].Transactions[0]
	var utxoPool UTXOPool
	txID := hex.EncodeToString(tx.ID)
	utxoPool.utxo = make(map[string]map[string]int)
	for _, out := range tx.TxOutput {
		utxoPool.utxo[out.Address] = make(map[string]int)
		utxoPool.utxo[out.Address][txID] = DefaultCoinbaseValue
	}
	return &utxoPool
}