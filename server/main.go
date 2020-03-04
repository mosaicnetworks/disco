package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/disco"
)

// groupRepo is a global variable giving access to a database of groups
var _groupRepo disco.GroupRepository

func main() {
	// Create and populate and inmem group repo with a fake group
	_groupRepo = disco.NewInmemGroupRepository()
	_groupRepo.SetGroup(
		disco.NewGroup(
			"Group1",
			"Useless Group",
			[]*peers.Peer{
				peers.NewPeer("XXX", "alice@localhost", "Alice"),
				peers.NewPeer("YYY", "bob@localhost", "Bob"),
				peers.NewPeer("ZZZ", "charlie@localhost", "Charlie"),
			},
		))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/group", createGroup).Methods("POST")
	router.HandleFunc("/groups", getAllGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", getOneGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", updateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", deleteGroup).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
