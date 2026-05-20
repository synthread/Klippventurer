package flash

import (
	"golang.org/x/exp/constraints"
)

// checksum will create a STM-compatible XOR-based checksum of the provided data
func checksum(bs []byte) byte {
	if len(bs) == 0 {
		return 0x00
	}
	if len(bs) == 1 {
		return bs[0]
	}
	s := bs[0]
	for i := 1; i < len(bs); i++ {
		s ^= bs[i]
	}
	return s
}

// min will return the minimum of the two values
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
