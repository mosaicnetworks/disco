package commands

import (
	"fmt"
	"os"

	"github.com/mosaicnetworks/disco/client"
	"github.com/spf13/cobra"
)

func NewJoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join [group-id]",
		Short: "Join a channel",
		RunE:  join,
	}
	return cmd
}

func join(cmd *cobra.Command, args []string) error {
	discoClient := client.NewDiscoClient(_config.DiscoAddr)

	allGroups, err := discoClient.GetAllGroups()
	if err != nil {
		return fmt.Errorf("Error fetching groups: %v", err)
	}

	selectedGroup := selectGroup(allGroups)

	moniker := promptMoniker()

	configManager := NewConfigManager(_config.BaseDir)
	err = configManager.CreateForJoin(selectedGroup)
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

	bchat, err := NewBChat(_config, selectedGroup.ID, moniker)
	if err != nil {
		return err
	}

	bchat.Run()

	return nil
}
