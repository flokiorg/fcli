// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"fmt"

	"github.com/flokiorg/fcli/cli"
	"github.com/flokiorg/fcli/utils"
)

type TransferCommand struct {
	Passphrase string  `short:"p" long:"passphrase" description:"Spending passphrase"`
	Amount     float64 `short:"a" long:"amount" description:"Amount in FLC"`
	Address    string  `short:"d" long:"address" description:"Destiation address"`

	Handler *cli.WalletCliHandler
}

func (s *TransferCommand) Execute(args []string) error {

	electsrv := s.Handler.Config().ElectrumServer
	_, err := utils.ValidateAndNormalizeURI(electsrv, 50001)
	if err != nil {
		return fmt.Errorf("failed to validate electeum server address: %v", err)
	}

	s.Handler.RequireWallet()
	s.Handler.Transfer(s.Passphrase, s.Address, s.Amount)
	return nil
}
