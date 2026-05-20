package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/synthread/go-flash/flash"
	"go.bug.st/serial"
)

const defaultSerialBootRequest = "~ \x1c Request Serial Bootloader!! ~"

type result struct {
	Name       string            `json:"name"`
	OK         bool              `json:"ok"`
	Warning    bool              `json:"warning,omitempty"`
	Message    string            `json:"message,omitempty"`
	Probe      *probeResultJSON  `json:"probe,omitempty"`
	Command    string            `json:"command,omitempty"`
	ExitStatus int               `json:"exit_status,omitempty"`
	Extra      map[string]string `json:"extra,omitempty"`
}

type probeResultJSON struct {
	TargetFamily      string            `json:"target_family"`
	TTY               string            `json:"tty"`
	BaudRate          int               `json:"baud_rate"`
	BootloaderVersion string            `json:"bootloader_version"`
	CommandCodes      map[string]string `json:"command_codes"`
	CommandBytes      string            `json:"command_bytes"`
	ChipID            string            `json:"chip_id"`
	RawGet            string            `json:"raw_get"`
	RawGetID          string            `json:"raw_get_id"`
}

func main() {
	var tty string
	var baud int
	var family string
	var bootEntry string
	var serialBootBaud int
	var serialBootRequest string
	var machine string
	var jsonOut bool
	var noGPIO bool
	var expectedIDs string
	var allowUnknown bool
	var nationsCmd string
	var mainboardTTY string
	var probeMainboardROM bool
	var scanTTYs string
	var bootDelay time.Duration

	flag.StringVar(&tty, "tty", flash.DefaultTTY, "serial device for N32 ROM probe")
	flag.IntVar(&baud, "baud", flash.DefaultBaud, "baud rate for N32 ROM probe")
	flag.StringVar(&family, "family", string(flash.TargetFamilyN32G45), "target family: n32g45x or stm32")
	flag.StringVar(&bootEntry, "boot-entry", "auto", "boot entry mode: auto, gpio, serial-request, none")
	flag.IntVar(&serialBootBaud, "serial-boot-request-baud", 230400, "baud for FlashForge serial bootloader request")
	flag.StringVar(&serialBootRequest, "serial-boot-request", defaultSerialBootRequest, "FlashForge serial bootloader request string")
	flag.StringVar(&machine, "machine", "auto", "machine/profile: auto, adv3-mips, arm-n32-single, ad5m")
	flag.BoolVar(&jsonOut, "json", false, "emit JSON")
	flag.BoolVar(&noGPIO, "no-gpio-boot-entry", false, "do not configure or toggle GPIO boot pins")
	flag.StringVar(&expectedIDs, "expected-chip-id", "", "comma-separated expected chip IDs")
	flag.BoolVar(&allowUnknown, "allow-unknown-chip-id", false, "do not fail when no expected chip ID is supplied")
	flag.StringVar(&nationsCmd, "nations-command", "", "path to NationsCommand for 5M mainboard contact/reset validation")
	flag.StringVar(&mainboardTTY, "mainboard-tty", "/dev/ttyS5", "5M mainboard serial device")
	flag.BoolVar(&probeMainboardROM, "probe-mainboard-rom", false, "best-effort ROM probe on mainboard TTY after NationsCommand reset")
	flag.StringVar(&scanTTYs, "scan-ttys", "", "comma-separated serial devices to try after the primary TTY")
	flag.DurationVar(&bootDelay, "boot-delay", time.Second, "delay after boot entry before ROM probe")
	flag.Parse()

	results := []result{}
	if machine == "auto" {
		machine = detectMachine()
	}

	if bootEntry == "auto" {
		if machine == "adv3-mips" {
			bootEntry = "gpio"
		} else {
			bootEntry = "serial-request"
		}
	}

	if machine == "ad5m" {
		results = append(results, validateMainboard(nationsCmd, mainboardTTY, family, baud, probeMainboardROM, expectedIDs, allowUnknown, jsonOut))
	}

	probeName := "mainboard"
	if machine == "ad5m" {
		probeName = "eboard"
	}
	results = append(results, validateProbe(probeName, tty, baud, family, bootEntry, serialBootBaud, serialBootRequest, noGPIO, expectedIDs, allowUnknown, bootDelay))
	if !results[len(results)-1].OK {
		for _, candidate := range splitCSV(scanTTYs) {
			if candidate == tty {
				continue
			}
			scanResult := validateProbe(probeName+"-scan-"+candidate, candidate, baud, family, "none", serialBootBaud, serialBootRequest, true, expectedIDs, allowUnknown, 0)
			results = append(results, scanResult)
			if scanResult.OK {
				break
			}
		}
	}

	ok := true
	for _, r := range results {
		if !r.OK {
			ok = false
		}
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]any{"ok": ok, "machine": machine, "results": results})
	} else {
		fmt.Printf("n32probe read-only validator\n")
		fmt.Printf("machine=%s ok=%v\n", machine, ok)
		for _, r := range results {
			printResult(r)
		}
	}

	if !ok {
		os.Exit(1)
	}
}

