package wav

import (
	"fmt"
)

func fromBit24(sample24 []byte) int32 {
	_ = sample24[2]
	sample := int32(sample24[0]) | int32(sample24[1])<<8 | int32(sample24[2])<<16

	if sample&0x800000 != 0 {
		sample |= ^0xffffff
	}
	return sample

}

// convert max 4 byte size little-endian into int32
func toInt32(sample []byte) (int32, error) {
	var n int32
	switch len(sample) {
	case 2:
		//16bit
		n = int32(sample[0]) | int32(sample[1])<<8
	case 3:
		//24bit
		n = fromBit24(sample)

	case 4:
		//32bit
		n = int32(sample[0]) | int32(sample[1])<<8 | int32(sample[2])<<16 | int32(sample[3])<<24
	default:
		return 0, fmt.Errorf("expected 4 byte size")
	}

	return n, nil
}

func putUint24(b []byte, v int32) {
	_ = b[3]
	if (v & 0x800000) > 0 {
		v |= ^0xffffff
	}
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)

}
