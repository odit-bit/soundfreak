package wav

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

func (wd *Wave) ReadSample() ([]int32, error) {
	if wd.isReaded {
		return nil, fmt.Errorf("sample has been readed")
	}
	wd.isReaded = true

	var b []int32
	br := bufio.NewReader(wd.r)
	peek := int(wd.Header.BitsPerSample / 8)
	end := wd.Header.DataSize
	switch peek {
	case 2:
		for i := 0; i < int(end); i += peek {
			sample16, err := br.Peek(peek)
			if err != nil {
				return nil, err
			}
			n := int32(binary.LittleEndian.Uint16(sample16))
			b = append(b, n)
			br.Discard(peek)
		}
		return b, nil

	case 3:
		for i := 0; i < int(end); i += peek {
			sample24, err := br.Peek(peek)
			if err != nil {
				return nil, err
			}
			n := fromBit24(sample24)
			b = append(b, n)
			br.Discard(peek)
		}
		return b, nil

	default:
		return nil, fmt.Errorf("bit depth not supported %d", wd.Header.BitsPerSample)
	}

}

func (wv *Wave) ReadFloat() ([]float64, error) {
	// if wv.isReaded {
	// 	return nil, fmt.Errorf("sample has been readed")
	// }
	// wv.isReaded = true

	sample32, err := wv.ReadSample()
	if err != nil {
		return nil, err
	}

	sampleF := make([]float64, len(sample32))
	var f float64
	for i, v := range sample32 {

		switch wv.Header.BitsPerSample {
		case 16:
			f = float64(v) / _MAX16_SAMPLE
		case 24:
			f = float64(v) / _MAX24_SAMPLE
		default:
			return nil, fmt.Errorf("invalid bit")
		}
		sampleF[i] = f
	}

	return sampleF, nil
}

// return sample as available channel
// error return only when samples has been read
func (wv *Wave) ReadChannels() ([][]float64, error) {
	sample64, err := wv.ReadFloat()
	if err != nil {
		return nil, err
	}

	ch := make([][]float64, wv.Header.NumChannels)
	if wv.Header.AudioFormat == 1 {
		//PCM
		for i := 0; i < len(ch); i++ {
			for j := 0; j < len(sample64); j += len(ch) {
				ch[i] = append(ch[i], sample64[j])
			}
		}
	}
	return ch, nil
}