func validateProbe(name, tty string, baud int, family, bootEntry string, serialBootBaud int, serialBootRequest string, noGPIO bool, expectedIDs string, allowUnknown bool, bootDelay time.Duration) result {
	if bootEntry == "serial-request" {
		if err := sendSerialBootRequest(tty, serialBootBaud, serialBootRequest); err != nil {
			return result{Name: name, OK: false, Message: "serial boot request failed: " + err.Error()}
		}
		time.Sleep(bootDelay)
	} else if bootEntry == "gpio" && bootDelay > 0 {
		tf, err := parseFamily(family)
		if err != nil {
			return result{Name: name, OK: false, Message: err.Error()}
		}
		mc, err := flash.NewMicrocontroller(&flash.Config{TargetFamily: tf, TTY: tty, BootloaderBaud: baud})
		if err != nil {
			return result{Name: name, OK: false, Message: err.Error()}
		}
		mc.PrepareBootloader()
		time.Sleep(bootDelay)
		noGPIO = true
		bootEntry = "none"
	}

	tf, err := parseFamily(family)
	if err != nil {
		return result{Name: name, OK: false, Message: err.Error()}
	}

	mc, err := flash.NewMicrocontroller(&flash.Config{
		TargetFamily:       tf,
		TTY:                tty,
		BootloaderBaud:     baud,
		SkipGPIOSetup:      noGPIO || bootEntry != "gpio",
		SkipBootEntry:      noGPIO || bootEntry != "gpio",
		ExpectedChipIDs:    splitCSV(expectedIDs),
		AllowUnknownChipID: allowUnknown,
	})
	if err != nil {
		return result{Name: name, OK: false, Message: err.Error()}
	}

	probe, err := mc.Probe()
	if err != nil {
		return result{Name: name, OK: false, Message: err.Error()}
	}
	if err := validateChipID(probe.ChipID, splitCSV(expectedIDs), allowUnknown); err != nil {
		return result{Name: name, OK: false, Message: err.Error(), Probe: jsonifyProbe(probe)}
	}

	return result{Name: name, OK: true, Message: "read-only probe passed", Probe: jsonifyProbe(probe)}
}

func validateMainboard(nationsCmd, tty, family string, baud int, probeROM bool, expectedIDs string, allowUnknown bool, jsonOut bool) result {
	cmdPath := nationsCmd
	if cmdPath == "" {
		cmdPath = findNationsCommand()
	}
	if cmdPath == "" {
		return result{Name: "mainboard", OK: true, Warning: true, Message: "NationsCommand not found; best-effort mainboard contact/reset skipped"}
	}

	cmd := exec.Command(cmdPath, "-c", "--pn", tty, "-r")
	out, err := cmd.CombinedOutput()
	r := result{Name: "mainboard", OK: err == nil, Command: strings.Join(cmd.Args, " "), Message: strings.TrimSpace(string(out)), Extra: map[string]string{"tty": tty}}
	if cmd.ProcessState != nil {
		r.ExitStatus = cmd.ProcessState.ExitCode()
	}
	if err != nil {
		r.Message = strings.TrimSpace(string(out)) + "\n" + err.Error()
		r.OK = true
		r.Warning = true
		return r
	}

	if probeROM {
		pr := validateProbe("mainboard-rom", tty, baud, family, "none", 0, "", true, expectedIDs, allowUnknown, 0)
		if pr.OK {
			r.Probe = pr.Probe
		} else {
			r.Warning = true
			r.Extra["rom_probe_warning"] = pr.Message
		}
	}

	return r
}

