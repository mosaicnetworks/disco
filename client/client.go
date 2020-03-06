package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mosaicnetworks/disco"
)

// DiscoClient is a client for the Discovery API
type DiscoClient struct {
	url string
}

// NewDiscoClient creates a new DiscoClient for a server hosted at the provided
// url
func NewDiscoClient(url string) *DiscoClient {
	return &DiscoClient{
		url: url,
	}
}

// GetAllGroups returs all groups in a map where the key is the ID of the group
// and the value is a pointer to the corresponding Group object.
func (c *DiscoClient) GetAllGroups() (map[string]*disco.Group, error) {
	path := fmt.Sprintf("%s/groups", c.url)
	fmt.Println("path: ", path)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var allGroups map[string]*disco.Group
	err = json.Unmarshal(body, &allGroups)
	if err != nil {
		return nil, fmt.Errorf("Error parsing groups: %v", err)
	}

	return allGroups, nil
}

// GetGroupByID gets a single group by ID
func (c *DiscoClient) GetGroupByID(id string) (*disco.Group, error) {
	path := fmt.Sprintf("%s/groups/%s", c.url, id)
	fmt.Println("path: ", path)

	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var group *disco.Group
	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, fmt.Errorf("Error parsing group: %v", err)
	}

	return group, nil
}

// CreateGroup adds a group to the discovery server. The group's ID field should
// be empty as it will be set by the server.
func (c *DiscoClient) CreateGroup(group disco.Group) (string, error) {
	path := fmt.Sprintf("%s/group", c.url)
	fmt.Println("path: ", path)

	jsonValue, err := json.Marshal(group)
	if err != nil {
		return "", fmt.Errorf("Error marshalling group: %v", err)
	}

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var id string
	err = json.Unmarshal(body, &id)
	if err != nil {
		return "", fmt.Errorf("Error parsing id: %v", err)
	}

	return id, nil
}
