// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	. "github.com/flokiorg/fcli/utils"
	"github.com/flokiorg/go-flokicoin/chainjson"
	"github.com/flokiorg/go-flokicoin/chainutil"
	"github.com/flokiorg/walletd/chain/electrum"
	"github.com/flokiorg/walletd/waddrmgr"
	"github.com/flokiorg/walletd/wallet"
)

func (wch *WalletCliHandler) GenerateNewAddress() {
	addr, err := wch.CreateNewAddress()
	if err != nil {
		log.Fatalf("failed creating address: %v", err)
	}

	log.Printf("Address: %s", addr.EncodeAddress())
}

func (wch *WalletCliHandler) ListAddresses() {
	addrs, err := wch.AccountAddresses(wch.cfg.AccountID)
	if err != nil {
		log.Fatalf("Failed to list addresses: %v", err)
	}

	if len(addrs) == 0 {
		log.Printf("No addresses found. Your wallet does not contain any generated addresses yet.")
	} else {
		log.Printf("List of available addresses:")
		for _, addr := range addrs {
			log.Printf("- %s", addr)
		}
	}
}

func (wch *WalletCliHandler) ListAccounts() {
	accounts, err := wch.Accounts(defaultAddressScope)
	if err != nil {
		log.Fatalf("unable to fetch accounts: %v", err)
	}

	log.Printf("Height: %d", accounts.CurrentBlockHeight)
	log.Printf("Blockhash: %v", accounts.CurrentBlockHash.String())

	log.Printf("Accounts:")
	for _, acc := range accounts.Accounts {
		log.Printf(" - %s", acc.AccountName)
		log.Printf("   ID: %v", acc.AccountNumber)
		log.Printf("   Balance: %v", acc.TotalBalance)
		addrs, err := wch.AccountAddresses(acc.AccountNumber)
		if err != nil {
			log.Fatalf("Failed to list addresses: %v", err)
		}

		if len(addrs) == 0 {
			log.Printf("   No addresses found")
		} else {
			log.Printf("   list of available addresses (%d):", len(addrs))
			for i, addr := range addrs {
				if i >= 5 {
					log.Printf("    ...and %d more", len(addrs)-i)
					break
				}
				log.Printf("    - %s", addr)
			}
		}
	}
	fmt.Print("==============================\n\n")

}

func (wch *WalletCliHandler) ImportWithWIF() {
	strWif, err := ReadLine("Enter WIF: ", func(s string) error {
		_, err := chainutil.DecodeWIF(s)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("unable to decode wif: %v", err)
	}

	wif, _ := chainutil.DecodeWIF(strWif)

	privPass := ReadPassword("Enter the private password to unlock the wallet: ", false)
	if err := wch.Unlock(privPass, nil); err != nil {
		log.Fatalf("Failed to unlock wallet: %v", err)
	}
	defer wch.Lock()
	fmt.Println("importing...")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		n1, n2, n3, _, err := wch.Watch()
		for {
			select {
			case <-n1:
				return
			case <-n2:
				// fmt.Printf("n2")
			case <-n3:
				// fmt.Printf("n3")
			case <-err:
			}
		}
	}()
	addr, err := wch.Wallet.ImportPrivateKey(defaultAddressScope, wif, nil, true)
	if err != nil {
		log.Fatalf("importation failed: %v", err)
	}
	log.Printf("address imported: %v", addr)
	wg.Wait()
}

func (wch *WalletCliHandler) ShowXpub(branch uint32, withPrivateData bool, printAddress bool) {

	if withPrivateData {
		privPass := ReadPassword("Enter the private password to unlock the wallet: ", false)
		if err := wch.Unlock(privPass, nil); err != nil {
			log.Fatalf("Failed to unlock wallet: %v", err)
		}
		defer wch.Lock()
	}

	props, err := wch.AccountProperties(defaultAddressScope, wch.cfg.AccountID)
	if err != nil {
		log.Fatalf("Failed to get account props: %v", err)
	}

	if props.AccountPubKey == nil {
		log.Printf("unable to retrieve derivation data for this account: %d\n", wch.cfg.AccountID)
		return
	}

	branchKey, err := props.AccountPubKey.Derive(waddrmgr.ExternalBranch)
	if err != nil {
		log.Fatalf("Failed deriving branch: %v", err)
	}

	xpub := branchKey.String()

	var xpriv string
	if withPrivateData {
		xpriv = props.AccountPrivKey.String()
	}

	if branch == 0 {
		branch = uint32(time.Since(wch.ChainParams().GenesisBlock.Header.Timestamp).Minutes()) // elapsed minutes
	}

	address, priv, err := DeriveKeysFromXpub(wch.ChainParams(), xpriv, xpub, branch)
	if err != nil {
		log.Fatalf("Failed deriving add/wif: %v", err)
	}

	if !printAddress {
		if withPrivateData {
			log.Printf("Your xpriv key: %s\n", xpriv)
			log.Printf("Derived WIF (index %d): %s\n", branch, priv)
		}
		log.Printf("Your xpub key: %s\n", xpub)
		log.Printf("Derived public key (index %d): %s", branch, address)
	} else {
		fmt.Print(address)
	}
}

