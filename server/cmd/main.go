package main

import (
	"github.com/mosaicnetworks/disco/group"
	"github.com/mosaicnetworks/disco/server"
)

func main() {
	// Group entries have a TTL of 600s, old items are pruned at not more than
	// 60 sec intervals
	groupRepo := group.NewInmemGroupRepository(600, 60)

	discoServer := server.NewDiscoServer(groupRepo)

	discoServer.Serve(":8080", ":9090", "office")
}
