package commands

import (
	"fmt"
	"os"
	"path"

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

	title := promptTitle()
	moniker := promptMoniker()

	configManager := NewConfigManager(_config.BaseDir)

	// Create a group ID and dump private key and peers.json in a corresponding
	// config directory
	groupID, err := configManager.CreateNew(moniker)
	if err != nil {
		return err
	}

	bchat, err := newBChat(_config, groupID, moniker)
	if err != nil {
		return err
	}

	go bchat.Engine.Run()
	defer waitForInterrupt(bchat.Engine)

	err = advertiseGroup(groupID, title, bchat.Engine.Node.GetPeers())
	if err != nil {
		return err
	}

	chat(
		bchat.Proxy.SubmitCh(),
		path.Join(_config.BaseDir, groupID, "info.log"),
	) //XXX

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
