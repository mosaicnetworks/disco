package commands

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/disco"
)

const (
	// Input box height.
	ih = 3
)

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

func chat() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}

	vdst, err := g.View("main")
	if err != nil {
		log.Fatalln(err)
	}

	dumper := hex.Dumper(vdst)

	go io.Copy(dumper, os.Stdout)

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
		v.Autoscroll = true
		v.Wrap = true
	}

	if v, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

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
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(v, 1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmdline", gocui.KeyEnter, gocui.ModNone, vcopy("main")); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func vcopy(dst string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		vdst, err := g.View(dst)
		if err != nil {
			return err
		}
		fmt.Fprint(vdst, v.Buffer())

		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		v.Clear()

		return nil
	}
}

func scrollView(v *gocui.View, dy int) error {
	if v != nil {
		v.Autoscroll = false
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+dy); err != nil {
			return err
		}
	}
	return nil
}
