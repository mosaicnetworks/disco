package group

import "github.com/mosaicnetworks/babble/src/peers"

// Group represents a Babble group of peers
type Group struct {
	ID             string
	Name           string
	AppID          string
	PubKey         string
	LastUpdated    int64
	LastBlockIndex int64
	Peers          []*peers.Peer
	GenesisPeers   []*peers.Peer
}

// NewGroup generates a new Group
func NewGroup(id string, name string, appID string, peers []*peers.Peer) *Group {
	return &Group{
		ID:           id,
		Name:         name,
		AppID:        appID,
		Peers:        peers,
		GenesisPeers: peers,
	}
}
