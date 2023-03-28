// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package chunksPkg

import (
	"context"
	"fmt"
	"io"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (opts *ChunksOptions) HandleAddresses(blockNums []uint64) error {
	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawIndexAddress], errorChan chan error) {
		showAddresses := func(walker *index.IndexWalker, path string, first bool) (bool, error) {
			if path != cache.ToBloomPath(path) {
				return false, fmt.Errorf("should not happen in showFinalizedStats")
			}

			path = cache.ToIndexPath(path)
			indexChunk, err := index.NewChunkData(path)
			if err != nil {
				return false, err
			}
			defer indexChunk.Close()

			_, err = indexChunk.File.Seek(int64(index.HeaderWidth), io.SeekStart)
			if err != nil {
				return false, err
			}

			cnt := 0
			for i := 0; i < int(indexChunk.Header.AddressCount); i++ {
				if opts.Globals.TestMode && i > walker.MaxTests() {
					continue
				}

				obj := index.AddressRecord{}
				err := obj.ReadAddress(indexChunk.File)
				if err != nil {
					return false, err
				}

				s := types.SimpleIndexAddress{
					Address: hexutil.Encode(obj.Address.Bytes()),
					Range:   indexChunk.Range.String(),
					Offset:  obj.Offset,
					Count:   obj.Count,
				}

				modelChan <- &s
				cnt++
			}
			return true, nil
		}

		walker := index.NewIndexWalker(
			opts.Globals.Chain,
			opts.Globals.TestMode,
			10, /* maxTests */
			showAddresses,
		)
		if err := walker.WalkBloomFilters(blockNums); err != nil {
			errorChan <- err
			cancel()
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
