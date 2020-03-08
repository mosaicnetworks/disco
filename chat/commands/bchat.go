package commands

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/mosaicnetworks/babble/src/babble"
	"github.com/mosaicnetworks/babble/src/proxy"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type BChat struct {
	Proxy  proxy.AppProxy
	Engine *babble.Babble
}

// newBChat returns a BChat node which is composed of a Babble node and an
// inmemory DummyApp proxy
func newBChat(config *BChatConfig, groupID string, moniker string) (*BChat, error) {
	// sets babble configuration and app proxy
	babbleConfig := config.BabbleConfig(groupID, moniker)

	// Set a special kind of logger for the proxy
	chatLogger := newLogger(babbleConfig.DataDir).WithField("component", "chat")
	proxy := dummy.NewInmemDummyClient(chatLogger)

	babbleConfig.Proxy = proxy

	engine := babble.NewBabble(babbleConfig)

	if err := engine.Init(); err != nil {
		babbleConfig.Logger().Error("Cannot initialize engine:", err)
		return nil, err
	}

	return &BChat{
		Proxy:  proxy,
		Engine: engine,
	}, nil
}

func newLogger(dir string) *logrus.Logger {
	logger := logrus.New()

	pathMap := lfshook.PathMap{}

	infoPath := path.Join(dir, "info.log")

	_, err := os.OpenFile(infoPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Info("Failed to open info.log file, using default stderr")
	} else {
		pathMap[logrus.InfoLevel] = infoPath
	}

	debugPath := path.Join(dir, "debug.log")

	_, err = os.OpenFile(debugPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Info("Failed to open debug.log file, using default stderr")
	} else {
		pathMap[logrus.DebugLevel] = debugPath
	}

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))

	return logger
}

// waitForInterrupt listens for an interrupt signal and politely leaves the
// BChat (Babble leave request) before returning
func waitForInterrupt(engine *babble.Babble) {
	//Prepare sigCh to relay SIGINT and SIGTERM system calls
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	engine.Node.Leave()
}
