// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package rpc

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc/query"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/ethereum/go-ethereum"
)

// GetTracesByBlockNumber returns a slice of traces in the given block
func (conn *Connection) GetTracesByBlockNumber(bn base.Blknum) ([]types.Trace, error) {
	if conn.StoreReadable() {
		traceGroup := &types.TraceGroup{
			BlockNumber:      bn,
			TransactionIndex: base.NOPOSN,
		}
		if err := conn.Store.Read(traceGroup, nil); err == nil {
			return traceGroup.Traces, nil
		}
	}

	method := "trace_block"
	params := query.Params{fmt.Sprintf("0x%x", bn)}

	if rawTraces, err := query.Query[[]types.RawTrace](conn.Chain, method, params); err != nil {
		return []types.Trace{}, err

	} else if rawTraces == nil || len(*rawTraces) == 0 {
		return []types.Trace{}, nil

	} else {
		curApp := types.Appearance{BlockNumber: uint32(^uint32(0))}
		curTs := conn.GetBlockTimestamp(bn)
		var traceIndex base.Tracenum

		// TODO: This could be loadTrace in the same way load Blocks works
		var ret []types.Trace
		for _, rawTrace := range *rawTraces {
			traceAction := types.TraceAction{
				Address:        base.HexToAddress(rawTrace.Action.Address),
				Author:         base.HexToAddress(rawTrace.Action.Author),
				Balance:        base.MustParseWei(rawTrace.Action.Balance),
				CallType:       rawTrace.Action.CallType,
				From:           base.HexToAddress(rawTrace.Action.From),
				Gas:            base.MustParseGas(rawTrace.Action.Gas),
				Init:           rawTrace.Action.Init,
				Input:          rawTrace.Action.Input,
				RefundAddress:  base.HexToAddress(rawTrace.Action.RefundAddress),
				RewardType:     rawTrace.Action.RewardType,
				SelfDestructed: base.HexToAddress(rawTrace.Action.SelfDestructed),
				To:             base.HexToAddress(rawTrace.Action.To),
				Value:          base.MustParseWei(rawTrace.Action.Value),
			}
			traceResult := types.TraceResult{}
			if rawTrace.Result != nil {
				traceResult.Address = base.HexToAddress(rawTrace.Result.Address)
				traceResult.Code = rawTrace.Result.Code
				traceResult.GasUsed = base.MustParseGas(rawTrace.Result.GasUsed)
				traceResult.Output = rawTrace.Result.Output
			}
			trace := types.Trace{
				Error:            rawTrace.Error,
				BlockHash:        base.HexToHash(rawTrace.BlockHash),
				BlockNumber:      rawTrace.BlockNumber,
				TransactionHash:  base.HexToHash(rawTrace.TransactionHash),
				TransactionIndex: rawTrace.TransactionIndex,
				TraceAddress:     rawTrace.TraceAddress,
				Subtraces:        rawTrace.Subtraces,
				TraceType:        rawTrace.TraceType,
				Timestamp:        curTs,
				Action:           &traceAction,
				Result:           &traceResult,
			}
			if trace.TransactionIndex != base.Txnum(curApp.TransactionIndex) {
				curApp = types.Appearance{
					BlockNumber:      uint32(trace.BlockNumber),
					TransactionIndex: uint32(trace.TransactionIndex),
				}
				traceIndex = 0
			}
			trace.TraceIndex = traceIndex
			traceIndex++
			trace.SetRaw(&rawTrace)
			ret = append(ret, trace)
		}

		if conn.StoreWritable() && conn.EnabledMap["traces"] && base.IsFinal(conn.LatestBlockTimestamp, curTs) {
			traceGroup := &types.TraceGroup{
				Traces:           ret,
				BlockNumber:      bn,
				TransactionIndex: base.NOPOSN,
			}
			_ = conn.Store.Write(traceGroup, nil)
		}
		return ret, nil
	}
}

// GetTracesByTransactionId returns a slice of traces in a given transaction
func (conn *Connection) GetTracesByTransactionId(bn base.Blknum, txid base.Txnum) ([]types.Trace, error) {
	if conn.StoreReadable() {
		traceGroup := &types.TraceGroup{
			BlockNumber:      bn,
			TransactionIndex: txid,
		}
		if err := conn.Store.Read(traceGroup, nil); err == nil {
			return traceGroup.Traces, nil
		}
	}

	tx, err := conn.GetTransactionByNumberAndId(bn, txid)
	if err != nil {
		return []types.Trace{}, err
	}

	return conn.GetTracesByTransactionHash(tx.Hash.Hex(), tx)
}

