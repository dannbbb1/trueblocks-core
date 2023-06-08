package blocksPkg

import (
	"context"
	"errors"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum"
)

func (opts *BlocksOptions) HandleUniq() (err error) {
	chain := opts.Globals.Chain

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawAppearance], errorChan chan error) {
		procFunc := func(s *types.SimpleAppearance) error {
			modelChan <- s
			return nil
		}

		for _, br := range opts.BlockIds {
			blockNums, err := br.ResolveBlocks(opts.Globals.Chain)
			if err != nil {
				errorChan <- err
				if errors.Is(err, ethereum.NotFound) {
					continue
				}
				cancel()
				return
			}

			for _, bn := range blockNums {
				addrMap := make(index.AddressBooleanMap)
				ts := rpc.GetBlockTimestamp(chain, bn)
				if err := opts.ProcessBlockUniqs(chain, procFunc, bn, addrMap, ts); err != nil {
					errorChan <- err
					if errors.Is(err, ethereum.NotFound) {
						continue
					}
					cancel()
					return
				}
			}
		}
	}

	extra := map[string]interface{}{
		"uniq": true,
	}
	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOptsWithExtra(extra))
}

func (opts *BlocksOptions) ProcessBlockUniqs(chain string, procFunc index.UniqProcFunc, bn uint64, addrMap index.AddressBooleanMap, ts int64) error {
	if bn == 0 {
		if namesArray, err := names.LoadNamesArray(opts.Globals.Chain, names.Prefund, names.SortByAddress, []string{}); err != nil {
			return err
		} else {
			for i, name := range namesArray {
				address := name.Address.Hex()
				index.StreamAppearance(procFunc, opts.Flow, "genesis", address, bn, uint64(i), utils.NOPOS, ts, addrMap)
			}
		}

	} else {
		if block, err := rpcClient.GetBlockByNumberWithTxs(chain, bn, true /* isFinal doesn't matter */); err != nil {
			return err
		} else {
			miner := block.Miner.Hex()
			txid := uint64(99999)
			if block.Miner.IsZero() {
				// Early clients allowed misconfigured miner settings with address 0x0 (reward got
				// burned). We enter a false record with a false tx_id to account for this.
				miner = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
				txid = 99997
			}
			index.StreamAppearance(procFunc, opts.Flow, "miner", miner, bn, txid, utils.NOPOS, ts, addrMap)

			if uncles, err := rpcClient.GetUnclesByNumber(chain, bn); err != nil {
				return err
			} else {
				for _, uncle := range uncles {
					unc := uncle.Miner.Hex()
					txid = uint64(99998)
					if uncle.Miner.IsZero() {
						// Early clients allowed misconfigured miner settings with address 0x0 (reward got
						// burned). We enter a false record with a false tx_id to account for this.
						unc = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
						txid = 99998 // do not change this!
					}
					index.StreamAppearance(procFunc, opts.Flow, "uncle", unc, bn, txid, utils.NOPOS, ts, addrMap)
				}
			}

			for _, trans := range block.Transactions {
				if trans.Traces, err = rpcClient.GetTracesByTransactionId(opts.Globals.Chain, trans.BlockNumber, trans.TransactionIndex); err != nil {
					return err
				}
				if err = index.UniqFromTransDetails(chain, procFunc, opts.Flow, &trans, ts, addrMap); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
