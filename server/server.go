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
	"github.com/sirupsen/logrus"
)

// DiscoServer is a peer-discovery and webrtc-signaling service for Babble.
// Peer-discovery enables users to advertise groups that other people can join
// and is exposed over a regular HTTP REST API.
// WebRTC-signaling enables users to exchange connection metadata (SDP) to
// create direct p2p connections, and relies on the WAMP protocol which is
// basically RPC over web-sockets.
type DiscoServer struct {
	repo     group.GroupRepository
	certFile string
	keyFile  string
	logger   *logrus.Entry
}

// NewDiscoServer instantiates a new DiscoServer with a GroupRepository.
func NewDiscoServer(
	repo group.GroupRepository,
	certFile string,
	keyFile string,
	logger *logrus.Entry,
) *DiscoServer {

	return &DiscoServer{
		repo:     repo,
		certFile: certFile,
		keyFile:  keyFile,
		logger:   logger,
	}
}

// Serve starts the peer-discovery and signaling services
func (s *DiscoServer) Serve(discoAddr string, signalAddr string, realm string) {
	wampServer, err := wamp.NewServer(signalAddr,
		realm,
		s.certFile,
		s.keyFile,
		s.logger)
	if err != nil {
		log.Fatal(err)
	}
	go wampServer.Run()
	defer wampServer.Shutdown()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/group", s.createGroup).Methods("POST")
	router.HandleFunc("/groups", s.getGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", s.getGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", s.updateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", s.deleteGroup).Methods("DELETE")
	log.Fatal(http.ListenAndServeTLS(discoAddr, s.certFile, s.keyFile, router))
	return
}

func (s *DiscoServer) createGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup group.Group
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading request body: %v", err)
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

func (s *DiscoServer) getGroups(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("app-id")

	groups := make(map[string]*group.Group)
	var err error

	if appID == "" {
		groups, err = s.repo.GetAllGroups()
	} else {
		groups, err = s.repo.GetAllGroupsByAppID(appID)
	}

	if err != nil {
		fmt.Fprintf(w, "Error getting groups: %v", err)
	}

	json.NewEncoder(w).Encode(groups)
}

func (s *DiscoServer) getGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	group, err := s.repo.GetGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error getting group %s: %v", groupID, err)
	}

	json.NewEncoder(w).Encode(group)
}

func (s *DiscoServer) updateGroup(w http.ResponseWriter, r *http.Request) {
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

func (s *DiscoServer) deleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := mux.Vars(r)["id"]

	err := s.repo.DeleteGroup(groupID)
	if err != nil {
		fmt.Fprintf(w, "Error deleting group: %v", err)
	}

	fmt.Fprintf(w, "The group with ID %v has been deleted successfully", groupID)
}
