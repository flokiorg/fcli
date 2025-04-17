// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"fmt"

	"github.com/flokiorg/fcli/cli"
	. "github.com/flokiorg/fcli/utils"
)

type VersionCommand struct {
	Handler *cli.WalletCliHandler
}

func (s *VersionCommand) Execute(args []string) error {
	fmt.Println("Version:", Version)
	return nil
}
