package receiptsPkg

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
)

type ReceiptsOptionsType struct {
	Transactions []string
	Articulate   bool
	Globals      globals.GlobalOptionsType
}

var Options ReceiptsOptionsType

func (opts *ReceiptsOptionsType) TestLog() {
	logger.TestLog(len(opts.Transactions) > 0, "Transactions: ", opts.Transactions)
	logger.TestLog(opts.Articulate, "Articulate: ", opts.Articulate)
	opts.Globals.TestLog()
}

func FromRequest(w http.ResponseWriter, r *http.Request) *ReceiptsOptionsType {
	opts := &ReceiptsOptionsType{}
	for key, value := range r.URL.Query() {
		switch key {
		case "transactions":
			opts.Transactions = append(opts.Transactions, value...)
		case "articulate":
			opts.Articulate = true
		}
	}
	opts.Globals = *globals.FromRequest(w, r)

	return opts
}
