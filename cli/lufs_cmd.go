package main

import (
	"fmt"
	"io"
	"os"

	"github.com/odit-bit/soundfreak/fx"
	"github.com/odit-bit/soundfreak/wav"
	"github.com/spf13/cobra"
)

var lufsCMD = cobra.Command{
	Use:     "lufs ./path/to/audio.wav",
	Short:   "use for find the lufs of audio content",
	Long:    "",
	Example: "lufs ./myAwesomeAudio.wav",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		in, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer in.Close()

		if err := lufs(in); err != nil {
			return err
		}

		return nil
	},
}

func lufs(r io.Reader) error {

	wv, err := wav.New(r)
	if err != nil {
		return err
	}

	sampleF, err := wv.ReadFloat()
	if err != nil {
		return err
	}

	lufs := fx.LUFS(sampleF, int(wv.Header.Frequence))
	fmt.Printf("loudness: %f\n", lufs)

	return nil
}
