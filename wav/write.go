package wav

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

func WriteFloat64(w io.Writer, frequence, bitsPerSample, ch int, sample []float64) error {
	bw := bufio.NewWriter(w)

	//write header
	dataSize := len(sample) * (bitsPerSample / 8)
	if err := writeHeader(bw, frequence, bitsPerSample, ch, uint32(dataSize)); err != nil {
		return err
	}

	//write data
	switch bitsPerSample {
	case 16:
		for _, f := range sample {
			//upscale float64 into larger bit integer than int16
			// in this case we use 24bit portion of int32 before downsacle into int16 (>> 8)
			n32 := int32(f * _MAX24_SAMPLE)
			n := int16(n32 >> 8)
			if err := binary.Write(bw, binary.LittleEndian, n); err != nil {
				return err
			}
		}
		return bw.Flush()

	case 24:
		buf := make([]byte, 4)
		for _, f := range sample {
			n := int32(f * _MAX24_SAMPLE)
			// n := int32(math.Round(f * _MAX24_SAMPLE))
			putUint24(buf, n)
			if err := binary.Write(bw, binary.LittleEndian, buf[:len(buf)-1]); err != nil {
				return err
			}
		}
		return bw.Flush()

	default:
		return fmt.Errorf("unimplemented bits %d", bitsPerSample)
	}
}

func WriteInt32(w io.Writer, frequence, bitsPerSample, ch int, sample []int32) error {
	bw := bufio.NewWriter(w)

	dataSize := len(sample) * (bitsPerSample / 8)
	if err := writeHeader(bw, frequence, bitsPerSample, ch, uint32(dataSize)); err != nil {
		return err
	}

	switch bitsPerSample {
	case 16:
		for _, v := range sample {
			sample16 := int16(v >> 8)
			if err := binary.Write(bw, binary.LittleEndian, sample16); err != nil {
				return err
			}
		}

	case 24:
		buf := make([]byte, 3)
		for _, v := range sample {
			putUint24(buf, v)
			if err := binary.Write(bw, binary.LittleEndian, buf); err != nil {
				return err
			}
		}
	}

	return bw.Flush()
}

func writeHeader(w io.Writer, frequence, bitsPerSample, ch int, dataSize uint32) error {

	bytePerBloc := ch * bitsPerSample / 8
	fileSize := 36 + dataSize
	header := Header{
		FileType:   _RIFF,
		FileSize:   uint32(fileSize),
		FileFormat: _WAVE,

		//
		FormatBlocID: _FMT,
		BlocSize:     16,
		AudioFormat:  _PCM,
		NumChannels:  uint16(ch),
		Frequence:    uint32(frequence),

		BytePerSec:    uint32(frequence * bytePerBloc),
		BytePerBloc:   uint16(bytePerBloc),
		BitsPerSample: uint16(bitsPerSample),

		//
		DataBloc: _SAMPLE_BLOC,
		DataSize: uint32(dataSize), //update
	}

	// slog.Info("wav", "info", header)

	if err := binary.Write(w, binary.LittleEndian, header); err != nil {
		return err
	}
	return nil
}
