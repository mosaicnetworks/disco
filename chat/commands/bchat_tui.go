package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
)

const (
	inputViewName  = "input"
	outputViewName = "output"
)

/*******************************************************************************
InputWidget
*******************************************************************************/

// InputWidget is a text input widget that implements the gocui.Manager
// interface. It wires the ENTER key to submit a msg to the babble proxy.
type InputWidget struct {
	Proxy   *dummy.InmemDummyClient
	Moniker string
}

// NewInputWidget instantiates a new InputWidget.
func NewInputWidget(proxy *dummy.InmemDummyClient, moniker string) *InputWidget {
	return &InputWidget{
		Proxy:   proxy,
		Moniker: moniker,
	}
}

// Layout implements the gocui.Manager interface
func (w *InputWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(inputViewName, -1, maxY-5, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true
		v.Frame = true

		// Wire up the Enter key to submit a transaction to Babble through the
		// proxy
		if err := g.SetKeybinding(inputViewName,
			gocui.KeyEnter,
			gocui.ModNone,
			func(g *gocui.Gui, v *gocui.View) error {
				msg := v.ViewBuffer()
				if len(msg) == 0 {
					return nil
				}
				prefixedMsg := fmt.Sprintf("%s: %s", w.Moniker, msg)
				w.Proxy.SubmitTx([]byte(prefixedMsg))
				v.SetCursor(0, 0)
				v.SetOrigin(0, 0)
				v.Clear()
				return nil
			},
		); err != nil {
			return err
		}

		if _, err := g.SetCurrentView(inputViewName); err != nil {
			return err
		}
	}

	return nil
}

/*******************************************************************************
OutputWidget
*******************************************************************************/

// OutputWidget is a text display widget that implements the gocui.Manager
// interface. It is usd to diplay text.
type OutputWidget struct {
}

// NewOutputWidget instantiates a new MainWidget
func NewOutputWidget() *OutputWidget {
	return &OutputWidget{}
}

// Layout implments the gocui.Manager interface
func (w *OutputWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(outputViewName, -1, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = true
	}

	return nil
}

/*******************************************************************************
BChatGui
*******************************************************************************/

// BChatGui represents a terminal-base graphical user interface for BChat
type BChatGui struct {
	gui   *gocui.Gui
	proxy *dummy.InmemDummyClient
}

// NewBChatGui instantiates a new BChatGui
func NewBChatGui(proxy *dummy.InmemDummyClient, moniker string) BChatGui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	inputWidget := NewInputWidget(proxy, moniker)
	outputWidget := NewOutputWidget()
	g.SetManager(inputWidget, outputWidget)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return BChatGui{
		gui:   g,
		proxy: proxy,
	}
}

// Loop runs the main loop where and regularly updates the output widget to
// reflect the latest committed messages
func (b *BChatGui) Loop() error {
	defer b.gui.Close()

	go func() {
		for {
			b.gui.Update(func(g *gocui.Gui) error {
				v, err := g.View(outputViewName)
				if err != nil {
					return err
				}
				v.Clear()
				for _, tx := range b.proxy.GetCommittedTransactions() {
					fmt.Fprintf(v, "%s", tx)
				}
				return nil
			})
			time.Sleep(100 * time.Millisecond)
		}
	}()

	if err := b.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