// GetTracesByTransactionHash returns a slice of traces in a given transaction's hash
func (conn *Connection) GetTracesByTransactionHash(txHash string, transaction *types.Transaction) ([]types.Trace, error) {
	if conn.StoreReadable() && transaction != nil {
		traceGroup := &types.TraceGroup{
			Traces:           make([]types.Trace, 0, len(transaction.Traces)),
			BlockNumber:      transaction.BlockNumber,
			TransactionIndex: transaction.TransactionIndex,
		}
		if err := conn.Store.Read(traceGroup, nil); err == nil {
			return traceGroup.Traces, nil
		}
	}

	method := "trace_transaction"
	params := query.Params{txHash}

	var ret []types.Trace
	if rawTraces, err := query.Query[[]types.RawTrace](conn.Chain, method, params); err != nil {
		return ret, ethereum.NotFound

	} else if rawTraces == nil || len(*rawTraces) == 0 {
		return []types.Trace{}, nil

	} else {
		curApp := types.Appearance{BlockNumber: uint32(^uint32(0))}
		var traceIndex base.Tracenum

		for _, rawTrace := range *rawTraces {
			value := base.NewWei(0)
			value.SetString(rawTrace.Action.Value, 0)
			balance := base.NewWei(0)
			balance.SetString(rawTrace.Action.Balance, 0)

			action := types.TraceAction{
				CallType:       rawTrace.Action.CallType,
				From:           base.HexToAddress(rawTrace.Action.From),
				Gas:            base.MustParseGas(rawTrace.Action.Gas),
				Input:          rawTrace.Action.Input,
				To:             base.HexToAddress(rawTrace.Action.To),
				Value:          *value,
				Balance:        *balance,
				Address:        base.HexToAddress(rawTrace.Action.Address),
				RefundAddress:  base.HexToAddress(rawTrace.Action.RefundAddress),
				SelfDestructed: base.HexToAddress(rawTrace.Action.SelfDestructed),
				Init:           rawTrace.Action.Init,
			}
			action.SetRaw(&rawTrace.Action)

			var result *types.TraceResult
			if rawTrace.Result != nil {
				result = &types.TraceResult{
					GasUsed: base.MustParseGas(rawTrace.Result.GasUsed),
					Output:  rawTrace.Result.Output,
					Code:    rawTrace.Result.Code,
				}
				if len(rawTrace.Result.Address) > 0 {
					result.Address = base.HexToAddress(rawTrace.Result.Address)
				}
				result.SetRaw(rawTrace.Result)
			}

			// TODO: This could be loadTrace in the same way load Blocks works
			trace := types.Trace{
				Error:            rawTrace.Error,
				BlockHash:        base.HexToHash(rawTrace.BlockHash),
				BlockNumber:      rawTrace.BlockNumber,
				TransactionHash:  base.HexToHash(rawTrace.TransactionHash),
				TransactionIndex: rawTrace.TransactionIndex,
				TraceAddress:     rawTrace.TraceAddress,
				Subtraces:        rawTrace.Subtraces,
				TraceType:        rawTrace.TraceType,
				Action:           &action,
				Result:           result,
			}
			if transaction != nil {
				trace.Timestamp = transaction.Timestamp
			}
			if trace.BlockNumber != base.Blknum(curApp.BlockNumber) || trace.TransactionIndex != base.Txnum(curApp.TransactionIndex) {
				curApp = types.Appearance{
					BlockNumber:      uint32(trace.BlockNumber),
					TransactionIndex: uint32(trace.TransactionIndex),
				}
				traceIndex = 0
			}
			trace.TraceIndex = traceIndex
			traceIndex++

			trace.SetRaw(&rawTrace)
			ret = append(ret, trace)
		}

		if transaction != nil {
			if conn.StoreWritable() && conn.EnabledMap["traces"] && base.IsFinal(conn.LatestBlockTimestamp, transaction.Timestamp) {
				traceGroup := &types.TraceGroup{
					Traces:           ret,
					BlockNumber:      transaction.BlockNumber,
					TransactionIndex: transaction.TransactionIndex,
				}
				_ = conn.Store.Write(traceGroup, nil)
			}
		}

		return ret, nil
	}
}

