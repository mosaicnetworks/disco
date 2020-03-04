package commands

import (
	"github.com/spf13/cobra"
)

var _url = "http://localhost:8080"

// RootCmd is the root command for the CLI discovery client
var RootCmd = &cobra.Command{
	Use:              "disco-client",
	Short:            "disco client",
	TraverseChildren: true,
}
