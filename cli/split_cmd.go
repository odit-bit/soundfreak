package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/odit-bit/soundfreak/wav"
	"github.com/spf13/cobra"
)

var splitCMD = cobra.Command{
	Use:     "split input",
	Short:   "split multichannel wav into multi single channel wav",
	Example: "split ./stereo.wav",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return split(args[0])
	},
}

func split(inputFile string) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// read wav file
	wv, err := wav.New(f)
	if err != nil {
		return err
	}
	//split the channels
	chs, err := wv.ReadChannels()
	if err != nil {
		return err
	}

	//write the channels
	dir, file := filepath.Split(f.Name())
	split := strings.Split(file, ".")
	ext := filepath.Ext(file)
	fmt.Println(dir, file)
	var wg sync.WaitGroup
	for i, c := range chs {
		wg.Add(1)
		go func(i int, c []float64) {
			defer wg.Done()
			outpath := fmt.Sprintf("%s%s_%d.%s", dir, split[0], i, ext)
			out, err := os.Create(outpath)
			if err != nil {
				log.Println(err)
				return
			}
			if err := wav.WriteFloat64(out, int(wv.Header.Frequence), int(wv.Header.BitsPerSample), 1, c); err != nil {
				log.Println(err)
				return
			}
		}(i, c)
	}

	wg.Wait()
	return nil
}
