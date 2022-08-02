package scrapePkg

// TODO: BOGUS - MULTI-CHAIN SPECIFIC SLEEPING IS
// TODO: BOGUS - THINGS TO DO:
// TODO: BOGUS - CLEANUP BETWEEN RUNS IF THINGS DON'T WORK OUT AND SOMETHING HAD TO QUIT EARLY
// TODO: BOGUS - MAKE SURE WE'RE NOT RUNNING IF ACCTSCRAPE (OR ANYTHING ELSE?) IS RUNNINGs
// TODO: BOGUS - RESPOND WELL TO RUNNING AGAINST NON-ARCHIVE NODES see isArchiveNode() in C++ code
// TODO: BOGUS - RETRY ON FAILURE OF BLAZE - SEE NOTES BELOW
// TODO: BOGUS - THINGS TO DO:
// TODO: BOGUS - CLEANUP BETWEEN RUNS IF THINGS DON'T WORK OUT AND SOMETHING HAD TO QUIT EARLY
// TODO: BOGUS - MAKE SURE WE'RE NOT RUNNING IF ACCTSCRAPE (OR ANYTHING ELSE?) IS RUNNINGs
// TODO: BOGUS - RESPOND WELL TO RUNNING AGAINST NON-ARCHIVE NODES see isArchiveNode() in C++ code

// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
)

func (opts *ScrapeOptions) HandleScrape() error {
	progress, err := rpcClient.GetMetaData(opts.Globals.Chain, opts.Globals.TestMode)
	if err != nil {
		return err
	}

	if ok, err := opts.HandlePrepare(progress); !ok || err != nil {
		return err
	}

	for {
		blazeOpts := BlazeOptions{}
		if ok, err := opts.HandleScrapeBlaze(progress, &blazeOpts); !ok || err != nil {
			if !ok {
				break
			}
			goto PAUSE
		}

		logger.Enter("HandleScrapeConsolidate")
		if ok, err := opts.HandleScrapeConsolidate(progress, &blazeOpts); !ok || err != nil {
			logger.ExitError("HandleScrapeConsolidate", err)
			if !ok {
				break
			}
			goto PAUSE
		}
		logger.Exit("HandleScrapeConsolidate")

		// Clean up after this run of the blockScraper
		logger.Enter("HandleScrapePin")
		if ok, err := opts.HandleScrapePin(progress); !ok || err != nil {
			logger.ExitError("HandleScrapePin", err)
			if !ok {
				break
			}
			goto PAUSE
		}
		logger.Exit("HandleScrapePin")

	PAUSE:
		opts.Z_6_pause(progress)
	}

	return nil
}

// TODO: BOGUS - NOTES ON RE-RUN AFTER BLAZE FAILURE
// cleanFolder(indexFolder_ripe);
// static bool failed_already = false;
// if (!failed_already) {
// 	failed_already = true;
// 	LOG_WARN(cYellow, "Blaze quit without finishing. Retrying...", cOff);
// 	sleep(3);
// 	return scrape_blocks();
// }
// LOG_WARN(cYellow, "Blaze quit without finishing twice. Reprocessing...", cOff);

var Debugging = file.FileExists("./testing")
var ForceFail = file.FileExists("./forcefail")
