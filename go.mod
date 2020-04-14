module github.com/mosaicnetworks/disco

go 1.14

// XXX reauires the webrtc-stream branch of Babble
// uncomment this line and point to your local version of the Babble repo
replace github.com/mosaicnetworks/babble => /home/jon/go/src/github.com/mosaicnetworks/babble

require (
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/jroimartin/gocui v0.4.0 // indirect
	github.com/manifoldco/promptui v0.7.0 // indirect
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/mosaicnetworks/babble v0.7.0
	github.com/nsf/termbox-go v0.0.0-20200204031403-4d2b513ad8be // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
)
