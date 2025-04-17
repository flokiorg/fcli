// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

type Config struct {
	// Create         bool          `short:"c" long:"create" description:"Create a new wallet"`
	// Restore        bool          `short:"r" long:"restore" description:"Restore a wallet using mnemonic"`
	// NewAddress     bool          `short:"n" long:"newaddress" description:"Create a new address (command)"`
	// ListAddrs      bool          `short:"a" long:"listaddrs" description:"List all available addresses."`
	// ShowXpub       bool          `short:"x" long:"showxpub" description:"Display the extended public key (command)"`
	// XpubBranch     uint32        `short:"i" long:"xpubbranch" description:"Xpub branch index"`
	// ListAccounts   bool          `short:"l" long:"listaccounts" description:"Display all accounts (command)"`
	// WithPrivate    bool          `long:"withprivate" description:"Include private data (e.g., xpriv) in the output. Use with caution."`
	// RegressionTest bool          `long:"regtest" description:"Use the regression test network"`
	// Testnet        bool          `long:"testnet" description:"Use the test network"`
	// PublicPassword string        `long:"pubpass" description:"Public password used to encrypt public data."`
	// DBTimeout      time.Duration `short:"t" long:"timeout" description:"Timeout duration (in seconds) for database connections."`
	// ElectrumServer string        `short:"e" long:"electserver" description:"Electrum server host:port"`
	// AccountID      uint32        `long:"id" description:"Account ID (default is used instead)"`
	// AccountName    string        `long:"name" description:"Account Name (default is used instead)"`
	// Balance bool `short:"b" description:"Show balance"`
	// reported
	// Import         bool          `long:"import" description:"Import command (default: 2147483647)"`
	// Sync           bool          `short:"s" description:"Sync command"`
	// Version bool `short:"v" description:"Print version"`
}
