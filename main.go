package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/babble/src/net/signal/wamp"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var routing = ":9090"
var realm = "office"
var discovery = ":8080"

func init() {
	RootCmd.Flags().StringVar(&routing, "routing", routing, "Routing Listen IP:Port")
	RootCmd.Flags().StringVar(&discovery, "discovery", discovery, "Discovery Listen IP:Port")
	viper.BindPFlags(RootCmd.Flags())
}

//RootCmd is the root command for the signaling server
var RootCmd = &cobra.Command{
	Use:   "signal",
	Short: "WebRTC signaling server using WebSockets",
	RunE:  runServer,
}

type group struct {
	GroupUID     string        `json:"GroupUID"`
	GroupName    string        `json:"GroupName"`
	AppID        string        `json:"AppID"`
	PubKey       string        `json:"PubKey"`
	LastUpdated  int64         `json:"LastUpdated"`
	Peers        []*peers.Peer `json:"Peers"`
	GenesisPeers []*peers.Peer `json:"InitialPeers"`
}

var allGroups = []group{
	{
		GroupUID:    "1",
		GroupName:   "Introduction to Golang",
		PubKey:      "ABCDEF",
		AppID:       "io.mosaicnetworks.sample",
		LastUpdated: 1583418648,
		Peers: peers.NewPeerSet([]*peers.Peer{
			peers.NewPeer("XXX", "Peer0Addr", "Peer0")}).Peers,
		GenesisPeers: peers.NewPeerSet([]*peers.Peer{
			peers.NewPeer("XXX", "Peer0Addr", "Peer0")}).Peers,
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup group
	log.Print("createGroup")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Kindly enter data with the group title and description only in order to update")
		fmt.Fprintf(w, "Kindly enter data with the group title and description only in order to update")
	}

	log.Print(string(reqBody))

	err = json.Unmarshal(reqBody, &newGroup)
	if err != nil {
		log.Print("Error unmarshalling")
		fmt.Fprintf(w, "Error unmarshalling group: %v", err)
	}

	// Overwrite the last updated time with current server time
	newGroup.LastUpdated = time.Now().Unix()

	allGroups = append(allGroups, newGroup)
	//	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newGroup)
}

func getOneGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	for _, singleGroup := range allGroups {
		if singleGroup.GroupUID == groupID {
			json.NewEncoder(w).Encode(singleGroup)
		}
	}
}

func getAllGroups(w http.ResponseWriter, r *http.Request) {
	log.Print("getAllGroups")
	json.NewEncoder(w).Encode(allGroups)
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]
	var updatedGroup group

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &updatedGroup)
	if err != nil {
		fmt.Fprintf(w, "Error parsing group: %v", err)
	}

	for i, singleGroup := range allGroups {
		if singleGroup.GroupUID == groupID {
			singleGroup.GroupName = updatedGroup.GroupName
			singleGroup.PubKey = updatedGroup.PubKey
			singleGroup.AppID = updatedGroup.AppID
			singleGroup.LastUpdated = time.Now().Unix()
			singleGroup.Peers = updatedGroup.Peers
			singleGroup.GenesisPeers = updatedGroup.GenesisPeers

			allGroups = append(allGroups[:i], singleGroup)
			json.NewEncoder(w).Encode(singleGroup)
		}
	}
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	for i, singleGroup := range allGroups {
		if singleGroup.GroupUID == groupID {
			allGroups = append(allGroups[:i], allGroups[i+1:]...)
			fmt.Fprintf(w, "The group with ID %v has been deleted successfully", groupID)
		}
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	server, err := wamp.NewServer(routing, realm)
	if err != nil {
		log.Fatal(err)
	}

	go server.Run()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/group", createGroup).Methods("POST")
	router.HandleFunc("/groups", getAllGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", getOneGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", updateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")

	log.Print("Starting Discovery Server")
	go log.Fatal(http.ListenAndServe(discovery, router))

	//Prepare sigCh to relay SIGINT and SIGTERM system calls
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	server.Shutdown()

	return nil
}

func main() {
	//Do not print usage when error occurs
	RootCmd.SilenceUsage = true

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