func (wch *WalletCliHandler) Transactions(limit int) {
	if err := printAllTransactionHistory(wch.Wallet, limit); err != nil {
		log.Fatalf("failed to fetch transactions: %v", err)
	}
}

func (wch *WalletCliHandler) Dashboard() {

	if _, err := ValidateAndNormalizeURI(wch.cfg.ElectrumServer, 80); err != nil {
		log.Fatalf("invalid electrum server URL: %v", err)
	}

	if _, err := wch.Synchronize(); err != nil {
		log.Fatalf("unable to sync: %v", err)
	}

	log.Printf("Network: %s", wch.ChainParams().Name)
	log.Printf("Account ID: %d", wch.cfg.AccountID)

	log.Printf("Address type: %s", StrAddrType(waddrmgr.ScopeAddrMap[defaultAddressScope].ExternalAddrType))
	log.Printf("Address path: %s", defaultAddressScope.String())

	wch.ListAddresses()
	wch.ShowXpub(1, false, false)
	wch.Balance()

	printAllTransactionHistory(wch.Wallet, -1)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	accountNotif := wch.NtfnServer.AccountNotifications()
	accountSpentnessNotif := wch.NtfnServer.AccountSpentnessNotifications(wch.cfg.AccountID)
	txtNotif := wch.NtfnServer.TransactionNotifications()

	log.Println("waiting for new updates...")
	for {
		select {

		case n := <-accountNotif.C:
			if n == nil {
				continue
			}
			log.Println("account updated")
			wch.Balance()

		case n := <-accountSpentnessNotif.C:
			if n == nil {
				continue
			}
			log.Println("got a new UTXO")
			wch.Balance()

		case n := <-txtNotif.C:
			if n == nil {
				continue
			}
			log.Println("got new transaction")
			wch.Balance()

		case <-sigChan:
			return
		}
	}

}

func (wch *WalletCliHandler) Balance() {
	balance, err := wch.CalculateAccountBalances(wch.cfg.AccountID, 0)
	if err != nil {
		log.Printf("unable to fetch balance: %v\n", err)
		return
	}
	log.Printf("balance: %f", balance.Total.ToFLC())
}

func (wch *WalletCliHandler) currentState() (int32, string) {
	accounts, err := wch.Accounts(defaultAddressScope)
	if err != nil {
		log.Fatalf("unable to fetch accounts: %v", err)
	}

	return accounts.CurrentBlockHeight, accounts.CurrentBlockHash.String()
}

func (wch *WalletCliHandler) Sync() {

	startupBlock, err := wch.Synchronize()
	if err != nil {
		log.Fatalf("unable to sync: %v", err)
	}

	IsUpToDate := func() bool {

		var height int32
		wBlock, err := wch.WalletService.CurrentWalletBlock()
		if err == nil {
			height = wBlock.Height
		}

		return height >= startupBlock.Height
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		n1, n2, n3, n4, err := wch.Watch()
		for {
			select {
			case <-n1:
				if IsUpToDate() {
					return
				}
			case <-n2:
				if IsUpToDate() {
					return
				}
			case <-n3:
				if IsUpToDate() {
					return
				}
			case <-n4:
				if IsUpToDate() {
					return
				}
			case e := <-err:
				if !errors.Is(e, electrum.NerrHealthPong) {
					log.Fatalf("unexpected error: %v", e)
				}
			}

		}
	}()

	log.Println("Syncing...")
	wg.Wait()
	wch.ListAccounts()
}

func (wch *WalletCliHandler) Transfer(password string, strAddress string, inAmount float64) {

	amount, err := chainutil.NewAmount(inAmount)
	if err != nil {
		log.Fatalf("invalid amount: %v", err)
	}

	address, err := chainutil.DecodeAddress(strAddress, wch.network)
	if err != nil {
		log.Fatalf("invalid address: %v", err)
	}

	var privPass []byte
	if len(password) == 0 {
		privPass = ReadPassword("Enter the private password to unlock the wallet: ", false)
	} else {
		privPass = []byte(password)
	}

	if _, err = wch.Synchronize(); err != nil {
		log.Fatalf("unable to sync with electrum: %v", err)
	}

	_, ts, _, ns, _ := wch.Watch()
	select {
	case <-ts:
	case <-ns:
	}

	tx, err := wch.SimpleTransfer(privPass, address, amount, 0)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	fmt.Printf("%s", tx.TxHash())

}

