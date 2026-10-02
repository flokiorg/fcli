# Changelog

## [0.1.7]

### Changed

- Built with Go 1.26.5, picking up the stdlib security fixes released since 1.26.1. (#3)

## [0.1.6-beta]

### Fixed

- `wallet/service.go`: `RelayFee` and `EstimateFee` both discarded the cancel function
  from `context.WithTimeout`, leaking the timer until the 10s timeout fired on its own
  instead of releasing it when the function returned.

## [0.1.5-beta]

### Fixed

- The Electrum wallet-sync path called `chainClient.Start` without the now-required
  `context.Context` argument. This surfaced as a build failure against walletd
  v0.2.0-beta rather than as a runtime bug — the package simply could not compile
  against the new walletd interface until it was fixed.

### Changed

- Updated `go-flokicoin` to v0.26.0-alpha and `walletd` to v0.2.0-beta.

## [0.1.4-beta]

### Changed

- Updated dependencies to align with `go-flokicoin` v0.25.13-alpha and `walletd`
  v0.1.8-beta.

## [0.1.3-beta]

### Changed

- Updated `go-flokicoin` to v0.25.13-alpha, which includes the TestNet4 P2P port
  correction.

## [0.1.2-beta]

### Changed

- Updated `flokicoin-neutrino` to v0.16.2-beta and `go-flokicoin` to v0.25.7-beta.

## [0.1.1-alpha]

- First pre-release, published for testing and feedback.
