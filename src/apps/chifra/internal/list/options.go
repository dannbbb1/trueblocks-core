package listPkg

/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

import (
	"net/http"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/globals"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type ListOptionsType struct {
	Addrs       []string
	Count       bool
	Appearances bool
	FirstBlock  uint64
	LastBlock   uint64
	Globals     globals.GlobalOptionsType
}

var Options ListOptionsType

func (opts *ListOptionsType) TestLog() {
	logger.TestLog(len(opts.Addrs) > 0, "Addrs: ", opts.Addrs)
	logger.TestLog(opts.Count, "Count: ", opts.Count)
	logger.TestLog(opts.Appearances, "Appearances: ", opts.Appearances)
	logger.TestLog(opts.FirstBlock != 0, "FirstBlock: ", opts.FirstBlock)
	logger.TestLog(opts.LastBlock != utils.NOPOS, "LastBlock: ", opts.LastBlock)
	opts.Globals.TestLog()
}

func FromRequest(w http.ResponseWriter, r *http.Request) *ListOptionsType {
	opts := &ListOptionsType{}
	for key, value := range r.URL.Query() {
		switch key {
		case "addrs":
			opts.Addrs = append(opts.Addrs, value...)
		case "count":
			opts.Count = true
		case "appearances":
			opts.Appearances = true
		case "firstblock":
			opts.FirstBlock = globals.ToUint64(value[0])
		case "lastblock":
			opts.LastBlock = globals.ToUint64(value[0])
		}
	}
	opts.Globals = *globals.FromRequest(w, r)

	return opts
}
