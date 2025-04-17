// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package cli

import "time"

type Config struct {
	WalletDir      string        `short:"w" long:"walletdir" description:"Directory for the wallet.db"`
	RegressionTest bool          `long:"regtest" description:"Use the regression test network"`
	Testnet        bool          `long:"testnet" description:"Use the test network"`
	PublicPassword string        `long:"pubpass" description:"Public password used to encrypt public data."`
	DBTimeout      time.Duration `short:"t" long:"timeout" description:"Timeout duration (in seconds) for wallet connection"`
	ElectrumServer string        `short:"e" long:"electserver" description:"Electrum server host:port"`
	AccountID      uint32        `long:"id" description:"Account ID (default '1' is used instead)"`
	AccountName    string        `long:"name" description:"Account Name (default 'myfloki' is used instead)"`
}
