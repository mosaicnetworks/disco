package commands

import (
	"github.com/manifoldco/promptui"
	"github.com/mosaicnetworks/disco"
)

func promptTitle() string {
	titlePrompt := promptui.Prompt{
		Label: "Title",
	}

	title, err := titlePrompt.Run()
	if err != nil {
		panic(1)
	}

	return title
}

func promptMoniker() string {
	monikerPrompt := promptui.Prompt{
		Label: "Moniker",
	}

	moniker, err := monikerPrompt.Run()
	if err != nil {
		panic(1)
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
		panic(0)
	}

	selectedGroup := items[selectedGroupIndex]

	return selectedGroup
}
