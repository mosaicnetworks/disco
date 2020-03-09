package group

import "github.com/mosaicnetworks/babble/src/peers"

// Group represents a Babble group of peers
type Group struct {
	ID    string        `json:"id"`
	Title string        `json:"title"`
	Peers []*peers.Peer `json:"peers"`
}

// NewGroup generates a new Group
func NewGroup(id string, title string, peers []*peers.Peer) *Group {
	return &Group{
		ID:    id,
		Title: title,
		Peers: peers,
	}
}
