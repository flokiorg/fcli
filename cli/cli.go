// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package cli

import (
	"fmt"
	"log"
	"strings"

	. "github.com/flokiorg/fcli/utils"
	"github.com/flokiorg/go-flokicoin/chaincfg"
	"github.com/flokiorg/go-flokicoin/chainutil/hdkeychain"
	"github.com/flokiorg/walletd/waddrmgr"
	"github.com/flokiorg/walletd/walletmgr"
	"github.com/flokiorg/walletd/walletseed/bip39"
)

var (
	defaultAddressScope = waddrmgr.KeyScopeBIP0044
)

type WalletCliHandler struct {
	*walletmgr.WalletService
	cfg     *Config
	network *chaincfg.Params
}

func NewWalletCliHandler(network *chaincfg.Params, cfg *Config) *WalletCliHandler {

	params := &walletmgr.WalletParams{
		Network:        network,
		Path:           cfg.WalletDir,
		Timeout:        cfg.DBTimeout,
		PublicPassword: cfg.PublicPassword,
		AddressScope:   defaultAddressScope,
		ElectrumServer: cfg.ElectrumServer,
		AccountID:      cfg.AccountID,
	}

	return &WalletCliHandler{
		WalletService: walletmgr.NewWalletService(params),
		network:       network,
		cfg:           cfg,
	}
}

func (wch *WalletCliHandler) CreateWallet() {
	exists, err := wch.WalletExists()
	if err != nil {
		log.Fatalf("unable to load wallet: %v", err)
	}
	if exists {
		log.Fatalf("A wallet already exists in the specified directory (%s)", wch.cfg.WalletDir)
	}

	privPass := ReadPassword("Enter a private password to secure your wallet: ", true)

	hex, words, err := wch.WalletService.Create(hdkeychain.RecommendedSeedLen, wch.cfg.AccountName, string(privPass))
	if err != nil {
		log.Fatalf("unable to create wallet: %v", err)
	}

	fmt.Println("\n========== Your Wallet Mnemonic ==========")
	fmt.Println("|                                         |")
	for i := 0; i < len(words); i += 4 {
		line := words[i:min(i+4, len(words))]
		fmt.Printf("| %-39s |\n", strings.Join(line, " "))
	}
	fmt.Println("|                                         |")
	fmt.Println("==========================================")
	fmt.Printf("  Hex: %s\n", hex)
	fmt.Println("==========================================")

	fmt.Println("\nKeep this mnemonic safe! If lost, you cannot recover your wallet.")

	fmt.Println("Wallet created successfully!")
}

func (wch *WalletCliHandler) RestoreWallet() {
	exists, err := wch.WalletExists()
	if err != nil {
		log.Fatalf("unable to load wallet: %v", err)
	}
	if exists {
		fmt.Println("A wallet already exists in the specified directory.")
		log.Fatalf("wallet restoration failed: %v", err)
	}

	mnemonic := ReadMnemonic()

	seed, err := bip39.EntropyFromMnemonic(mnemonic)
	if err != nil {
		log.Fatalf("Invalid mnemonic: %v", err)
	}

	privPass := ReadPassword("Enter a private password to secure your wallet: ", true)

	if err := wch.WalletService.RestoreWallet(seed, privPass, wch.cfg.AccountName); err != nil {
		log.Fatalf("wallet restoration failed: %v", err)
	}

	fmt.Println("Wallet restored successfully!")
}

func (wch *WalletCliHandler) RequireWallet() {

	// wallet commands
	exists, err := wch.WalletExists()
	if err != nil {
		log.Fatalf("unable to load wallet: %v", err)
	}
	if !exists {
		log.Fatal("wallet not found")
	}

	if err := wch.OpenWallet(); err != nil {
		log.Fatalf("opening failed: %v", err)
	}
}

func (wch *WalletCliHandler) Config() *Config {
	return wch.cfg
}

func (wch *WalletCliHandler) ProcessCommand() bool {

	// // wallet commands
	// exists, err := wch.WalletExists()
	// if err != nil {
	// 	log.Fatalf("unable to load wallet: %v", err)
	// }
	// if !exists {
	// 	log.Fatal("wallet not found")
	// }

	// if err := wch.OpenWallet(); err != nil {
	// 	log.Fatalf("opening failed: %v", err)
	// }

	// if wch.cfg.NewAddress {
	// 	wch.GenerateNewAddress()
	// } else if wch.cfg.ListAddrs {
	// 	wch.ListAddresses()
	// } else if wch.cfg.ShowXpub {
	// 	wch.ShowXpub(wch.cfg.XpubBranch)
	// } else if wch.cfg.ListAccounts {
	// 	wch.ListAccounts()
	// } else if wch.cfg.Balance {
	// 	wch.Balance()
	// } else {
	// 	wch.Dashboard(wch.ChainParams(), wch.cfg)
	// }

	// reported
	//  else if wch.cfg.Import {
	// wch.ImportWithWIF()
	// } else if wch.cfg.Sync {
	// 	wch.Sync()
	// }

	return true
}
