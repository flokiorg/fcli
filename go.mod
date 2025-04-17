module github.com/flokiorg/fcli

go 1.23.4

require (
	github.com/flokiorg/go-flokicoin v0.23.5-0.20230711222809-7faa9b266231
	github.com/flokiorg/walletd v0.0.0-20250227162953-b73c954e8fab
	github.com/jessevdk/go-flags v1.6.1
	golang.org/x/term v0.30.0
)

require (
	github.com/aead/siphash v1.0.1 // indirect
	github.com/btcsuite/btcd v0.24.2 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/decred/dcrd/lru v1.1.3 // indirect
	github.com/flokiorg/flokicoin-neutrino v0.0.0-00010101000000-000000000000 // indirect
	github.com/flokiorg/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/kkdai/bstream v1.0.0 // indirect
	github.com/lightningnetwork/lnd/clock v1.1.1 // indirect
	github.com/lightningnetwork/lnd/fn/v2 v2.0.8 // indirect
	github.com/lightningnetwork/lnd/queue v1.1.1 // indirect
	github.com/lightningnetwork/lnd/tlv v1.3.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.etcd.io/bbolt v1.4.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250210185358-939b2ce775ac // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

require (
	github.com/decred/dcrd/crypto/blake256 v1.1.0 // indirect
	github.com/lightninglabs/gozmq v0.0.0-20191113021534-d20a764486bf // indirect
	github.com/lightningnetwork/lnd/ticker v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/flokiorg/go-flokicoin => ../go-flokicoin

replace github.com/flokiorg/walletd => ../walletd

replace github.com/flokiorg/flokicoin-neutrino => ../flokicoin-neutrino
