package main

import (
	"github.com/spf13/cobra"
)

func main() {

	app := cobra.Command{}
	app.AddCommand(&lufsCMD, &splitCMD, &normCMD)
	app.Execute()

}
