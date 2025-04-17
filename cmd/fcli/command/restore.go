// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"github.com/flokiorg/fcli/cli"
)

type RestoreCommand struct {
	Handler *cli.WalletCliHandler
}

func (s *RestoreCommand) Execute(args []string) error {
	s.Handler.RestoreWallet()
	return nil
}
