package commands

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	_config = NewDefaultBChatConfig()
)

func init() {
	RootCmd.PersistentFlags().String("basedir", _config.BaseDir, "Top-level directory for configuration and data")
	RootCmd.PersistentFlags().String("disco-addr", _config.DiscoAddr, "Address of the discovery server")
	RootCmd.PersistentFlags().String("signal-addr", _config.SignalAddr, "Address of the webrtc signaling server")
	RootCmd.PersistentFlags().String("log", _config.LogLevel, "debug, info, warn, error, fatal, panic")
}

// RootCmd is the root command for the Babble Chat Client
var RootCmd = &cobra.Command{
	Use:              "bchat",
	Short:            "Babble Chat Client",
	TraverseChildren: true,
	PreRunE:          loadConfig,
}

/*******************************************************************************
* CONFIG
*******************************************************************************/

func loadConfig(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	// first unmarshal to read from CLI flags
	return viper.Unmarshal(_config)
}

func logLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
