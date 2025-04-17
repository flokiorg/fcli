// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/flokiorg/fcli/cli"
	"github.com/flokiorg/fcli/utils"
)

type TransactionInput struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type BulkTransferCommand struct {
	Passphrase string    `short:"p" long:"passphrase" description:"Spending passphrase"`
	Amounts    []float64 `short:"a" long:"amount" description:"Amounts in FLC"`
	Addresses  []string  `short:"d" long:"address" description:"Destination addresses"`
	InputFile  string    `short:"f" long:"file" description:"JSON file with addresses and amounts"`

	Handler *cli.WalletCliHandler
}

func (s *BulkTransferCommand) Execute(args []string) error {

	electsrv := s.Handler.Config().ElectrumServer
	_, err := utils.ValidateAndNormalizeURI(electsrv, 50001)
	if err != nil {
		return fmt.Errorf("failed to validate electeum server address: %v", err)
	}

	// If a JSON file is provided, parse and use its values:
	if s.InputFile != "" {
		file, err := os.Open(s.InputFile)
		if err != nil {
			log.Fatalf("Cannot open input file: %v", err)
		}
		defer file.Close()

		var txInputs []TransactionInput
		if err := json.NewDecoder(file).Decode(&txInputs); err != nil {
			log.Fatalf("Error decoding JSON file: %v", err)
		}

		// Overwrite CLI fields from JSON content
		s.Addresses = nil
		s.Amounts = nil
		for _, input := range txInputs {
			s.Addresses = append(s.Addresses, input.Address)
			s.Amounts = append(s.Amounts, input.Amount)
		}
	}

	s.Handler.RequireWallet()
	s.Handler.BulkTransfer(s.Passphrase, s.Addresses, s.Amounts)
	return nil
}
