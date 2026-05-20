package flash

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const n32FlashBase uint32 = 0x08000000
const n32FlashSize uint32 = 0x00020000

var ErrN32FlashOutOfRange = errors.New("n32 flash request is outside known N32G45x flash range")
var ErrN32MissingExpectedChipID = errors.New("n32 flashing requires expected chip id unless AllowUnknownChipID is set")
var ErrN32UnexpectedChipID = errors.New("n32 target chip id did not match expected id list")

func (mc *Microcontroller) n32Init() error {
	// Evidence from Klipper #6116 describes N32G452/G455 as mostly compatible
	// with STM32F103. Start with the STM ROM-bootloader handshake and command
	// discovery, but keep N32 behind an explicit target family.
	return mc.stmInit()
}

func (mc *Microcontroller) n32CmdGetID() (string, error) {
	return mc.stmCmdGetId()
}

func (mc *Microcontroller) n32Probe() (*ProbeResult, error) {
	return mc.stmCompatibleProbe(TargetFamilyN32G45)
}

func (mc *Microcontroller) n32Flash(bs []byte, addr uint32) error {
	if err := mc.n32ValidateFlashRange(bs, addr); err != nil {
		return err
	}
	if err := mc.n32ValidateChipID(); err != nil {
		return err
	}

	return mc.stmFlash(bs, addr)
}

func (mc *Microcontroller) n32ExitBootloader() {
	// Use the same boot-pin reset sequence as the STM path until board-level
	// validation proves N32 Adventurer 3 variants require different behavior.
	mc.exitSTBL()
}

func (mc *Microcontroller) n32ValidateFlashRange(bs []byte, addr uint32) error {
	if len(bs) == 0 {
		return errors.New("n32 flash payload is empty")
	}

	end := uint64(addr) + uint64(len(bs))
	flashBase, flashSize := mc.n32FlashBounds()
	flashStart := uint64(flashBase)
	flashEnd := uint64(flashBase) + uint64(flashSize)
	if uint64(addr) < flashStart || end > flashEnd {
		return errors.Wrapf(
			ErrN32FlashOutOfRange,
			"addr=0x%08x len=%d allowed=0x%08x..0x%08x",
			addr,
			len(bs),
			flashBase,
			flashBase+flashSize,
		)
	}

	return nil
}

func (mc *Microcontroller) n32FlashBounds() (base uint32, size uint32) {
	base = n32FlashBase
	if mc.config.TargetFlashBase != 0 {
		base = mc.config.TargetFlashBase
	}

	size = n32FlashSize
	if mc.config.TargetFlashSize != 0 {
		size = mc.config.TargetFlashSize
	}

	return base, size
}

func (mc *Microcontroller) n32ValidateChipID() error {
	pid, err := mc.n32CmdGetID()
	if err != nil {
		return errors.Wrap(err, "could not read n32 chip id before flash")
	}

	if len(mc.config.ExpectedChipIDs) == 0 && mc.config.AllowUnknownChipID {
		logrus.Warnf("n32 chip id %s has not been checked against an expected id list", pid)
		return nil
	}
	if len(mc.config.ExpectedChipIDs) == 0 {
		return errors.Wrapf(ErrN32MissingExpectedChipID, "got %s", pid)
	}

	for _, expected := range mc.config.ExpectedChipIDs {
		if strings.EqualFold(pid, expected) {
			return nil
		}
	}

	return errors.Wrapf(ErrN32UnexpectedChipID, "got %s expected one of %s", pid, strings.Join(mc.config.ExpectedChipIDs, ","))
}
