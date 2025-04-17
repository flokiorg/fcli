// Copyright (c) 2024 The Flokicoin developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package utils

import (
	"fmt"

	"github.com/flokiorg/go-flokicoin/chaincfg"
	"github.com/flokiorg/go-flokicoin/chainutil"
	"github.com/flokiorg/go-flokicoin/chainutil/hdkeychain"
	"github.com/flokiorg/walletd/waddrmgr"
)

func StrAddrType(t waddrmgr.AddressType) string {
	switch t {
	case waddrmgr.Script:
		return "Script"
	case waddrmgr.RawPubKey:
		return "RawPubKey"
	case waddrmgr.NestedWitnessPubKey:
		return "NestedWitnessPubKey"
	case waddrmgr.WitnessPubKey:
		return "WitnessPubKey"
	case waddrmgr.WitnessScript:
		return "WitnessScript"
	case waddrmgr.TaprootPubKey:
		return "TaprootPubKey"
	case waddrmgr.TaprootScript:
		return "TaprootScript"

	// case waddrmgr.PubKeyHash:
	default:
		return "PubKeyHash"
	}
}

func DeriveKeysFromXpub(network *chaincfg.Params, xpriv, xpub string, childIndex uint32) (string, string, error) {

	extKey, err := hdkeychain.NewKeyFromString(xpub)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse xpub: %v", err)
	}

	childKey, err := extKey.Derive(childIndex)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive child key: %v", err)
	}

	add, err := childKey.Address(network)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate address: %v", err)
	}

	var strWif string

	if xpriv != "" {
		privExtKey, err := hdkeychain.NewKeyFromString(xpriv)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse xpriv: %v", err)
		}

		privBranchKey, err := privExtKey.Derive(waddrmgr.ExternalBranch)
		if err != nil {
			return "", "", err
		}

		privChildKey, err := privBranchKey.Derive(childIndex)
		if err != nil {
			return "", "", fmt.Errorf("failed to derive private key: %v", err)
		}

		privKey, err := privChildKey.ECPrivKey()
		if err != nil {
			return "", "", fmt.Errorf("failed to extract EC private key: %v", err)
		}
		wif, err := chainutil.NewWIF(privKey, network, true)
		if err != nil {
			return "", "", fmt.Errorf("failed to convert private key to WIF: %v", err)
		}
		strWif = wif.String()
	}

	return add.EncodeAddress(), strWif, nil
}
