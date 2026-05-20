package flash

import (
	"fmt"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const b_STM_ACK byte = 0x79
const b_STM_NACK byte = 0x1f
const b_STM_SYNC byte = 0x7f
const stmFlashBlockMax = 256

var STMTimeout = 5 * time.Second

var ErrSTMFailedToAck = errors.New("failed to read ack or nack from stm microcontroller")
var ErrSTMNACK = errors.New("received nack from stm microcontroller")

type CommandCode int

// these must be the index of the bytes as received in the get data call
const (
	CommandCodeSync             CommandCode = -1
	CommandCodeGet              CommandCode = 0
	CommandCodeGetVersion       CommandCode = 1
	CommandCodeGetID            CommandCode = 2
	CommandCodeReadMemory       CommandCode = 3
	CommandCodeGo               CommandCode = 4
	CommandCodeWriteMemory      CommandCode = 5
	CommandCodeErase            CommandCode = 6
	CommandCodeWriteProtect     CommandCode = 7
	CommandCodeWriteUnprotect   CommandCode = 8
	CommandCodeReadoutProtect   CommandCode = 9
	CommandCodeReadoutUnprotect CommandCode = 10
)

// these are the default command codes
type commandCodeMap map[CommandCode]byte

var defaultCmdCodeMap map[CommandCode]byte = map[CommandCode]byte{
	CommandCodeGet:              0x00,
	CommandCodeGetVersion:       0x01,
	CommandCodeGetID:            0x02,
	CommandCodeReadMemory:       0x11,
	CommandCodeGo:               0x21,
	CommandCodeWriteMemory:      0x31,
	CommandCodeErase:            0x43,
	CommandCodeWriteProtect:     0x63,
	CommandCodeWriteUnprotect:   0x73,
	CommandCodeReadoutProtect:   0x82,
	CommandCodeReadoutUnprotect: 0x92,
}

func (mc *Microcontroller) stmInit() error {
	mc.enterSTBL()

	if err := mc.stmCmdSync(); err != nil {
		return err
	}
	return mc.stmCmdGet()
}

// enterSTBL will execute the GPIO sequence to enter the STM bootloader
func (mc *Microcontroller) enterSTBL() {
	if !mc.canControlBootPins() {
		return
	}

	mc.pinPower.Low()

	// BOOT0 low and BOOT1 high when reapplying PWR will go into the bootloader
	// mode on STM32 chips
	mc.pinBoot0.High()
	mc.pinBoot1.Low()
	time.Sleep(10 * time.Millisecond)
	mc.pinPower.High()
	time.Sleep(10 * time.Millisecond)
}

// exitSTBL will execute the GPIO sequence to exit the STM bootloader
func (mc *Microcontroller) exitSTBL() {
	if !mc.canControlBootPins() {
		return
	}

	mc.pinPower.Low()
	mc.pinBoot0.Low()
	mc.pinBoot1.Low()
	time.Sleep(10 * time.Millisecond)
	mc.pinPower.High()
	time.Sleep(10 * time.Millisecond)
}

// stmCommandSequence will return the byte sequence required for the requested
// command
func (mc *Microcontroller) stmCommandSequence(c CommandCode) []byte {
	if c == CommandCodeSync {
		return []byte{b_STM_SYNC}
	}

	cmdb, ok := mc.stmCmdCodes[c]
	if !ok {
		cmdb, ok = defaultCmdCodeMap[c]
		if !ok {
			panic("unknown command code")
		}
	}

	return []byte{cmdb, 0xff ^ cmdb}
}

// stmReadWithLength will read the next bytes based on a STM formatted message
// which is prefixed by a single byte that represents the length of the
// expected message
func (mc *Microcontroller) stmReadWithLength() ([]byte, error) {
	n, err := mc.ReadN(1, STMTimeout)
	if err != nil {
		return nil, err
	}
	if len(n) != 1 {
		return nil, errors.New("could not get length from stm microcontroller")
	}
	return mc.ReadN(int(n[0])+1, STMTimeout)
}

// stmWriteWithChecksum will write the requested data with a checksum at the end
func (mc *Microcontroller) stmWriteWithChecksum(bs []byte) error {
	cs := checksum(bs)
	return mc.Write(append(bs, cs))
}

// stmWriteWithNAndChecksum will write the data prefixed with the length in a
// single byte and suffixed with the checksum of the entire message
func (mc *Microcontroller) stmWriteWithNAndChecksum(bs []byte) error {
	n := byte(len(bs) - 1)
	return mc.stmWriteWithChecksum(append([]byte{n}, bs...))
}

// stmReadAckOrNack reads whether the pending byte is ACK, NACK, or neither
// and returns the ACK/NACK status, whether it is valid, and an optional error
// if it could not be read or timed out
func (mc *Microcontroller) stmReadAckOrNack() (err error) {
	bs, err := mc.ReadN(1, STMTimeout)
	if err != nil || len(bs) != 1 {
		return
	}

	if bs[0] == b_STM_ACK {
		return nil
	} else if bs[0] == b_STM_NACK {
		return ErrSTMNACK
	}

	return ErrSTMFailedToAck
}

func (mc *Microcontroller) stmFlash(bs []byte, addr uint32) error {
	// if err := mc.stmCmdWriteUnprotect(); err != nil {
	// 	return errors.Wrap(err, "could not write unprotect")
	// }

	if err := mc.stmCmdEraseMemory(); err != nil {
		return errors.Wrap(err, "could not erase memory")
	}

	nseg := int(math.Ceil(float64(len(bs)) / stmFlashBlockMax))

	for i := 0; i < nseg; i++ {

		offset := i * stmFlashBlockMax
		segAddr := addr + uint32(offset)
		endIndex := min(len(bs), offset+stmFlashBlockMax)

		logrus.Debugf("wm: %d -> %d @ %x [l=%d]", offset, endIndex, segAddr, len(bs[offset:endIndex]))

		if err := mc.stmCmdWriteMemory(segAddr, bs[offset:endIndex]); err != nil {
			return errors.Wrap(err, fmt.Sprintf("could not write segment %d", i))
		}
	}

	return nil
}
