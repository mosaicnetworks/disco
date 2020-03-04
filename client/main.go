package main

import (
	_ "net/http/pprof"
	"os"

	cmd "github.com/mosaicnetworks/disco/client/commands"
)

func main() {
	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewGetCmd(),
		cmd.NewCreateCmd(),
	)

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
