package commands

import (
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
)

// BChat is a group-chat application that uses Babble. The structure wraps an
// InmemDummyClient that acts as the app-side that plugs into Babble Babble
// consensus using an inmemory proxy.
type BChat struct {
	GroupID string
	Moniker string
	Proxy   *dummy.InmemDummyClient
	Engine  *babble.Babble
	Gui     BChatGui
}

// NewBChat returns a BChat node which is composed of a Babble node and an
// inmemory DummyApp proxy
func NewBChat(config *BChatConfig, groupID string, moniker string) (*BChat, error) {
	// produce babble configuration based on BChatConfig and default values
	babbleConfig := config.BabbleConfig(groupID, moniker)

	// set the app proxy
	proxy := dummy.NewInmemDummyClient(babbleConfig.Logger())
	babbleConfig.Proxy = proxy

	engine := babble.NewBabble(babbleConfig)

	if err := engine.Init(); err != nil {
		babbleConfig.Logger().Error("Cannot initialize engine:", err)
		return nil, err
	}

	bchatTui := NewBChatGui(proxy, moniker)

	return &BChat{
		GroupID: groupID,
		Moniker: moniker,
		Proxy:   proxy,
		Engine:  engine,
		Gui:     bchatTui,
	}, nil
}

// Run Starts Babble and the GUI. It keeps running until ctrl-C is pressed.
func (b *BChat) Run() error {
	go b.Engine.Run()
	defer b.Engine.Node.Leave()

	return b.Gui.Loop()
}
