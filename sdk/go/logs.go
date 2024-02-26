// Copyright 2024 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package sdk

import (
	// EXISTING_CODE
	"io"
	"net/url"

	logs "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/sdk"
	// EXISTING_CODE
)

type LogsOptions struct {
	TransactionIds []string
	Emitter        []string // allow for ENS names and addresses
	Topic          []string // topics are strings
	Articulate     bool
	Globals

	// EXISTING_CODE
	// EXISTING_CODE
}

// Logs implements the chifra logs command for the SDK.
func (opts *LogsOptions) Logs(w io.Writer) error {
	values := make(url.Values)

	// EXISTING_CODE
	for _, v := range opts.TransactionIds {
		values.Add("transactions", v)
	}
	for _, v := range opts.Emitter {
		values.Add("emitter", v)
	}
	for _, v := range opts.Topic {
		values.Add("topic", v)
	}
	if opts.Articulate {
		values.Set("articulate", "true")
	}
	// EXISTING_CODE
	opts.Globals.mapGlobals(values)

	return logs.Logs(w, values)
}

// EXISTING_CODE
// EXISTING_CODE
