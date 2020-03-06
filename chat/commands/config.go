package commands

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/mosaicnetworks/babble/src/config"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
)

// BChatConfig contains configuration for the Run command
type BChatConfig struct {
	BaseDir    string `mapstructure:"basedir"`
	DiscoAddr  string `mapstructure:"disco-addr"`
	SignalAddr string `mapstructure:"signal-addr"`
	LogLevel   string `mapstructure:"log"`
}

// NewDefaultBChatConfig creates a default BChatConfig
func NewDefaultBChatConfig() *BChatConfig {
	return &BChatConfig{
		BaseDir:    DefaultBaseDir(),
		DiscoAddr:  "http://localhost:8080", // XXX why do we need http:// ?
		SignalAddr: "localhost:8888",
	}
}

// BabbleConfig returns the Babble Configuration object associated with a BChat
// Config
func (c *BChatConfig) BabbleConfig(groupID string, moniker string) *config.Config {
	babbleConfig := config.NewDefaultConfig()
	babbleConfig.SetDataDir(path.Join(c.BaseDir, groupID))
	babbleConfig.LogLevel = c.LogLevel
	babbleConfig.Moniker = moniker
	babbleConfig.SignalAddr = c.SignalAddr
	babbleConfig.WebRTC = true
	babbleConfig.Proxy = dummy.NewInmemDummyClient(babbleConfig.Logger())
	return babbleConfig
}

// DefaultBaseDir return the default directory name for top-level BChat config
// based on the underlying OS, attempting to respect conventions.
func DefaultBaseDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".BChat")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "BChat")
		} else {
			return filepath.Join(home, ".bchat")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
