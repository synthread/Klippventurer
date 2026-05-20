package flash

import "os"

// FlashPayloadFromFile will flash the requested file to the flash memory at the
// provided address
func (mc *Microcontroller) FlashPayloadFromFile(filePath string, addr uint32) error {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return mc.FlashPayload(bs, addr)
}

// FlashPayload will flash the payload provided to the flash at the provided
// address
func (mc *Microcontroller) FlashPayload(bs []byte, addr uint32) error {
	if !mc.IsOpen() {
		if err := mc.Open(); err != nil {
			return err
		}
		defer mc.Close()
	}

	if err := mc.familyFlash(bs, addr); err != nil {
		return err
	}

	return nil
}
