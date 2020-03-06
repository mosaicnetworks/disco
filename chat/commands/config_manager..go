package commands

import (
	"crypto/ecdsa"
	"path"

	"github.com/google/uuid"
	"github.com/mosaicnetworks/babble/src/crypto/keys"
	"github.com/mosaicnetworks/babble/src/peers"
	"github.com/mosaicnetworks/disco"
)

// ConfigManager manages configuration directories for bchat groups. These
// config directories are used as datadirs for the underlying Babble nodes.
type ConfigManager struct {
	dir string
}

// NewConfigManager instantiates a ConfigManager with a base directory which
// will contain all the group configurations
func NewConfigManager(dir string) *ConfigManager {
	return &ConfigManager{
		dir: dir,
	}
}

// CreateNew generates a new random group-id and creates a configuration
// directory for it. In this directory, it will dump a new priv_key and
// peers.json file with a single peer (corresponding to the priv_key).
func (cm *ConfigManager) CreateNew(moniker string) (string, error) {
	id := uuid.New().String()

	privKey, err := cm.dumpNewPrivKey(id)
	if err != nil {
		return "", err
	}

	peers := []*peers.Peer{
		peers.NewPeer(
			keys.PublicKeyHex(&privKey.PublicKey),
			keys.PublicKeyHex(&privKey.PublicKey),
			moniker,
		),
	}

	err = cm.dumpPeers(id, peers)
	if err != nil {
		return "", err
	}

	return id, nil
}

// CreateForJoin is used when joining a group. Given a Group object, it creates
// a config directory for that group (based on ID), and writes a new priv_key
// file. The group's peers is dumped to a peers.json file in the config
// directory.
func (cm *ConfigManager) CreateForJoin(group *disco.Group) error {
	_, err := cm.dumpNewPrivKey(group.ID)
	if err != nil {
		return err
	}

	return cm.dumpPeers(group.ID, group.Peers)
}

func (cm *ConfigManager) dumpNewPrivKey(groupID string) (*ecdsa.PrivateKey, error) {
	// generate new private key
	privKey, _ := keys.GenerateECDSAKey()

	// write it to dir/groupID
	keyFile := keys.NewSimpleKeyfile(path.Join(cm.dir, groupID, "priv_key"))
	err := keyFile.WriteKey(privKey)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func (cm *ConfigManager) dumpPeers(groupID string, peerList []*peers.Peer) error {
	jsonPeerSet := peers.NewJSONPeerSet(path.Join(cm.dir, groupID), true)
	return jsonPeerSet.Write(peerList)
}
