package main

import (
	"io"
	"os"

	"github.com/odit-bit/soundfreak/wav"
)

func convert(r io.Reader, bitDepth int, outFile string) error {
	wv, err := wav.New(r)
	if err != nil {
		return err
	}

	samples, err := wv.ReadSample()
	if err != nil {
		return err
	}

	// WRITE
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := wav.WriteInt32(out, int(wv.Header.Frequence), bitDepth, int(wv.Header.NumChannels), samples); err != nil {
		os.RemoveAll(outFile)
		return err
	}
	return nil
}
