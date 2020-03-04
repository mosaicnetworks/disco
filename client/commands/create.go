package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/disco"
	"github.com/spf13/cobra"
)

// NewCreateCmd produces a CreateCmd which posts a new group to the disco server
func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create new group",
		RunE:  create,
	}

	return cmd
}

func create(cmd *cobra.Command, args []string) error {

	path := fmt.Sprintf("%s/group", _url)
	fmt.Println("path: ", path)

	jsonValue, _ := json.Marshal(disco.Group{
		ID:          "10",
		Title:       "Office",
		Description: "This is another group",
		Peers: []*peers.Peer{
			peers.NewPeer("XXX", "alice@locahost", "alice"),
			peers.NewPeer("YYY", "bob@locahost", "bob"),
			peers.NewPeer("ZZZ", "charlie@locahost", "charlie"),
		},
	})

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var group disco.Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		return fmt.Errorf("Error parsing group: %v", err)
	}

	fmt.Printf("%#v\n", group)

	return nil
}
