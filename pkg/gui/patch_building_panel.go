package gui

import (
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) refreshPatchBuildingPanel() error {
	if gui.GitCommand.PatchManager.IsEmpty() {
		return gui.handleEscapePatchBuildingPanel(gui.g, nil)
	}

	gui.State.SplitMainPanel = true

	// get diff from commit file that's currently selected
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	diff, err := gui.GitCommand.ShowCommitFile(commitFile.Sha, commitFile.Name, true)
	if err != nil {
		return err
	}

	secondaryDiff := gui.GitCommand.PatchManager.RenderPatchForFile(commitFile.Name, true, false, true)
	if err != nil {
		return err
	}

	empty, err := gui.refreshLineByLinePanel(diff, secondaryDiff, false)
	if err != nil {
		return err
	}

	if empty {
		return gui.handleEscapePatchBuildingPanel(gui.g, nil)
	}

	return nil
}

func (gui *Gui) handleAddSelectionToPatch(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.AddFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshPatchBuildingPanel(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleRemoveSelectionFromPatch(g *gocui.Gui, v *gocui.View) error {
	state := gui.State.Panels.LineByLine

	// add range of lines to those set for the file
	commitFile := gui.getSelectedCommitFile(gui.g)
	if commitFile == nil {
		return gui.renderString(gui.g, "commitFiles", gui.Tr.SLocalize("NoCommiteFiles"))
	}

	gui.GitCommand.PatchManager.RemoveFileLineRange(commitFile.Name, state.FirstLineIdx, state.LastLineIdx)

	if err := gui.refreshCommitFilesView(); err != nil {
		return err
	}

	if err := gui.refreshPatchBuildingPanel(); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) handleEscapePatchBuildingPanel(g *gocui.Gui, v *gocui.View) error {
	gui.State.Panels.LineByLine = nil
	gui.State.Contexts["main"] = "normal"

	return gui.switchFocus(gui.g, nil, gui.getCommitFilesView())
}

func (gui *Gui) refreshSecondaryPatchPanel() error {
	if !gui.GitCommand.PatchManager.IsEmpty() {
		gui.State.SplitMainPanel = true
		secondaryView := gui.getSecondaryView()
		secondaryView.Highlight = true
		secondaryView.Wrap = false

		gui.g.Update(func(*gocui.Gui) error {
			return gui.setViewContent(gui.g, gui.getSecondaryView(), gui.GitCommand.PatchManager.RenderAggregatedPatchColored(false))
		})
	} else {
		gui.State.SplitMainPanel = false
	}

	return nil
}
