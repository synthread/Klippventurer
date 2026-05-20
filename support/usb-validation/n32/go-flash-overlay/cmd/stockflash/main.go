package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/synthread/go-flash/flash"
)

const requiredConfirmToken = "FLASH_STOCK_FIRMWARE"

func main() {
	var firmware string
	var family string
	var tty string
	var baud int
	var addrText string
	var expectedIDs string
	var allowUnknown bool
	var machine string
	var execute bool
	var probeOnly bool
	var dryRun bool
	var confirmToken string
	var noGPIO bool

	flag.StringVar(&firmware, "firmware", "", "stock firmware file discovered on printer storage")
	flag.StringVar(&family, "family", string(flash.TargetFamilyN32G45), "target family: n32g45x or stm32")
	flag.StringVar(&tty, "tty", flash.DefaultTTY, "serial device")
	flag.IntVar(&baud, "baud", flash.DefaultBaud, "bootloader baud rate")
	flag.StringVar(&addrText, "addr", "0x08000000", "flash address")
	flag.StringVar(&expectedIDs, "expected-chip-id", "", "comma-separated expected chip IDs")
	flag.BoolVar(&allowUnknown, "allow-unknown-chip-id", false, "allow flashing without expected chip ID list")
	flag.StringVar(&machine, "machine", "auto", "machine/profile label for logs")
	flag.BoolVar(&execute, "execute", false, "perform the flash after all safeguards pass")
	flag.BoolVar(&probeOnly, "probe-only", false, "probe target and exit without flashing")
	flag.BoolVar(&dryRun, "dry-run", false, "validate inputs and probe target without flashing")
	flag.StringVar(&confirmToken, "confirm-token", "", "must equal FLASH_STOCK_FIRMWARE when --execute is used")
	flag.BoolVar(&noGPIO, "no-gpio-boot-entry", false, "do not configure or toggle GPIO boot pins")
	flag.Parse()

	if err := run(firmware, family, tty, baud, addrText, expectedIDs, allowUnknown, machine, execute, probeOnly, dryRun, confirmToken, noGPIO); err != nil {
		fmt.Fprintf(os.Stderr, "stockflash: %v\n", err)
		os.Exit(1)
	}
}

func run(firmware, family, tty string, baud int, addrText, expectedIDs string, allowUnknown bool, machine string, execute, probeOnly, dryRun bool, confirmToken string, noGPIO bool) error {
	if firmware == "" {
		return errors.New("--firmware is required")
	}
	info, err := os.Stat(firmware)
	if err != nil {
		return fmt.Errorf("stat firmware: %w", err)
	}
	if info.Size() <= 0 {
		return errors.New("firmware file is empty")
	}

	addr, err := parseUint32(addrText)
	if err != nil {
		return err
	}
	tf, err := parseFamily(family)
	if err != nil {
		return err
	}
	expected := splitCSV(expectedIDs)
	if execute {
		if confirmToken != requiredConfirmToken {
			return errors.New("--execute requires --confirm-token FLASH_STOCK_FIRMWARE")
		}
		if len(expected) == 0 && !allowUnknown {
			return errors.New("--execute requires --expected-chip-id unless --allow-unknown-chip-id is set")
		}
	}

	started := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("event=start time=%s machine=%s firmware=%s size=%d family=%s tty=%s baud=%d addr=0x%08x execute=%t probe_only=%t dry_run=%t\n", started, machine, firmware, info.Size(), tf, tty, baud, addr, execute, probeOnly, dryRun)

	mc, err := flash.NewMicrocontroller(&flash.Config{
		TargetFamily:       tf,
		TTY:                tty,
		BootloaderBaud:     baud,
		SkipGPIOSetup:      noGPIO,
		SkipBootEntry:      noGPIO,
		ExpectedChipIDs:    expected,
		AllowUnknownChipID: allowUnknown,
	})
	if err != nil {
		return fmt.Errorf("create microcontroller: %w", err)
	}

	probe, err := mc.Probe()
	if err != nil {
		return fmt.Errorf("probe failed: %w", err)
	}
	fmt.Printf("event=probe-ok chip_id=%s bootloader_version=0x%02x command_bytes=%s\n", probe.ChipID, probe.BootloaderVersion, hex.EncodeToString(probe.CommandBytes))
	if err := validateChipID(probe.ChipID, expected, allowUnknown); err != nil {
		return err
	}
	if probeOnly || dryRun || !execute {
		fmt.Printf("result=probe-only note=no flash performed\n")
		return nil
	}

	if err := mc.FlashPayloadFromFile(firmware, addr); err != nil {
		return fmt.Errorf("flash failed: %w", err)
	}
	fmt.Printf("result=flashed firmware=%s addr=0x%08x size=%d\n", firmware, addr, info.Size())
	return nil
}

func parseUint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid --addr %q: %w", s, err)
	}
	return uint32(v), nil
}

func parseFamily(s string) (flash.TargetFamily, error) {
	switch flash.TargetFamily(s) {
	case flash.TargetFamilyN32G45:
		return flash.TargetFamilyN32G45, nil
	case flash.TargetFamilySTM32:
		return flash.TargetFamilySTM32, nil
	default:
		return "", fmt.Errorf("unsupported family %q", s)
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func validateChipID(got string, expected []string, allowUnknown bool) error {
	if len(expected) == 0 && allowUnknown {
		fmt.Printf("event=chip-id-warning chip_id=%s note=no expected id supplied\n", got)
		return nil
	}
	if len(expected) == 0 {
		return fmt.Errorf("missing expected chip id; got %s", got)
	}
	for _, want := range expected {
		if strings.EqualFold(got, want) {
			return nil
		}
	}
	return fmt.Errorf("unexpected chip id %s; expected one of %s", got, strings.Join(expected, ","))
}
