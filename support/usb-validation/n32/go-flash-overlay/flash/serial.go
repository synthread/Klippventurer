package flash

import (
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var ErrTimeout = errors.New("timed out reading from microcontroller")
var ErrClosed = errors.New("serial port is closed")

func (mc *Microcontroller) Open() (err error) {
	if err := mc.openSerialPort(); err != nil {
		return err
	}

	if err = errors.Wrap(mc.familyInit(), "could not init target chip"); err != nil {
		mc.Close()
		return
	}

	logrus.Debug("mcu open")

	return nil
}

func (mc *Microcontroller) openSerialPort() (err error) {
	mc.ttyPort, err = serial.Open(mc.TTY(), &serial.Mode{
		BaudRate: mc.BaudRate(),
		DataBits: 8,
		Parity:   serial.EvenParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return errors.Wrap(err, "could not open serial")
	}

	mc.ttyRx = make(chan byte, 64)
	go mc.rx()

	return nil
}

// Close will close the connection and reset the MCU
func (mc *Microcontroller) Close() error {
	mc.familyExitBootloader()

	mc.ttyActive = false

	if mc.ttyPort != nil {
		mc.ttyPort.Close()
		mc.ttyPort = nil
	}

	return nil
}

func (mc *Microcontroller) IsOpen() bool {
	// TODO: probably a better way to handle this
	return mc.ttyPort != nil
}

// rx is the loop that will forever read from the port and write the incoming
// bytes to the rx chan
func (mc *Microcontroller) rx() {
	mc.ttyActive = true
	buf := make([]byte, 64)

	defer mc.Close()

	mc.ttyPort.SetReadTimeout(1 * time.Millisecond)

	for {
		n, err := mc.ttyPort.Read(buf)
		if err != nil {

			// don't write out if we're just complaining about it being closed
			if perr, ok := err.(*serial.PortError); ok {
				if perr.Code() == serial.PortClosed {
					mc.ttyPort = nil
					return
				}
			}

			if errors.Is(err, syscall.EBADF) {
				return
			}

			logrus.Error("rx err: ", err.Error())
			return
		}

		for _, b := range buf[:n] {
			mc.ttyRx <- b
		}
		if n > 0 {
			logrus.Debugf("mcu rx: %x", buf[:n])
		}
	}
}

// Write will write the specified bytes to the microcontroller
func (mc *Microcontroller) Write(bs ...[]byte) (err error) {
	if !mc.IsOpen() {
		return ErrClosed
	}

	if len(bs) == 0 {
		panic("must provide at least one []byte")
	}

	for _, b := range bs {
		_, err = mc.ttyPort.Write(b)
		if err != nil {
			return
		}
		logrus.Debugf("mcu tx: %x", b)
	}

	return
}

// ReadN will read exactly N bytes from the rx chan
func (mc *Microcontroller) ReadN(n int, to time.Duration) ([]byte, error) {
	if !mc.IsOpen() {
		return nil, ErrClosed
	}

	bs := make([]byte, n)

	for i := 0; i < n; i++ {
		select {
		case <-time.After(to):
			return nil, ErrTimeout
		case b := <-mc.ttyRx:
			bs[i] = b
		}
	}

	return bs, nil
}