// GetTracesByFilter returns a slice of traces in a given transaction's hash
func (conn *Connection) GetTracesByFilter(filter string) ([]types.Trace, error) {
	var f types.TraceFilter
	ff, _ := f.ParseBangString(conn.Chain, filter)

	method := "trace_filter"
	params := query.Params{ff}

	var ret []types.Trace
	if rawTraces, err := query.Query[[]types.RawTrace](conn.Chain, method, params); err != nil {
		return ret, fmt.Errorf("trace filter %s returned an error: %w", filter, ethereum.NotFound)

	} else if rawTraces == nil || len(*rawTraces) == 0 {
		return []types.Trace{}, nil

	} else {
		curApp := types.Appearance{BlockNumber: uint32(^uint32(0))}
		curTs := conn.GetBlockTimestamp(f.FromBlock)
		var traceIndex base.Tracenum

		// TODO: This could be loadTrace in the same way load Blocks works
		for _, rawTrace := range *rawTraces {
			value := base.NewWei(0)
			value.SetString(rawTrace.Action.Value, 0)
			balance := base.NewWei(0)
			balance.SetString(rawTrace.Action.Balance, 0)

			action := types.TraceAction{
				CallType:       rawTrace.Action.CallType,
				From:           base.HexToAddress(rawTrace.Action.From),
				Gas:            base.MustParseGas(rawTrace.Action.Gas),
				Input:          rawTrace.Action.Input,
				To:             base.HexToAddress(rawTrace.Action.To),
				Value:          *value,
				Balance:        *balance,
				Address:        base.HexToAddress(rawTrace.Action.Address),
				RefundAddress:  base.HexToAddress(rawTrace.Action.RefundAddress),
				SelfDestructed: base.HexToAddress(rawTrace.Action.SelfDestructed),
				Init:           rawTrace.Action.Init,
			}
			action.SetRaw(&rawTrace.Action)

			var result *types.TraceResult
			if rawTrace.Result != nil {
				result = &types.TraceResult{
					GasUsed: base.MustParseGas(rawTrace.Result.GasUsed),
					Output:  rawTrace.Result.Output,
					Code:    rawTrace.Result.Code,
				}
				if len(rawTrace.Result.Address) > 0 {
					result.Address = base.HexToAddress(rawTrace.Result.Address)
				}
				result.SetRaw(rawTrace.Result)
			}

			trace := types.Trace{
				Error:            rawTrace.Error,
				BlockHash:        base.HexToHash(rawTrace.BlockHash),
				BlockNumber:      rawTrace.BlockNumber,
				TransactionHash:  base.HexToHash(rawTrace.TransactionHash),
				TransactionIndex: rawTrace.TransactionIndex,
				TraceAddress:     rawTrace.TraceAddress,
				Subtraces:        rawTrace.Subtraces,
				TraceType:        rawTrace.TraceType,
				Timestamp:        curTs,
				Action:           &action,
				Result:           result,
			}

			if trace.BlockNumber != base.Blknum(curApp.BlockNumber) {
				curTs = conn.GetBlockTimestamp(trace.BlockNumber)
			}

			if trace.BlockNumber != base.Blknum(curApp.BlockNumber) || trace.TransactionIndex != base.Txnum(curApp.TransactionIndex) {
				curApp = types.Appearance{
					BlockNumber:      uint32(trace.BlockNumber),
					TransactionIndex: uint32(trace.TransactionIndex),
				}
				traceIndex = 0
			}

			trace.TraceIndex = traceIndex
			traceIndex++

			trace.SetRaw(&rawTrace)
			ret = append(ret, trace)
		}
	}
	return ret, nil
}

// GetTracesCountInBlock returns the number of traces in a block
func (conn *Connection) GetTracesCountInBlock(bn base.Blknum) (uint64, error) {
	if traces, err := conn.GetTracesByBlockNumber(bn); err != nil {
		return base.NOPOS, err
	} else {
		return uint64(len(traces)), nil
	}
}
