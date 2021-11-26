package whenPkg

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
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

// EXISTING_CODE
import (
	"net/http"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

func RunWhen(cmd *cobra.Command, args []string) error {
	opts := WhenFinishParse(args)

	err := opts.ValidateWhen()
	if err != nil {
		return err
	}

	// EXISTING_CODE
	if opts.List {
		err := opts.ListInternal()
		if err != nil {
			return err
		}
		if len(opts.Blocks) == 0 {
			return nil
		}
		// continue but don't show headers
		opts.List = false
		opts.Globals.NoHeader = true
	}

	return opts.Globals.PassItOn("whenBlock", opts.ToCmdLine())
	// EXISTING_CODE
}

func ServeWhen(w http.ResponseWriter, r *http.Request) bool {
	opts := FromRequest(w, r)

	err := opts.ValidateWhen()
	if err != nil {
		opts.Globals.RespondWithError(w, http.StatusInternalServerError, err)
		return true
	}

	// EXISTING_CODE
	if opts.List {
		err := opts.ListInternal()
		if err != nil {
			logger.Fatal("Cannot open local manifest file", err)
			return false
		}
		if len(opts.Blocks) == 0 {
			return true
		}
		// continue but don't show headers or --list
		r.URL.RawQuery = strings.Replace(r.URL.RawQuery, "list", "noop", -1)
		r.URL.RawQuery += "&no_header"
	}
	return false
	// EXISTING_CODE
}

// EXISTING_CODE
// EXISTING_CODE