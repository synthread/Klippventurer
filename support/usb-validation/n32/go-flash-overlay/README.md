# go-flash

`go-flash` lets you flash microcontrollers over UART, using GPIO pins to
control boot0/boot1/pwr lines.

Supported target families:

- `stm32` (default)
- `n32g45x` (experimental; uses the STM-compatible ROM bootloader flow with
  N32-specific target selection and range checks)

N32 support is intended for Nation N32G452/N32G455-class devices. It should be
used with explicit board-profile data for the application address, allowed flash
range, and expected chip ID.

## Installation

```sh
go get github.com/synthread/go-flash
```

## Usage

See https://pkg.go.dev/github.com/synthread/go-flash/flash for some docs
on how to use it.

For N32G45x targets, set `Config.TargetFamily` to `flash.TargetFamilyN32G45`.
Set `Config.ExpectedChipIDs` so flashing fails closed when the target is not the
expected MCU. During hardware bring-up only, `Config.AllowUnknownChipID` can be
set to permit flashing after logging the observed chip ID. Use
`Config.TargetFlashBase` and `Config.TargetFlashSize` when the board profile
requires bounds different from the N32 defaults.
