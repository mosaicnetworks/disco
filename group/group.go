package group

import "github.com/mosaicnetworks/babble/src/peers"

// Group represents a Babble group of peers
type Group struct {
	GroupUID     string        `json:"GroupUID"`
	GroupName    string        `json:"GroupName"`
	AppID        string        `json:"AppID"`
	PubKey       string        `json:"PubKey"`
	LastUpdated  int64         `json:"LastUpdated"`
	Peers        []*peers.Peer `json:"Peers"`
	GenesisPeers []*peers.Peer `json:"InitialPeers"`
}

// NewGroup generates a new Group
func NewGroup(id string, name string, appID string, peers []*peers.Peer) *Group {
	return &Group{
		GroupUID:     id,
		GroupName:    name,
		AppID:        appID,
		Peers:        peers,
		GenesisPeers: peers,
	}
}
