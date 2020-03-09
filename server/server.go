package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/babble/src/net/signal/wamp"
	"github.com/mosaicnetworks/disco/group"
)

type DiscoServer struct {
	repo group.GroupRepository
}

func NewDiscoServer(repo group.GroupRepository) *DiscoServer {
	return &DiscoServer{
		repo: repo,
	}
}

func (s *DiscoServer) Serve(discoAddr string, signalAddr string, realm string) {
	// XXX
	wampServer, err := wamp.NewServer(signalAddr, realm)
	if err != nil {
		log.Fatal(err)
	}
	go wampServer.Run()
	defer wampServer.Shutdown()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/group", s.CreateGroup).Methods("POST")
	router.HandleFunc("/groups", s.GetAllGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", s.GetOneGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", s.UpdateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", s.DeleteGroup).Methods("DELETE")
	log.Fatal(http.ListenAndServe(discoAddr, router))
	return
}

func (s *DiscoServer) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup group.Group
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the group title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &newGroup)
	if err != nil {
		fmt.Fprintf(w, "Error unmarshalling group: %v", err)
	}

	id, err := s.repo.SetGroup(&newGroup)
	if err != nil {
		fmt.Fprintf(w, "Error saving group: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func (s *DiscoServer) GetOneGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	group, err := s.repo.GetGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error getting group %s: %v", groupID, err)
	}

	json.NewEncoder(w).Encode(group)
}

func (s *DiscoServer) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	allGroups, err := s.repo.GetAllGroups()
	if err != nil {
		fmt.Fprintf(w, "Error getting groups: %v", err)
	}
	json.NewEncoder(w).Encode(allGroups)
}

func (s *DiscoServer) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var updatedGroup group.Group

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	err = json.Unmarshal(reqBody, &updatedGroup)
	if err != nil {
		fmt.Fprintf(w, "Error parsing group: %v", err)
	}

	id, err := s.repo.SetGroup(&updatedGroup)
	if err != nil {
		fmt.Fprintf(w, "Error setting group: %v", err)
	}

	json.NewEncoder(w).Encode(id)
}

func (s *DiscoServer) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	err := s.repo.DeleteGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error deleting group: %v", err)
	}

	fmt.Fprintf(w, "The group with ID %v has been deleted successfully", groupID)
}
