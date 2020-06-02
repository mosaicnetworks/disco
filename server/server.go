package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mosaicnetworks/babble/src/net/signal/wamp"
	"github.com/mosaicnetworks/disco/group"
	"github.com/pion/logging"
	"github.com/pion/turn/v2"
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

// Serve starts the peer-discovery, signaling, and TURN servers.
func (s *DiscoServer) Serve(
	discoAddr string,
	signalAddr string,
	turnAddr string,
	turnUsername string,
	turnPassword string,
	realm string,
	ttl time.Duration,
	ttlHearbeat time.Duration) {

	// Create and start WAMP server
	wampServer, err := wamp.NewServer(
		signalAddr,
		realm,
		s.certFile,
		s.keyFile,
		s.logger)
	if err != nil {
		log.Fatal(err)
	}
	go wampServer.Run()
	defer wampServer.Shutdown()

	// Create and start TURN server
	turnServer, err := createAndStartTURNServer(
		turnAddr,
		turnUsername,
		turnPassword,
		realm,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer turnServer.Close()

	// Start the TTL routine that deletes groups when the exceed their Time To
	// Live
	go s.processTTL(ttlHearbeat, ttl)

	// Configure and start discovery API
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/group", s.createGroup).Methods("POST")
	router.HandleFunc("/groups", s.getGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", s.getGroup).Methods("GET")
	router.HandleFunc("/groups/{id}", s.updateGroup).Methods("PATCH")
	router.HandleFunc("/groups/{id}", s.deleteGroup).Methods("DELETE")
	log.Fatal(http.ListenAndServeTLS(discoAddr, s.certFile, s.keyFile, router))
	return
}

func createAndStartTURNServer(
	turnAddr string,
	turnUsername string,
	turnPassword string,
	realm string) (*turn.Server, error) {

	// Populate the map of authorised users with the single user defined by
	// turnUsername and turnPassword.
	usersMap := map[string][]byte{}
	usersMap[turnUsername] = turn.GenerateAuthKey(turnUsername, realm, turnPassword)

	// Split the turnAddr into IP and Port.
	split := strings.Split(turnAddr, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("Invalid ICE address format")
	}
	bindAddr := split[0]
	icePort := split[1]

	// Create a UDP listener to pass into pion/turn
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+icePort)
	if err != nil {
		return nil, fmt.Errorf("Failed to create TURN server listener: %s", err)
	}

	// Override the default log level
	logFactory := logging.NewDefaultLoggerFactory()
	logFactory.DefaultLogLevel = logging.LogLevelInfo

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: realm,
		// Set AuthHandler callback
		// This is called everytime a user tries to authenticate with the TURN
		// server. Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			if key, ok := usersMap[username]; ok {
				return key, true
			}
			return nil, false
		},
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP(bindAddr), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",             // But actually be listening on every interface
				},
			},
		},
		LoggerFactory: logFactory,
	})
	if err != nil {
		return nil, fmt.Errorf("Fail to create TURN server: %s", err)
	}

	return s, nil
}

// processTTL deletes groups that have exceeded their Time To Live (TTL). It
// will check each group at event intervals defined by the heartbeat parameter.
func (s *DiscoServer) processTTL(heartbeat time.Duration, ttl time.Duration) {
	for now := range time.Tick(heartbeat) {
		allGroups, _ := s.repo.GetAllGroups()
		for gid, g := range allGroups {
			if g.LastUpdated+int64(ttl.Seconds()) < now.Unix() {
				s.repo.DeleteGroup(gid)
				s.logger.Debugf("Deletet group %s, TTL exceeded", g.Name)
			}
		}
	}
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
