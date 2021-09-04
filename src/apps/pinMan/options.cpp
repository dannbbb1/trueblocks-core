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
 * Parts of this file were generated with makeClass --options. Edit only those parts of
 * the code outside of the BEG_CODE/END_CODE sections
 */
#define LOGGING_LEVEL_TEST
#include "options.h"

//---------------------------------------------------------------------------------------------------
static const COption params[] = {
    // BEG_CODE_OPTIONS
    // clang-format off
    COption("list", "l", "", OPT_SWITCH, "list the index and Bloom filter hashes from local manifest or pinning service"),  // NOLINT
    COption("init", "i", "", OPT_SWITCH, "initialize local index by downloading Bloom filters from pinning service"),
    COption("init_all", "n", "", OPT_SWITCH, "initialize local index by downloading both Bloom filters and index chunks"),  // NOLINT
    COption("sleep", "s", "<double>", OPT_HIDDEN | OPT_FLAG, "the number of seconds to sleep between requests during init (default .25)"),  // NOLINT
    COption("remote", "r", "", OPT_HIDDEN | OPT_SWITCH, "applicable only to --list mode, recover the manifest from pinning service"),  // NOLINT
    COption("pin_locally", "p", "", OPT_SWITCH, "pin all local files in the index to an IPFS store (requires IPFS)"),
    COption("", "", "", OPT_DESCRIPTION, "Manage pinned index of appearances and associated Bloom filters."),
    // clang-format on
    // END_CODE_OPTIONS
};
static const size_t nParams = sizeof(params) / sizeof(COption);

//---------------------------------------------------------------------------------------------------
bool COptions::parseArguments(string_q& command) {
    if (!standardOptions(command))
        return false;

    // BEG_CODE_LOCAL_INIT
    bool list = false;
    // END_CODE_LOCAL_INIT

    Init();
    explode(arguments, command, ' ');
    for (auto arg : arguments) {
        if (false) {
            // do nothing -- make auto code generation easier
            // BEG_CODE_AUTO
        } else if (arg == "-l" || arg == "--list") {
            list = true;

        } else if (arg == "-i" || arg == "--init") {
            init = true;

        } else if (arg == "-n" || arg == "--init_all") {
            init_all = true;

        } else if (startsWith(arg, "-s:") || startsWith(arg, "--sleep:")) {
            if (!confirmDouble("sleep", sleep, arg))
                return false;
        } else if (arg == "-s" || arg == "--sleep") {
            return flag_required("sleep");

        } else if (arg == "-r" || arg == "--remote") {
            remote = true;

        } else if (arg == "-p" || arg == "--pin_locally") {
            pin_locally = true;

        } else if (startsWith(arg, '-')) {  // do not collapse

            if (!builtInCmd(arg)) {
                return invalid_option(arg);
            }

            // END_CODE_AUTO
        }
    }

    // BEG_DEBUG_DISPLAY
    LOG_TEST_BOOL("list", list);
    LOG_TEST_BOOL("init", init);
    LOG_TEST_BOOL("init_all", init_all);
    LOG_TEST("sleep", sleep, (sleep == .25));
    LOG_TEST_BOOL("remote", remote);
    LOG_TEST_BOOL("pin_locally", pin_locally);
    // END_DEBUG_DISPLAY

    if (Mocked(""))
        return false;

    LOG_INFO("hashToIndexFormatFile:\t", cGreen, hashToIndexFormatFile, cOff);
    LOG_INFO("hashToBloomFormatFile:\t", cGreen, hashToBloomFormatFile, cOff);
    LOG_INFO("unchainedIndexAddr:\t", cGreen, unchainedIndexAddr, cOff);
    LOG_INFO("manifestHashEncoding:\t", cGreen, manifestHashEncoding, cOff);

    if (list && (init || init_all))
        return usage("Please choose only a single option.");

    if (!list && !init && !init_all)
        return usage("You must choose at least one of --list, --init, or --init_all.");

    if (init_all)
        init = true;

    if (pin_locally) {
        string_q res = doCommand("which ipfs");
        if (res.empty())
            return usage(
                "Could not find ipfs in your $PATH. You must install ipfs for the --pin_locally command to work.");
    }

    configureDisplay("pinMan", "CPinnedChunk", STR_DISPLAY_PINNEDCHUNK);

    return true;
}

//---------------------------------------------------------------------------------------------------
void COptions::Init(void) {
    registerOptions(nParams, params);

    // BEG_CODE_INIT
    init = false;
    init_all = false;
    sleep = .25;
    remote = false;
    pin_locally = false;
    // END_CODE_INIT
}

//---------------------------------------------------------------------------------------------------
COptions::COptions(void) {
    Init();

    // BEG_CODE_NOTES
    // clang-format off
    notes.push_back("One of `--list`, `--init`, or `--init_all` is required.");
    notes.push_back("the `--pin_locally` option only works if the IPFS executable is in your path.");
    // clang-format on
    // END_CODE_NOTES

    // BEG_ERROR_STRINGS
    // END_ERROR_STRINGS
}

//--------------------------------------------------------------------------------
COptions::~COptions(void) {
}
