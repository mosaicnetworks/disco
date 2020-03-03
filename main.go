package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/babble/src/peers"
)

type group struct {
	ID          string         `json:"ID"`
	Title       string         `json:"Title"`
	Description string         `json:"Description"`
	Peers       *peers.PeerSet `json:"Peers"`
}

var allGroups = []group{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
		Peers: peers.NewPeerSet([]*peers.Peer{
			peers.NewPeer("XXX", "Peer0Addr", "Peer0"),
		}),
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup group
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the group title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &newGroup)
	if err != nil {
		fmt.Fprintf(w, "Error unmarshalling group: %v", err)
	}
	allGroups = append(allGroups, newGroup)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newGroup)
}

func getOneGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	for _, singleGroup := range allGroups {
		if singleGroup.ID == groupID {
			json.NewEncoder(w).Encode(singleGroup)
		}
	}
}

func getAllGroups(w http.ResponseWriter, r *http.Request) {
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
		if singleGroup.ID == groupID {
			singleGroup.Title = updatedGroup.Title
			singleGroup.Description = updatedGroup.Description
			singleGroup.Peers = updatedGroup.Peers

			allGroups = append(allGroups[:i], singleGroup)
			json.NewEncoder(w).Encode(singleGroup)
		}
	}
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	for i, singleGroup := range allGroups {
		if singleGroup.ID == groupID {
			allGroups = append(allGroups[:i], allGroups[i+1:]...)
			fmt.Fprintf(w, "The group with ID %v has been deleted successfully", groupID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/group", createGroup).Methods("POST")
	router.HandleFunc("/groups", getAllGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", getOneGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", updateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
