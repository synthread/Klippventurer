package flash

import "github.com/pkg/errors"

type ProbeResult struct {
	TargetFamily      TargetFamily
	TTY               string
	BaudRate          int
	BootloaderVersion byte
	CommandCodes      map[CommandCode]byte
	CommandBytes      []byte
	ChipID            string
	RawGet            []byte
	RawGetID          []byte
}

// Probe performs a read-only ROM bootloader probe. It synchronizes with the
// target, runs GET, runs GET ID, captures raw responses, and exits bootloader
// mode. It must not call erase, write, unprotect, protect, read memory, or go.
func (mc *Microcontroller) Probe() (*ProbeResult, error) {
	opened := false
	if !mc.IsOpen() {
		if err := mc.openSerialPort(); err != nil {
			return nil, err
		}
		opened = true
	}

	defer func() {
		if opened {
			_ = mc.Close()
		} else {
			mc.familyExitBootloader()
		}
	}()

	return mc.familyProbe()
}

// PrepareBootloader performs only the configured boot-entry sequence. It does
// not open the serial port or run any bootloader command.
func (mc *Microcontroller) PrepareBootloader() {
	switch mc.TargetFamily() {
	case TargetFamilySTM32, TargetFamilyN32G45:
		mc.enterSTBL()
	}
}

func (mc *Microcontroller) familyProbe() (*ProbeResult, error) {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		return mc.stmProbe()
	case TargetFamilyN32G45:
		return mc.n32Probe()
	default:
		return nil, errors.Errorf("unsupported target family: %s", mc.TargetFamily())
	}
}

func (mc *Microcontroller) stmProbe() (*ProbeResult, error) {
	return mc.stmCompatibleProbe(TargetFamilySTM32)
}

func (mc *Microcontroller) stmCompatibleProbe(family TargetFamily) (*ProbeResult, error) {
	mc.enterSTBL()

	if err := mc.stmCmdSync(); err != nil {
		return nil, err
	}

	rawGet, err := mc.stmCmdGetRaw()
	if err != nil {
		return nil, err
	}

	chipID, rawGetID, err := mc.stmCmdGetIdRaw()
	if err != nil {
		return nil, err
	}

	result := &ProbeResult{
		TargetFamily:      family,
		TTY:               mc.TTY(),
		BaudRate:          mc.BaudRate(),
		BootloaderVersion: mc.stmBootloaderVersion,
		CommandCodes:      copyCommandCodeMap(mc.stmCmdCodes),
		CommandBytes:      commandBytesFromGet(rawGet),
		ChipID:            chipID,
		RawGet:            append([]byte(nil), rawGet...),
		RawGetID:          append([]byte(nil), rawGetID...),
	}

	return result, nil
}

func copyCommandCodeMap(in commandCodeMap) map[CommandCode]byte {
	out := make(map[CommandCode]byte, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func commandBytesFromGet(rawGet []byte) []byte {
	if len(rawGet) <= 1 {
		return nil
	}
	return append([]byte(nil), rawGet[1:]...)
}
