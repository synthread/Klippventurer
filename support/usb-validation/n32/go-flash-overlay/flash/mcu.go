package flash

import (
	"github.com/piotrjaromin/gpio"
	"github.com/pkg/errors"
	"go.bug.st/serial"
)

var DefaultBaud = 115200
var DefaultTTY = "/dev/ttyS1"

type TargetFamily string

const (
	TargetFamilySTM32  TargetFamily = "stm32"
	TargetFamilyN32G45 TargetFamily = "n32g45x"
)

// Config defines configuration for communicating and flashing the
// microcontroller
type Config struct {
	Boot0GPIO int
	Boot1GPIO int
	PowerGPIO int

	TargetFamily       TargetFamily
	ExpectedChipIDs    []string
	AllowUnknownChipID bool
	TargetFlashBase    uint32
	TargetFlashSize    uint32
	BootloaderBaud     int
	TTY                string
	SkipGPIOSetup      bool
	SkipBootEntry      bool
}

// Microcontroller represents an embedded microntroller chip that can be
// communicated with over UART
type Microcontroller struct {
	config *Config

	pinPower gpio.Pin
	pinBoot0 gpio.Pin
	pinBoot1 gpio.Pin

	stmCmdCodes          commandCodeMap
	stmBootloaderVersion byte

	ttyPort serial.Port
	ttyRx   chan byte

	ttyActive bool

	identity string
	target   TargetFamily
}

// NewMicrocontroller will create a new reference to a particular chip
func NewMicrocontroller(c *Config) (*Microcontroller, error) {
	if c == nil {
		c = &Config{}
	}

	if c.Boot0GPIO <= 0 {
		c.Boot0GPIO = 39
	}
	if c.Boot1GPIO <= 0 {
		c.Boot1GPIO = 41
	}
	if c.PowerGPIO <= 0 {
		c.PowerGPIO = 19
	}

	mc := &Microcontroller{
		config:      c,
		stmCmdCodes: commandCodeMap{},
		target:      c.targetFamilyOrDefault(),
	}

	if !c.SkipGPIOSetup {
		if err := mc.setupPins(); err != nil {
			return nil, errors.Wrap(err, "could not setup pins")
		}
	}

	return mc, nil
}

func (c *Config) targetFamilyOrDefault() TargetFamily {
	if c != nil && c.TargetFamily != "" {
		return c.TargetFamily
	}
	return TargetFamilySTM32
}

func (mc *Microcontroller) TargetFamily() TargetFamily {
	if mc.target != "" {
		return mc.target
	}
	return TargetFamilySTM32
}

func (mc *Microcontroller) setupPins() (err error) {
	mc.pinPower, err = gpio.NewOutput(uint(mc.config.PowerGPIO), true)
	if err != nil {
		return
	}
	mc.pinBoot0, err = gpio.NewOutput(uint(mc.config.Boot0GPIO), false)
	if err != nil {
		return
	}
	mc.pinBoot1, err = gpio.NewOutput(uint(mc.config.Boot1GPIO), false)
	if err != nil {
		return
	}

	return
}

func (mc *Microcontroller) canControlBootPins() bool {
	return mc.config != nil && !mc.config.SkipGPIOSetup && !mc.config.SkipBootEntry
}

// Identify will report back a unique string with the ID of the chip
func (mc *Microcontroller) Identify() (string, error) {
	if mc.identity != "" {
		return mc.identity, nil
	}

	if !mc.IsOpen() {
		if err := mc.Open(); err != nil {
			return "", err
		}
		defer mc.Close()
	}

	pid, err := mc.familyGetID()
	if err != nil {
		return "", err
	}
	mc.identity = mc.familyIdentityPrefix() + pid

	return mc.identity, nil
}

// TTY will return the TTY that will be used
func (mc *Microcontroller) TTY() string {
	if mc.config.TTY != "" {
		return mc.config.TTY
	}
	return DefaultTTY
}

// BaudRate will return the baud rate used to connect to the TTY
func (mc *Microcontroller) BaudRate() int {
	if mc.config.BootloaderBaud > 0 {
		return mc.config.BootloaderBaud
	}
	return DefaultBaud
}

// Reset will force a power cycle on the microcontroller
func (mc *Microcontroller) Reset() {
	mc.familyExitBootloader()
}

func (mc *Microcontroller) familyInit() error {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		return mc.stmInit()
	case TargetFamilyN32G45:
		return mc.n32Init()
	default:
		return errors.Errorf("unsupported target family: %s", mc.TargetFamily())
	}
}

func (mc *Microcontroller) familyGetID() (string, error) {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		return mc.stmCmdGetId()
	case TargetFamilyN32G45:
		return mc.n32CmdGetID()
	default:
		return "", errors.Errorf("unsupported target family: %s", mc.TargetFamily())
	}
}

func (mc *Microcontroller) familyFlash(bs []byte, addr uint32) error {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		return mc.stmFlash(bs, addr)
	case TargetFamilyN32G45:
		return mc.n32Flash(bs, addr)
	default:
		return errors.Errorf("unsupported target family: %s", mc.TargetFamily())
	}
}

func (mc *Microcontroller) familyExitBootloader() {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		mc.exitSTBL()
	case TargetFamilyN32G45:
		mc.n32ExitBootloader()
	default:
		mc.exitSTBL()
	}
}

func (mc *Microcontroller) familyIdentityPrefix() string {
	switch mc.TargetFamily() {
	case TargetFamilySTM32:
		return "STM_"
	case TargetFamilyN32G45:
		return "N32G45x_"
	default:
		return string(mc.TargetFamily()) + "_"
	}
}
