package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/disco/group"
	"github.com/mosaicnetworks/disco/server"
	"github.com/sirupsen/logrus"
)

/* TODO
- Test connecting with self-signed certificates and skip-verify
- Test all methods
*/

func TestSetGroup(t *testing.T) {

	// Init server and client

	server := server.NewDiscoServer(
		group.NewInmemGroupRepository(),
		"../test_data/cert.pem",
		"../test_data/key.pem",
		logrus.New().WithField("component", "disco-server"),
	)

	go server.Serve(
		"localhost:10443",
		"localhost:20443",
		"main",
	)

	time.Sleep(2 * time.Second)

	client, err := NewDiscoClient(
		"localhost:10443",
		"",
		true,
		logrus.New().WithField("component", "disco-client"),
	)

	if err != nil {
		t.Fatal(err)
	}

	// Insert group1

	group1 := group.NewGroup(
		"",
		"TestGroup1",
		"TestApp1",
		[]*peers.Peer{
			peers.NewPeer("pub1", "net1", "peer1"),
		},
	)

	group1ID, err := client.CreateGroup(*group1)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve group1 and check values

	retrievedGroup, err := client.GetGroupByID(group1ID)
	if err != nil {
		t.Fatal(err)
	}

	if group1.Name != retrievedGroup.Name {
		t.Fatalf("group Name should be %s, not %s", group1.Name, retrievedGroup.Name)
	}

	if group1.AppID != retrievedGroup.AppID {
		t.Fatalf("group AppID should be %s, not %s", group1.AppID, retrievedGroup.AppID)
	}

	if !reflect.DeepEqual(group1.Peers, retrievedGroup.Peers) {
		t.Fatalf("group Peers should be %#v, not %#v", group1.Peers, retrievedGroup.Peers)
	}

	if group1ID != retrievedGroup.ID {
		t.Fatalf("group ID should be %v, not %v", group1ID, retrievedGroup.ID)
	}

	if retrievedGroup.LastUpdated <= 0 {
		t.Fatalf("group LastUpdated should be positive")
	}

	// Insert another group from another app

	group2 := group.NewGroup(
		"",
		"TestGroup2",
		"TestApp2",
		[]*peers.Peer{
			peers.NewPeer("pub1", "net1", "peer1"),
		},
	)

	_, err = client.CreateGroup(*group2)
	if err != nil {
		t.Fatal(err)
	}

	// Get all groups

	allGroups, err := client.GetGroups("")
	if err != nil {
		t.Fatal(err)
	}

	if len(allGroups) != 2 {
		t.Fatalf("All groups should contain 2 groups, not %d", len(allGroups))
	}

	// Get App1 groups

	app1Groups, err := client.GetGroups("TestApp1")
	if err != nil {
		t.Fatal(err)
	}

	if len(app1Groups) != 1 {
		t.Fatalf("TestApp1 should contain 1 group, not %d", len(app1Groups))
	}

	// Get App2 groups

	app2Groups, err := client.GetGroups("TestApp2")
	if err != nil {
		t.Fatal(err)
	}

	if len(app2Groups) != 1 {
		t.Fatalf("TestApp2 should contain 1 group, not %d", len(app2Groups))
	}

	// Delete group 1

	err = client.DeleteGroup(group1ID)
	if err != nil {
		t.Fatal(err)
	}

	g, err := client.GetGroupByID(group1ID)
	if g != nil || err == nil {
		t.Fatalf("Retrieving deleted group should be return nil and error")
	}

}
