package commands

import (
	"fmt"
	"time"

	"github.com/mosaicnetworks/disco/group"
	"github.com/mosaicnetworks/disco/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var address = "0.0.0.0"
var discoPort = "1443"
var signalPort = "2443"
var icePort = "3478"
var iceUsername = "test"
var icePassword = "test"
var realm = "main"
var certFile = "cert.pem"
var keyFile = "key.pem"
var ttl = 5 * time.Minute
var ttlHeartbeat = 1 * time.Minute

func init() {
	RootCmd.Flags().StringVar(&address, "address", address, "Advertise address (use public address)")
	RootCmd.Flags().StringVar(&discoPort, "disco-port", discoPort, "Discovery API port")
	RootCmd.Flags().StringVar(&signalPort, "signal-port", signalPort, "WebRTC-Signaling port")
	RootCmd.Flags().StringVar(&icePort, "ice-port", icePort, "ICE server port")
	RootCmd.Flags().StringVar(&iceUsername, "ice-username", iceUsername, "ICE server userame. Only this user will be allowed to use the ICE server")
	RootCmd.Flags().StringVar(&icePassword, "ice-password", icePassword, "ICE server password corresponding to username")
	RootCmd.Flags().StringVar(&realm, "realm", realm, "Administrative routing domain within the WebRTC signaling")
	RootCmd.Flags().StringVar(&certFile, "cert-file", certFile, "File containing TLS certificate")
	RootCmd.Flags().StringVar(&keyFile, "key-file", keyFile, "File containing certificate key")
	RootCmd.Flags().DurationVar(&ttl, "ttl", ttl, "Group Time To Live, after which groups will be deleted")
	RootCmd.Flags().DurationVar(&ttlHeartbeat, "ttl-hearbeat", ttlHeartbeat, "Ticker frequency for checking group TTL")
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

	discoUrl := fmt.Sprintf("0.0.0.0:%s", discoPort)
	signalUrl := fmt.Sprintf("0.0.0.0:%s", signalPort)
	iceUrl := fmt.Sprintf("%s:%s", address, icePort)

	discoServer.Serve(
		discoUrl,
		signalUrl,
		iceUrl,
		iceUsername,
		icePassword,
		realm,
		ttl,
		ttlHeartbeat)

	return nil
}
