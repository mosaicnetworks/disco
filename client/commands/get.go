package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mosaicnetworks/disco"
	"github.com/spf13/cobra"
)

// NewGetCmd produces a GetCmd which retrieves groups from the discovery server
func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get groups",
		RunE:  get,
	}

	return cmd
}

func get(cmd *cobra.Command, args []string) error {

	path := fmt.Sprintf("%s/groups", _url)
	fmt.Println("path: ", path)

	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var allGroups []disco.Group
	err = json.Unmarshal(body, &allGroups)
	if err != nil {
		return fmt.Errorf("Error parsing groups: %v", err)
	}

	fmt.Printf("%#v\n", allGroups)

	return nil
}
