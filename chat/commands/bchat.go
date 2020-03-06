package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mosaicnetworks/babble/src/babble"
)

// newBChat returns a BChat node which is composed of a Babble node and an
// inmemory DummyApp proxy
func newBChat(config *BChatConfig, groupID string, moniker string) (*babble.Babble, error) {
	// sets babble configuration and app proxy
	babbleConfig := config.BabbleConfig(groupID, moniker)

	engine := babble.NewBabble(babbleConfig)

	if err := engine.Init(); err != nil {
		babbleConfig.Logger().Error("Cannot initialize engine:", err)
		return nil, err
	}

	return engine, nil
}

// waitForInterrupt listens for an interrupt signal and politely leaves the
// BChat (Babble leave request) before returning
func waitForInterrupt(engine *babble.Babble) {
	//Prepare sigCh to relay SIGINT and SIGTERM system calls
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	engine.Node.Leave()
}
