package commands

import (
	"fmt"

	"github.com/mosaicnetworks/disco/client"
	"github.com/spf13/cobra"
)

// NewListCmd produces a ListCmd which retrieves all the groups from the
// discovery server, and displays them.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE:  list,
	}

	return cmd
}

func list(cmd *cobra.Command, args []string) error {
	discoClient := client.NewDiscoClient(_config.DiscoAddr)

	allGroups, err := discoClient.GetAllGroups()
	if err != nil {
		return fmt.Errorf("Error fetching groups: %v", err)
	}

	for id, g := range allGroups {
		fmt.Printf("%v %v\n", id, g.Title)
	}

	return nil
}
