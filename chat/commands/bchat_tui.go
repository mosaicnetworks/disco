package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
)

type InputWidget struct {
	Proxy   *dummy.InmemDummyClient
	Moniker string
}

func NewInputWidget(proxy *dummy.InmemDummyClient, moniker string) *InputWidget {
	return &InputWidget{
		Proxy:   proxy,
		Moniker: moniker,
	}
}

func (w *InputWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("input", -1, maxY-5, maxX, maxY)
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
		if err := g.SetKeybinding("input",
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

		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

type MainWidget struct {
}

func NewMainWidget() *MainWidget {
	return &MainWidget{}
}

func (w *MainWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = true
	}

	return nil
}

type BChatTui struct {
	cui   *gocui.Gui
	proxy *dummy.InmemDummyClient
}

func NewBChatTui(proxy *dummy.InmemDummyClient, moniker string) BChatTui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	inputWidget := NewInputWidget(proxy, moniker)
	mainWidget := NewMainWidget()
	g.SetManager(inputWidget, mainWidget)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return BChatTui{
		cui:   g,
		proxy: proxy,
	}
}

func (b *BChatTui) Loop() error {
	defer b.cui.Close()

	go func() {
		for {
			b.cui.Update(func(g *gocui.Gui) error {
				v, err := g.View("main")
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

	if err := b.cui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
