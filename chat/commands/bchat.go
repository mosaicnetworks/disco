package commands

import (
	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
)

type BChat struct {
	GroupID string
	Moniker string
	Proxy   *dummy.InmemDummyClient
	Engine  *babble.Babble
}

// NewBChat returns a BChat node which is composed of a Babble node and an
// inmemory DummyApp proxy
func NewBChat(config *BChatConfig, groupID string, moniker string) (*BChat, error) {
	// sets babble configuration and app proxy
	babbleConfig := config.BabbleConfig(groupID, moniker)

	proxy := dummy.NewInmemDummyClient(babbleConfig.Logger())

	babbleConfig.Proxy = proxy

	engine := babble.NewBabble(babbleConfig)

	if err := engine.Init(); err != nil {
		babbleConfig.Logger().Error("Cannot initialize engine:", err)
		return nil, err
	}

	return &BChat{
		GroupID: groupID,
		Moniker: moniker,
		Proxy:   proxy,
		Engine:  engine,
	}, nil
}

func (b *BChat) Run() error {
	go b.Engine.Run()
	defer b.Engine.Node.Leave()

	bchatTui := NewBChatTui(b.Proxy, b.Moniker)
	return bchatTui.Loop()
}
