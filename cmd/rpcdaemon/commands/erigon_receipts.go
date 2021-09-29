package commands

import (
	"context"
	"fmt"

	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/core/types"
)

// GetLogsByHash implements erigon_getLogsByHash. Returns an array of arrays of logs generated by the transactions in the block given by the block's hash.
func (api *ErigonImpl) GetLogsByHash(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	chainConfig, err := api.chainConfig(tx)
	if err != nil {
		return nil, err
	}

	block, err := api.blockByHashWithSenders(tx, hash)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, nil
	}
	receipts, err := getReceipts(ctx, tx, chainConfig, block, block.Body().SendersFromTxs())
	if err != nil {
		return nil, fmt.Errorf("getReceipts error: %v", err)
	}

	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

// GetLogsByNumber implements erigon_getLogsByHash. Returns all the logs that appear in a block given the block's hash.
// func (api *ErigonImpl) GetLogsByNumber(ctx context.Context, number rpc.BlockNumber) ([][]*types.Log, error) {
// 	tx, err := api.db.Begin(ctx, false)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	number := rawdb.ReadHeaderNumber(tx, hash)
// 	if number == nil {
// 		return nil, fmt.Errorf("block not found: %x", hash)
// 	}

// 	receipts, err := getReceipts(ctx, tx, *number, hash)
// 	if err != nil {
// 		return nil, fmt.Errorf("getReceipts error: %v", err)
// 	}

// 	logs := make([][]*types.Log, len(receipts))
// 	for i, receipt := range receipts {
// 		logs[i] = receipt.Logs
// 	}
// 	return logs, nil
// }
