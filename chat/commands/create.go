package commands

import (
	"fmt"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/disco"
	"github.com/mosaicnetworks/disco/client"
	"github.com/spf13/cobra"
)

// NewCreateCmd produces a CreateCmd which which stats a new BChat group and
// advertises it in the discovery service
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Start a new chat channel",
		RunE:  create,
	}
	return cmd
}

// Initialise and start a new BChat node. Advertise the new group in the
// discovery service and wait for an interrup signal before politely leaving the
// BChat.
func create(cmd *cobra.Command, args []string) error {
	title := promptTitle()
	moniker := promptMoniker()

	configManager := NewConfigManager(_config.BaseDir)

	// Create a group ID and dump private key and peers.json in a corresponding
	// config directory
	groupID, err := configManager.CreateNew(moniker)
	if err != nil {
		return err
	}

	engine, err := newBChat(_config, groupID, moniker)
	if err != nil {
		return err
	}

	go engine.Run()
	defer waitForInterrupt(engine)

	err = advertiseGroup(groupID, title, engine.Node.GetPeers())
	if err != nil {
		return err
	}

	chat()

	return nil
}

func advertiseGroup(groupID string, title string, peers []*peers.Peer) error {
	newGroup := disco.NewGroup(
		groupID,
		title,
		peers, // equal to genesis peers
	)

	discoClient := client.NewDiscoClient(_config.DiscoAddr)

	id, err := discoClient.CreateGroup(*newGroup)
	if err != nil {
		return fmt.Errorf("Error creating group: %v", err)
	}

	fmt.Printf("Group ID: %v\n", id)

	return nil
}
