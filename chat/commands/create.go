package commands

import (
	"fmt"
	"os"

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
	// Prompt title of new group
	title := promptTitle()

	// Prompt moniker of user in group
	moniker := promptMoniker()

	// Create a group ID and dump private key and peers.json in a corresponding
	// config directory
	configManager := NewConfigManager(_config.BaseDir)
	groupID, err := configManager.CreateNew(moniker)
	if err != nil {
		return err
	}

	// XXX
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	_, wf, _ := os.Pipe()

	defer func() {
		wf.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	os.Stdout = wf
	os.Stderr = wf
	// end XXX

	bchat, err := NewBChat(_config, groupID, moniker)
	if err != nil {
		return err
	}

	err = advertiseGroup(groupID, title, bchat.Engine.Node.GetPeers())
	if err != nil {
		return err
	}

	bchat.Run()

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