func sendSerialBootRequest(tty string, baud int, request string) error {
	port, err := serial.Open(tty, &serial.Mode{BaudRate: baud, DataBits: 8, Parity: serial.NoParity, StopBits: serial.OneStopBit})
	if err != nil {
		return err
	}
	defer port.Close()
	_, err = io.WriteString(port, unescapeRequest(request))
	return err
}

func unescapeRequest(s string) string {
	return strings.ReplaceAll(s, "\\x1c", string([]byte{0x1c}))
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
	if len(expected) == 0 {
		if allowUnknown {
			return nil
		}
		return errors.New("no expected chip ID configured; use --allow-unknown-chip-id for bring-up")
	}
	for _, e := range expected {
		if strings.EqualFold(got, e) {
			return nil
		}
	}
	return fmt.Errorf("chip ID %s did not match expected IDs %s", got, strings.Join(expected, ","))
}

func jsonifyProbe(p *flash.ProbeResult) *probeResultJSON {
	commands := map[string]string{}
	for k, v := range p.CommandCodes {
		commands[commandName(k)] = fmt.Sprintf("0x%02x", v)
	}
	return &probeResultJSON{
		TargetFamily:      string(p.TargetFamily),
		TTY:               p.TTY,
		BaudRate:          p.BaudRate,
		BootloaderVersion: fmt.Sprintf("0x%02x", p.BootloaderVersion),
		CommandCodes:      commands,
		CommandBytes:      hex.EncodeToString(p.CommandBytes),
		ChipID:            p.ChipID,
		RawGet:            hex.EncodeToString(p.RawGet),
		RawGetID:          hex.EncodeToString(p.RawGetID),
	}
}

func commandName(c flash.CommandCode) string {
	switch c {
	case flash.CommandCodeGet:
		return "get"
	case flash.CommandCodeGetVersion:
		return "get_version"
	case flash.CommandCodeGetID:
		return "get_id"
	case flash.CommandCodeReadMemory:
		return "read_memory"
	case flash.CommandCodeGo:
		return "go"
	case flash.CommandCodeWriteMemory:
		return "write_memory"
	case flash.CommandCodeErase:
		return "erase"
	case flash.CommandCodeWriteProtect:
		return "write_protect"
	case flash.CommandCodeWriteUnprotect:
		return "write_unprotect"
	case flash.CommandCodeReadoutProtect:
		return "readout_protect"
	case flash.CommandCodeReadoutUnprotect:
		return "readout_unprotect"
	default:
		return fmt.Sprintf("command_%d", c)
	}
}

func printResult(r result) {
	status := "FAIL"
	if r.OK {
		status = "PASS"
	}
	if r.Warning {
		status += " (WARN)"
	}
	fmt.Printf("\n[%s] %s\n", status, r.Name)
	if r.Message != "" {
		fmt.Println(r.Message)
	}
	if r.Command != "" {
		fmt.Println("command:", r.Command)
	}
	if r.Probe != nil {
		fmt.Printf("target=%s tty=%s baud=%d bootloader=%s chip_id=%s\n", r.Probe.TargetFamily, r.Probe.TTY, r.Probe.BaudRate, r.Probe.BootloaderVersion, r.Probe.ChipID)
		fmt.Printf("command_bytes=%s raw_get=%s raw_get_id=%s\n", r.Probe.CommandBytes, r.Probe.RawGet, r.Probe.RawGetID)
	}
}

func detectMachine() string {
	if fileContains("/etc/machine", "Adventurer5M") || fileContains("/etc/hostname", "Adventurer5M") || os.Getenv("N32PROBE_MACHINE") == "ad5m" {
		return "ad5m"
	}
	if strings.Contains(runtimeArch(), "mips") {
		return "adv3-mips"
	}
	return "arm-n32-single"
}

func runtimeArch() string {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func fileContains(path, needle string) bool {
	bs, err := os.ReadFile(path)
	return err == nil && strings.Contains(string(bs), needle)
}

func findNationsCommand() string {
	for _, p := range []string{
		"/opt/PROGRAM/control/2.2.3/NationsCommand",
		"/mnt/orig_root/opt/PROGRAM/control/2.2.3/NationsCommand",
		"./NationsCommand",
	} {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}
