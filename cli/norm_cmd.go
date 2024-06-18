package main

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/odit-bit/soundfreak/fx"
	"github.com/odit-bit/soundfreak/wav"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	normCMD.Flags().Float64P("target", "t", -16.0, "target loudness (lufs) for normalizing audio content")
	viper.BindPFlag("target", normCMD.Flags().Lookup("target"))
}

func normArgValidate() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		if len(args) != 2 {
			return fmt.Errorf("requires exactly %d arg(s), received %d", 2, len(args))
		}

		return nil
	}
}

var normCMD = cobra.Command{
	Use:     "norm in out",
	Short:   "normalize to target lufs",
	Example: "norm ./stereo.wav ./stereo_norm.wav",
	Args:    normArgValidate(),
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer in.Close()
		if args[1] == "" {
			return fmt.Errorf("invalid output path")
		}
		target := viper.GetFloat64("target")
		return Normalize(in, args[1], target)
	},
}

func Normalize(r io.Reader, outFile string, target float64) error {
	wv, err := wav.New(r)
	if err != nil {
		return err
	}

	sampleF, err := wv.ReadFloat()
	if err != nil {
		return err
	}

	sampleCopy := make([]float64, len(sampleF))
	copy(sampleCopy, sampleF)

	//FIND LUFS
	lufs := fx.LUFS(sampleCopy, int(wv.Header.Frequence))
	fmt.Printf("loudness: %f\n", lufs)

	if lufs <= target {
		fmt.Println("loudness below the max target, no need adjustment")
		return nil
	}

	//adjust the gain of original sample
	gain := math.Pow(10, ((target - lufs) / 20))
	fmt.Printf("adjustment need: %f\n", gain)
	for i := range sampleF {
		sampleF[i] *= gain
	}

	// WRITE
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}

	if err := wav.WriteFloat64(out, int(wv.Header.Frequence), int(wv.Header.BitsPerSample), int(wv.Header.NumChannels), sampleF); err != nil {
		out.Close()
		os.RemoveAll(outFile)
		return err
	}
	return nil
}
