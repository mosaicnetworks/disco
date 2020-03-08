package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/disco"
)

const (
	// Input box height.
	ih = 3
)

// XXX
var _logFile string
var _submitCh chan []byte

func promptTitle() string {
	titlePrompt := promptui.Prompt{
		Label: "Title",
	}

	title, err := titlePrompt.Run()
	if err != nil {
		log.Panicln(err)
	}

	return title
}

func promptMoniker() string {
	monikerPrompt := promptui.Prompt{
		Label: "Moniker",
	}

	moniker, err := monikerPrompt.Run()
	if err != nil {
		log.Panicln(err)
	}

	return moniker
}

func selectGroup(allGroups map[string]*disco.Group) *disco.Group {
	items := []*disco.Group{}
	for _, g := range allGroups {
		items = append(items, g)
	}

	templates := &promptui.SelectTemplates{
		Active:   "{{ .ID | cyan }} {{ .Title | red }}",
		Inactive: "{{ .ID }} {{ .Title }}",
		Selected: "{{ .ID }} {{ .Title }}",
	}

	selector := promptui.Select{
		Label:     "Select Group",
		Templates: templates,
		Items:     items,
	}

	selectedGroupIndex, _, err := selector.Run()
	if err != nil {
		log.Panicln(err)
	}

	selectedGroup := items[selectedGroupIndex]

	return selectedGroup
}

func chat(submitCh chan []byte, file string) {
	//XXX
	_logFile = file
	_submitCh = submitCh

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("main")
				if err != nil {
					return err
				}
				v.Clear()
				b, err := ioutil.ReadFile(_logFile)
				if err != nil {
					panic(err)
				}
				fmt.Fprintf(v, "%s", b)
				return nil
			})
			time.Sleep(100 * time.Millisecond)
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", -1, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Frame = true

		b, err := ioutil.ReadFile(_logFile)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
	}

	if v, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		g.Cursor = true
		v.Frame = true

		if _, err := g.SetCurrentView("cmdline"); err != nil {
			return err
		}
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmdline",
		gocui.KeyEnter,
		gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_submitCh <- []byte(v.ViewBuffer())
			v.SetCursor(0, 0)
			v.SetOrigin(0, 0)
			v.Clear()
			return nil
		},
	); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
