// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"os"
	"time"

	"github.com/flokiorg/fcli/cli"
	"github.com/flokiorg/fcli/cmd/fcli/command"
	"github.com/flokiorg/go-flokicoin/chaincfg"
	"github.com/flokiorg/go-flokicoin/chainutil"
	"github.com/flokiorg/walletd/walletdb/bdb"
	"github.com/flokiorg/walletd/walletseed/bip39"
	"github.com/flokiorg/walletd/walletseed/bip39/wordlists"
	flags "github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {

	var cfg cli.Config

	parser := flags.NewParser(&cfg, flags.IgnoreUnknown)
	_, err := parser.ParseArgs(os.Args[1:])
	if err != nil {
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

	if opt := parser.FindOptionByShortName('w'); !optionDefined(opt) {
		cfg.WalletDir = chainutil.AppDataDir(defaultAppName, false)
	}

	if opt := parser.FindOptionByLongName("id"); !optionDefined(opt) {
		cfg.AccountID = defaultAccountID
	}

	if opt := parser.FindOptionByLongName("name"); !optionDefined(opt) {
		cfg.AccountName = defaultAccountName
	}

	handler := cli.NewWalletCliHandler(network, &cfg)

	parser = flags.NewParser(&cfg, flags.Default|flags.PassDoubleDash)

	parser.AddCommand("create", "Create new wallet", "", &command.CreateCommand{Handler: handler})
	parser.AddCommand("restore", "Restore wallet", "", &command.RestoreCommand{Handler: handler})

	parser.AddCommand("addresses", "List existing addresses", "", &command.ListAddressesCommand{Handler: handler})
	parser.AddCommand("newaddress", "Create new address", "", &command.NewAddressCommand{Handler: handler})
	parser.AddCommand("xpub", "Print extended public key (xpub)", "", &command.XpubCommand{Handler: handler})
	parser.AddCommand("accounts", "Display all accounts", "", &command.ListAccountsCommand{Handler: handler})
	parser.AddCommand("balance", "Print wallet balance", "", &command.BalanceCommand{Handler: handler})
	parser.AddCommand("transactions", "Print wallet transactions", "", &command.TransactionsCommand{Handler: handler})
	parser.AddCommand("sync", "Sync with network", "", &command.SyncCommand{Handler: handler})
	parser.AddCommand("transfer", "Send transaction", "", &command.TransferCommand{Handler: handler})
	parser.AddCommand("bulktransfer", "Send transaction", "", &command.BulkTransferCommand{Handler: handler})
	parser.AddCommand("version", "Show version", "", &command.VersionCommand{})

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

}

func optionDefined(opt *flags.Option) bool {
	return opt != nil && opt.IsSet()
}
