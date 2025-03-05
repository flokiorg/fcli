// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"log"
	"os"
	"time"

	"github.com/flokiorg/fcli/cli"
	. "github.com/flokiorg/fcli/common"
	"github.com/flokiorg/go-flokicoin/chaincfg"
	"github.com/flokiorg/go-flokicoin/chainutil"
	"github.com/flokiorg/walletd/walletdb/bdb"
	"github.com/flokiorg/walletd/walletseed/bip39"
	"github.com/flokiorg/walletd/walletseed/bip39/wordlists"
	flags "github.com/jessevdk/go-flags"
)

var (
	defaultDBTimeout = 10 * time.Second
	defaultPubPass   = "/flc/public"
	defaultWordList  = wordlists.English
	network          = &chaincfg.MainNetParams
	defaultAppName   = "flcwallet"

	defaultAccountID   uint32 = 1
	defaultAccountName string = "myfloki"
)

func main() {
	var cfg Config

	parser := flags.NewParser(&cfg, flags.Default|flags.PassDoubleDash)

	if len(os.Args) == 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	// Register the backend database
	bdb.Register()

	// init word list
	bip39.SetWordList(defaultWordList)

	if cfg.PublicPassword == "" {
		cfg.PublicPassword = defaultPubPass
	}

	if cfg.RegressionTest {
		network = &chaincfg.RegressionNetParams
	} else if cfg.Testnet {
		network = &chaincfg.TestNet3Params
	}

	if opt := parser.FindOptionByShortName('t'); !optionDefined(opt) {
		cfg.DBTimeout = defaultDBTimeout
	}

	if opt := parser.FindOptionByShortName('e'); !optionDefined(opt) {
		log.Fatal("electserver is required")
	}

	if opt := parser.FindOptionByShortName('d'); !optionDefined(opt) {
		cfg.Dir = chainutil.AppDataDir(defaultAppName, false)
	}

	if opt := parser.FindOptionByLongName("id"); !optionDefined(opt) {
		cfg.AccountID = defaultAccountID
	}

	if opt := parser.FindOptionByLongName("name"); !optionDefined(opt) {
		cfg.AccountName = defaultAccountName
	}

	handler := cli.NewWalletCliHandler(network, &cfg)

	if cfg.Create {
		handler.CreateWallet()
	} else if cfg.Restore {
		handler.RestoreWallet()
	} else if ok := handler.ProcessCommand(); !ok {
		parser.WriteHelp(os.Stdout)
	}
}

func optionDefined(opt *flags.Option) bool {
	return opt != nil && opt.IsSet()
}
