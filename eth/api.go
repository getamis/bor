// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthereumAPI provides an API to access Ethereum full node-related information.
type EthereumAPI struct {
	e *Ethereum
}

// NewEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewEthereumAPI(e *Ethereum) *EthereumAPI {
	return &EthereumAPI{e}
}

// Etherbase is the address that mining rewards will be sent to.
func (api *EthereumAPI) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}

// Coinbase is the address that mining rewards will be sent to (alias for Etherbase).
func (api *EthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate.
func (api *EthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().Hashrate())
}

// Mining returns an indication if this node is currently mining.
func (api *EthereumAPI) Mining() bool {
	return api.e.IsMining()
}

func getFinalizedBlockNumber(eth *Ethereum) (uint64, error) {
	currentBlockNum := eth.BlockChain().CurrentBlock()

	doExist, number, hash := eth.Downloader().GetWhitelistedMilestone()
	if doExist && number <= currentBlockNum.Number.Uint64() {
		block := eth.BlockChain().GetBlockByNumber(number)

		if block.Hash() == hash {
			return number, nil
		}
	}

	doExist, number, hash = eth.Downloader().GetWhitelistedCheckpoint()
	if doExist && number <= currentBlockNum.Number.Uint64() {
		block := eth.BlockChain().GetBlockByNumber(number)

		if block.Hash() == hash {
			return number, nil
		}
	}

	return 0, fmt.Errorf("No finalized block")
}

// GetBlockReceipts returns all transaction receipts of the specified block.
func (api *PrivateDebugAPI) GetBlockReceipts(ctx context.Context, blockHash common.Hash) ([]map[string]interface{}, error) {
	block := api.eth.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, errors.New("block not found")
	}

	blockNumber := block.Header().Number
	receipts := api.eth.blockchain.GetReceiptsByHash(blockHash)

	if receipts == nil {
		return nil, errors.New("receipts not found")
	}

	txs := block.Transactions()
	if len(txs) != len(receipts) {
		return nil, fmt.Errorf("txs length doesn't equal to receipts' length")
	}

	txReceipts := make([]map[string]interface{}, 0, len(txs))

	for idx, receipt := range receipts {
		tx := txs[idx]
		fields, err := ethapi.ToTransactionReceipt(ctx, api.eth.APIBackend, tx, receipt, blockHash, tx.Hash(), blockNumber.Uint64(), uint64(idx))

		if err != nil {
			return nil, err
		}

		txReceipts = append(txReceipts, fields)
	}

	receipt := rawdb.ReadBorReceipt(api.eth.chainDb, blockHash, blockNumber.Uint64())
	if receipt != nil {
		tx, _, _, _ := rawdb.ReadBorTransaction(api.eth.chainDb, receipt.TxHash)
		fields, err := ethapi.ToTransactionReceipt(ctx, api.eth.APIBackend, tx, receipt, blockHash, receipt.TxHash, blockNumber.Uint64(), uint64(receipt.TransactionIndex))

		if err != nil {
			return nil, err
		}

		txReceipts = append(txReceipts, fields)
	}

	return txReceipts, nil
}
