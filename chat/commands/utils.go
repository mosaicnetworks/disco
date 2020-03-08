package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/babble/src/proxy/dummy"
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

func chat(proxy *dummy.InmemDummyClient, file string) {
	//XXX
	_submitCh = proxy.SubmitCh()

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
				for _, tx := range proxy.GetCommittedTransactions() {
					fmt.Fprintf(v, "%s", tx)
				}
				return nil
			})
			time.Sleep(100 * time.Millisecond)
		}
	}()

}

func layout(g *gocui.Gui) error {
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

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}
