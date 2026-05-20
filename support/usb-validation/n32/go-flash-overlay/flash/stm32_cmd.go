package flash

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/pkg/errors"
)

// stmExecCmd will run the specified command and check that it is ACK'd
func (mc *Microcontroller) stmExecCmd(c CommandCode) error {
	if err := mc.Write(mc.stmCommandSequence(c)); err != nil {
		return err
	}
	if err := mc.stmReadAckOrNack(); err != nil {
		return err
	}
	return nil
}

// stmCmdSync will sync the bootloader
func (mc *Microcontroller) stmCmdSync() (err error) {
	return mc.stmExecCmd(CommandCodeSync)
}

// stmCmdGet will load information about the bootloader
func (mc *Microcontroller) stmCmdGet() error {
	_, err := mc.stmCmdGetRaw()
	return err
}

func (mc *Microcontroller) stmCmdGetRaw() ([]byte, error) {
	if err := mc.stmExecCmd(CommandCodeGet); err != nil {
		return nil, err
	}

	bs, err := mc.stmReadWithLength()
	if err != nil {
		return nil, err
	}

	if err = mc.stmReadAckOrNack(); err != nil {
		return nil, err
	}

	mc.stmApplyGetResponse(bs)

	return bs, nil
}

func (mc *Microcontroller) stmApplyGetResponse(bs []byte) {
	if len(bs) == 0 {
		return
	}

	mc.stmBootloaderVersion = bs[0]
	mc.stmCmdCodes = commandCodeMap{}

	// get the command codes from the response
	for i := 0; i < len(bs)-1; i++ {
		mc.stmCmdCodes[CommandCode(i)] = bs[i+1]
	}
}

// stmCmdGetId will return the PID of the microcontroller
func (mc *Microcontroller) stmCmdGetId() (string, error) {
	pid, _, err := mc.stmCmdGetIdRaw()
	return pid, err
}

func (mc *Microcontroller) stmCmdGetIdRaw() (string, []byte, error) {
	if err := mc.stmExecCmd(CommandCodeGetID); err != nil {
		return "", nil, err
	}

	bs, err := mc.stmReadWithLength()
	if err != nil {
		return "", nil, err
	}

	if err = mc.stmReadAckOrNack(); err != nil {
		return "", nil, err
	}

	return hex.EncodeToString(bs), bs, nil
}

// stmCmdEraseMemory will request that all flash memory be erased
func (mc *Microcontroller) stmCmdEraseMemory() error {
	if err := mc.stmExecCmd(CommandCodeErase); err != nil {
		return err
	}

	if err := mc.Write([]byte{0xff, 0x00}); err != nil {
		return err
	}

	return mc.stmReadAckOrNack()
}

// stmCmdWriteMemory will attempt to write the requested data at the provided
// address in memory
func (mc *Microcontroller) stmCmdWriteMemory(addr uint32, data []byte) error {
	if err := mc.stmExecCmd(CommandCodeWriteMemory); err != nil {
		return errors.Wrap(err, "err exec write mem")
	}

	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, addr); err != nil {
		return err
	}

	addrbs := buf.Bytes()

	// write the address and its checksum
	if err := mc.Write(append(addrbs, checksum(addrbs))); err != nil {
		return errors.Wrap(err, "err writing addr")
	}
	if err := mc.stmReadAckOrNack(); err != nil {
		return errors.Wrap(err, "addr ack fail")
	}

	// write the data with length and checksum
	if err := mc.stmWriteWithNAndChecksum(data); err != nil {
		return errors.Wrap(err, "err writing data")
	}

	return errors.Wrap(mc.stmReadAckOrNack(), "err ack after write data")
}

// stmCmdWriteUnprotect will set flash to be unprotected so that we can write it
func (mc *Microcontroller) stmCmdWriteUnprotect() error {
	if err := mc.stmExecCmd(CommandCodeWriteUnprotect); err != nil {
		return err
	}
	// this does ACK twice, once for the command and once for the unprotect
	if err := mc.stmReadAckOrNack(); err != nil {
		return err
	}

	// we want to resync after this since it will reset the chip
	return mc.stmCmdSync()
}
