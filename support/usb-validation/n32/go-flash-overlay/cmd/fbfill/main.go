package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const (
	fbIOGETVSCREENINFO = 0x4600
	fbIOGETFSCREENINFO = 0x4602
)

type fbBitfield struct {
	Offset   uint32
	Length   uint32
	MsbRight uint32
}

type fbVarScreeninfo struct {
	Xres         uint32
	Yres         uint32
	XresVirtual  uint32
	YresVirtual  uint32
	Xoffset      uint32
	Yoffset      uint32
	BitsPerPixel uint32
	Grayscale    uint32
	Red          fbBitfield
	Green        fbBitfield
	Blue         fbBitfield
	Transp       fbBitfield
	Nonstd       uint32
	Activate     uint32
	Height       uint32
	Width        uint32
	AccelFlags   uint32
	Pixclock     uint32
	LeftMargin   uint32
	RightMargin  uint32
	UpperMargin  uint32
	LowerMargin  uint32
	HsyncLen     uint32
	VsyncLen     uint32
	Sync         uint32
	Vmode        uint32
	Rotate       uint32
	Colorspace   uint32
	Reserved     [4]uint32
}

type fbFixScreeninfo struct {
	ID         [16]byte
	SmemStart  uintptr
	SmemLen    uint32
	Type       uint32
	TypeAux    uint32
	Visual     uint32
	Xpanstep   uint16
	Ypanstep   uint16
	Ywrapstep  uint16
	LineLength uint32
	MmioStart  uintptr
	MmioLen    uint32
	Accel      uint32
	Caps       uint16
	Reserved   [2]uint16
}

func main() {
	var fb string
	var color string
	var bytesToWrite int
	var width int
	var height int
	var bpp int
	var order string
	flag.StringVar(&fb, "fb", "/dev/fb0", "framebuffer device")
	flag.StringVar(&color, "color", "red", "color: black, white, blue, red, yellow, or green")
	flag.IntVar(&bytesToWrite, "bytes", 0, "bytes to write; default uses existing framebuffer size")
	flag.IntVar(&width, "width", 0, "framebuffer width override")
	flag.IntVar(&height, "height", 0, "framebuffer height override")
	flag.IntVar(&bpp, "bpp", 0, "bits-per-pixel override")
	flag.StringVar(&order, "order", "", "byte order override: rgba, bgra, argb, abgr, rgb565, bgr565")
	flag.Parse()

	st, err := os.Stat(fb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f, err := os.OpenFile(fb, os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	info := detectFramebuffer(f)
	if width > 0 {
		info.width = width
	}
	if height > 0 {
		info.height = height
	}
	if bpp > 0 {
		info.bpp = bpp
	}
	if order != "" {
		info.order = order
	}
	if bytesToWrite <= 0 {
		bytesToWrite = info.byteLen
		if bytesToWrite <= 0 {
			bytesToWrite = int(st.Size())
		}
		if bytesToWrite <= 0 {
			bytesToWrite = 800 * 480 * 4
		}
	}

	pattern, err := patternFor(color, info.order)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	buf := make([]byte, 4096)
	for i := 0; i < len(buf); i += len(pattern) {
		copy(buf[i:], pattern)
	}
	remaining := bytesToWrite
	wrote := 0
	for remaining > 0 {
		n := len(buf)
		if remaining < n {
			n = remaining
		}
		wn, err := f.Write(buf[:n])
		wrote += wn
		remaining -= wn
		if err != nil {
			if wrote > 0 && errors.Is(err, syscall.ENOSPC) {
				return
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

type framebufferInfo struct {
	width   int
	height  int
	bpp     int
	byteLen int
	order   string
}

func detectFramebuffer(f *os.File) framebufferInfo {
	info := framebufferInfo{width: 800, height: 480, bpp: 32, byteLen: 800 * 480 * 4, order: "rgba"}
	var v fbVarScreeninfo
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), fbIOGETVSCREENINFO, uintptr(unsafe.Pointer(&v))); errno == 0 {
		if v.Xres > 0 {
			info.width = int(v.Xres)
		}
		if v.Yres > 0 {
			info.height = int(v.Yres)
		}
		if v.BitsPerPixel > 0 {
			info.bpp = int(v.BitsPerPixel)
		}
		info.order = orderFromBitfields(v)
	}
	var fix fbFixScreeninfo
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), fbIOGETFSCREENINFO, uintptr(unsafe.Pointer(&fix))); errno == 0 {
		if fix.LineLength > 0 && info.height > 0 {
			info.byteLen = int(fix.LineLength) * info.height
		}
	} else {
		info.byteLen = info.width * info.height * ((info.bpp + 7) / 8)
	}
	return info
}

func orderFromBitfields(v fbVarScreeninfo) string {
	if v.BitsPerPixel == 16 {
		if v.Red.Offset == 11 && v.Green.Offset == 5 && v.Blue.Offset == 0 {
			return "rgb565"
		}
		if v.Blue.Offset == 11 && v.Green.Offset == 5 && v.Red.Offset == 0 {
			return "bgr565"
		}
	}
	if v.BitsPerPixel == 32 || v.BitsPerPixel == 24 {
		switch {
		case v.Red.Offset == 0 && v.Green.Offset == 8 && v.Blue.Offset == 16:
			return "rgba"
		case v.Blue.Offset == 0 && v.Green.Offset == 8 && v.Red.Offset == 16:
			return "bgra"
		case v.Red.Offset == 8 && v.Green.Offset == 16 && v.Blue.Offset == 24:
			return "argb"
		case v.Blue.Offset == 8 && v.Green.Offset == 16 && v.Red.Offset == 24:
			return "abgr"
		}
	}
	return "rgba"
}

func patternFor(color string, order string) ([]byte, error) {
	r, g, b := byte(0), byte(0), byte(0)
	switch strings.ToLower(color) {
	case "black":
		// default zero values
	case "white":
		r, g, b = 0xff, 0xff, 0xff
	case "blue":
		b = 0xff
	case "red":
		r = 0xff
	case "yellow":
		r, g = 0xff, 0xff
	case "green":
		g = 0xff
	default:
		return nil, fmt.Errorf("unsupported color %q", color)
	}

	switch strings.ToLower(order) {
	case "rgba":
		return []byte{r, g, b, 0x00}, nil
	case "bgra":
		return []byte{b, g, r, 0x00}, nil
	case "argb":
		return []byte{0x00, r, g, b}, nil
	case "abgr":
		return []byte{0x00, b, g, r}, nil
	case "rgb565":
		v := uint16(r&0xf8)<<8 | uint16(g&0xfc)<<3 | uint16(b)>>3
		return []byte{byte(v), byte(v >> 8)}, nil
	case "bgr565":
		v := uint16(b&0xf8)<<8 | uint16(g&0xfc)<<3 | uint16(r)>>3
		return []byte{byte(v), byte(v >> 8)}, nil
	default:
		return nil, fmt.Errorf("unsupported byte order %q", order)
	}
}
