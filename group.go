package disco

import "github.com/mosaicnetworks/babble/src/peers"

// Group represents a Babble group of peers
type Group struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Peers       []*peers.Peer `json:"peers"`
}

// NewGroup generates a new Group and leaves the ID field empty because IDs are
// controlled by GroupRepositories and assigned upon insertion.
func NewGroup(title string, description string, peers []*peers.Peer) *Group {
	return &Group{
		Title:       title,
		Description: description,
		Peers:       peers,
	}
}
