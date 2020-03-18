package commands

import (
	"github.com/mosaicnetworks/disco/group"
	"github.com/mosaicnetworks/disco/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var discoUrl = ":1443"
var signalUrl = ":2443"
var realm = "office"
var certFile = "cert.pem"
var keyFile = "key.pem"

func init() {
	RootCmd.Flags().StringVar(&discoUrl, "disco-url", discoUrl, "API IP:Port")
	RootCmd.Flags().StringVar(&signalUrl, "signal-url", signalUrl, "Signal IP:Port")
	RootCmd.Flags().StringVar(&realm, "realm", realm, "Administrative routing domain within the WebRTC signaling")
	RootCmd.Flags().StringVar(&certFile, "cert-file", certFile, "File containing TLS certificate")
	RootCmd.Flags().StringVar(&keyFile, "key-file", keyFile, "File containing certificate key")
	viper.BindPFlags(RootCmd.Flags())
}

//RootCmd is the root command for the disco server
var RootCmd = &cobra.Command{
	Use:   "disco",
	Short: "Discovery service for Babble",
	RunE:  runServer,
}

// runServer starts the disco server and waits for a SIGINT or SIGTERM
func runServer(cmd *cobra.Command, args []string) error {
	groupRepo := group.NewInmemGroupRepository()

	discoServer := server.NewDiscoServer(groupRepo,
		certFile,
		keyFile,
		logrus.New().WithField("component", "disco-server"))

	discoServer.Serve(discoUrl, signalUrl, "office")

	return nil
}
