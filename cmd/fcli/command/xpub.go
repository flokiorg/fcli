// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"github.com/flokiorg/fcli/cli"
)

type XpubCommand struct {
	Index        uint32 `long:"index" short:"i" description:"branch index"`
	WithPrivate  bool   `long:"withprivate" description:"Include private data (e.g., xpriv) in the output. Use with caution."`
	PrintAddress bool   `long:"print-address" description:"Print address only"`

	Handler *cli.WalletCliHandler
}

func (s *XpubCommand) Execute(args []string) error {
	s.Handler.RequireWallet()
	s.Handler.ShowXpub(s.Index, s.WithPrivate, s.PrintAddress)
	return nil
}