func (wch *WalletCliHandler) BulkTransfer(password string, strAddresses []string, inAmounts []float64) {

	if len(strAddresses) != len(inAmounts) {
		log.Fatalf("addresses (%d) != amounts (%d)", len(strAddresses), len(inAmounts))
	}

	amounts := make([]chainutil.Amount, 0, len(inAmounts))
	addresses := make([]chainutil.Address, 0, len(strAddresses))

	for _, inAmount := range inAmounts {
		amount, err := chainutil.NewAmount(inAmount)
		if err != nil {
			log.Fatalf("invalid amount: %v", err)
		}
		amounts = append(amounts, amount)
	}

	for _, strAddress := range strAddresses {
		address, err := chainutil.DecodeAddress(strAddress, wch.network)
		if err != nil {
			log.Fatalf("invalid address: %v", err)
		}
		addresses = append(addresses, address)
	}

	var privPass []byte
	if len(password) == 0 {
		privPass = ReadPassword("Enter the private password to unlock the wallet: ", false)
	} else {
		privPass = []byte(password)
	}

	if _, err := wch.Synchronize(); err != nil {
		log.Fatalf("unable to sync with electrum: %v", err)
	}

	_, ts, _, ns, _ := wch.Watch()
	select {
	case <-ts:
	case <-ns:
	}

	tx, err := wch.BulkSimpleTransfer(privPass, addresses, amounts, 0)
	if err != nil {
		log.Fatalf("failed: %v", err)
	}

	fmt.Printf("%s", tx.TxHash())

}

func printAddressTransactionHistory(w *wallet.Wallet, address string) error {

	pkHashes := map[string]struct{}{}
	pkHashes[address] = struct{}{}

	txDetails, err := w.ListAddressTransactions(pkHashes)
	if err != nil {
		return fmt.Errorf("failed to list transactions: %w", err)
	}
	cleanHistory := aggregateTransactions(txDetails)

	// Print transaction details
	fmt.Printf("Address Transactions\n--------------------\n")
	printTransactionHistory(cleanHistory)
	fmt.Println()

	return nil
}

func printAllTransactionHistory(w *wallet.Wallet, limit int) error {
	// Fetch the transaction history
	txDetails, err := w.ListAllTransactions()
	if err != nil {
		return fmt.Errorf("failed to list transactions: %w", err)
	}

	cleanHistory := aggregateTransactions(txDetails)
	if limit > 0 {
		cleanHistory = cleanHistory[:limit]
	}

	if len(cleanHistory) == 0 {
		fmt.Printf("No transactions found.")
		return nil
	}
	// Print transaction details
	fmt.Printf("Transactions\n-------------\n")
	printTransactionHistory(cleanHistory)

	return nil
}

type CleanTransactionType int

const (
	TRANSACTION_SENT CleanTransactionType = iota
	TRANSACTION_RECEIVED
	TRANSACTION_MINED
)

// CleanTransaction represents a simplified view of a transaction.
type CleanTransaction struct {
	Timestamp     string
	Amount        float64
	Type          CleanTransactionType
	TxID          string
	Address       string
	Confirmations int64
}

// printTransactionHistory processes and displays a clean transaction history.
func printTransactionHistory(txs []CleanTransaction) {

	for _, tx := range txs {
		fmt.Printf("%s %s %s %f %d\n",
			tx.Timestamp,
			tx.TxID,
			tx.Address,
			tx.Amount,
			tx.Confirmations,
		)
	}
}

// aggregateTransactions groups raw transactions into a simplified view.
func aggregateTransactions(txs []chainjson.ListTransactionsResult) []CleanTransaction {
	// Map to group transactions by TxID
	txMap := make(map[string]*CleanTransaction)

	for _, tx := range txs {
		// Format timestamp
		timestamp := time.Unix(tx.Time, 0).Format("2006-01-02 15:04:05")

		// Check if the transaction is already in the map
		cleanTx, exists := txMap[tx.TxID]
		if !exists {
			cleanTx = &CleanTransaction{
				Timestamp:     timestamp,
				Amount:        0,
				TxID:          tx.TxID,
				Address:       tx.Address,
				Confirmations: tx.Confirmations,
			}
			txMap[tx.TxID] = cleanTx
		}

		// Aggregate the amount (net debit/credit)
		cleanTx.Amount = tx.Amount

		// Set the transaction type and details
		if tx.Generated {
			cleanTx.Type = TRANSACTION_MINED
		} else if tx.Amount < 0 {
			cleanTx.Type = TRANSACTION_SENT
		} else {
			cleanTx.Type = TRANSACTION_RECEIVED
		}
	}

	// Collect the clean transactions into a slice
	cleanHistory := []CleanTransaction{}
	for _, tx := range txMap {
		cleanHistory = append(cleanHistory, *tx)
	}

	sort.Slice(cleanHistory, func(i, j int) bool {
		// Parse timestamps to compare
		timeI, _ := time.Parse("2006-01-02 15:04:05", cleanHistory[i].Timestamp)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", cleanHistory[j].Timestamp)
		return timeI.After(timeJ)
	})

	return cleanHistory
}
