package wav

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	_MAX24_SAMPLE = 1<<23 - 1
	_MAX16_SAMPLE = 1<<15 - 1
)

const (
	_PCM     = uint16(1)
	_IEEE754 = uint16(3)
	_48KHz   = 48000
	_441KHz  = 44100
)

var (
	_RIFF        = [4]byte{0x52, 0x49, 0x46, 0x46}
	_WAVE        = [4]byte{0x57, 0x41, 0x56, 0x45}
	_FMT         = [4]byte{0x66, 0x6D, 0x74, 0x20}
	_SAMPLE_BLOC = [4]byte{0x64, 0x61, 0x74, 0x61}
)

type Header struct {
	FileType   [4]byte // "RIFF"
	FileSize   uint32
	FileFormat [4]byte // "WAVE"

	//format pointer
	FormatBlocID  [4]byte // "fmt "
	BlocSize      uint32
	AudioFormat   uint16 //pcm=1 or float=3
	NumChannels   uint16
	Frequence     uint32 // hz frequency
	BytePerSec    uint32 // byte/s
	BytePerBloc   uint16
	BitsPerSample uint16

	//sampled data pointer
	DataBloc [4]byte // "data"
	DataSize uint32
}

func readHeader(r io.Reader) (*Header, error) {
	wh := Header{}
	if err := binary.Read(r, binary.LittleEndian, &wh); err != nil {
		return nil, err
	}
	switch wh.BitsPerSample {
	case 8, 16, 24, 32:
	default:
		return nil, fmt.Errorf("unsupported bit %d ", wh.BitsPerSample)
	}
	return &wh, nil
}

type Wave struct {
	Header *Header

	MaxSample int64
	isReaded  bool
	r         io.Reader
}

// Decode decode wav binary format
func New(r io.Reader) (*Wave, error) {

	//read header
	h, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	// maxSample = data / (bitdepth/8)
	maxSample := h.DataSize / (uint32(h.BitsPerSample) / 8)

	lr := io.LimitReader(r, int64(h.DataSize))
	wd := &Wave{
		Header:    h,
		MaxSample: int64(maxSample),
		r:         lr,
	}

	return wd, nil

}

// type PluginFunc func(sample int64, hdr *Header) int64

// func (wv *Wave) Render(w io.Writer, plugins ...PluginFunc) error {
// 	if err := binary.Write(w, binary.LittleEndian, wv.header); err != nil {
// 		return err
// 	}

// 	buf := make([]byte, 3)
// 	switch wv.header.BitsPerSample {
// 	case 24:
// 		for i := 0; i < len(wv.data); i += 3 {
// 			sample := int32(wv.data[i]) | int32(wv.data[i+1])<<8 | int32(wv.data[i+2])<<16
// 			if sample&0x800000 != 0 {
// 				sample |= ^0xffffff
// 			}
// 			for _, plugin := range plugins {
// 				sample = int32(plugin(int64(sample), wv.header))
// 			}
// 			PutUint24(buf, sample)
// 			binary.Write(w, binary.LittleEndian, buf)
// 		}
// 	}
// 	return nil
// }
