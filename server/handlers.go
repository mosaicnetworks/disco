package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/disco"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup disco.Group
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the group title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &newGroup)
	if err != nil {
		fmt.Fprintf(w, "Error unmarshalling group: %v", err)
	}

	id, err := _groupRepo.SetGroup(&newGroup)
	if err != nil {
		fmt.Fprintf(w, "Error saving group: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func getOneGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	group, err := _groupRepo.GetGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error getting group %s: %v", groupID, err)
	}

	json.NewEncoder(w).Encode(group)
}

func getAllGroups(w http.ResponseWriter, r *http.Request) {
	allGroups, err := _groupRepo.GetAllGroups()
	if err != nil {
		fmt.Fprintf(w, "Error getting groups: %v", err)
	}
	json.NewEncoder(w).Encode(allGroups)
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	var updatedGroup disco.Group

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &updatedGroup)
	if err != nil {
		fmt.Fprintf(w, "Error parsing group: %v", err)
	}

	id, err := _groupRepo.SetGroup(&updatedGroup)
	if err != nil {
		fmt.Fprintf(w, "Error setting group: %v", err)
	}

	json.NewEncoder(w).Encode(id)
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	err := _groupRepo.DeleteGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error deleting group: %v", err)
	}

	fmt.Fprintf(w, "The group with ID %v has been deleted successfully", groupID)
}
