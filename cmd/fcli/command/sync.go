// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"fmt"

	"github.com/flokiorg/fcli/cli"
	"github.com/flokiorg/fcli/utils"
)

type SyncCommand struct {
	Handler *cli.WalletCliHandler
}

func (s *SyncCommand) Execute(args []string) error {

	electsrv := s.Handler.Config().ElectrumServer
	_, err := utils.ValidateAndNormalizeURI(electsrv, 50001)
	if err != nil {
		return fmt.Errorf("failed to validate electeum server address: %v", err)
	}

	s.Handler.RequireWallet()
	s.Handler.Sync()
	return nil
}
