// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"github.com/flokiorg/fcli/cli"
)

type TransactionsCommand struct {
	Handler *cli.WalletCliHandler
	Limit   int `long:"limit" description:"limit the number of rows to display"`
}

func (s *TransactionsCommand) Execute(args []string) error {

	s.Handler.RequireWallet()
	s.Handler.Transactions(s.Limit)
	return nil
}
