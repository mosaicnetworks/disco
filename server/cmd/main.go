package main

import (
	"github.com/mosaicnetworks/disco/group"
	"github.com/mosaicnetworks/disco/server"
)

func main() {
	groupRepo := group.NewInmemGroupRepository()

	discoServer := server.NewDiscoServer(groupRepo)

	discoServer.Serve(":8080", ":9090", "office")
}
