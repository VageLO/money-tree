package cmd

import (
	"errors"
	"main/action"
	m "main/modal"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	defer m.ErrorModal(source.Pages, source.Modal)

	table := source.Table
	tree := source.CategoryTree
	accounts := source.AccountList
	pages := source.Pages

	switch key := event.Key(); key {
	case tcell.KeyCtrlA, tcell.KeyInsert:
		check(ifFormExist(pages))

		if table.HasFocus() {
			action.AddToTable(source)
			return nil
		}
		if tree.HasFocus() {
			action.FormAddCategory(source)
			return nil
		}
		if accounts.HasFocus() {
			action.FormAddAccount(source)
			return nil
		}
	case tcell.KeyCtrlD, tcell.KeyDelete, tcell.KeyBackspace:
		check(ifFormExist(pages))

		if table.HasFocus() {
			action.DeleteTransaction(source)
			return nil
		}
		if tree.HasFocus() {
			action.RemoveCategory(source)
			return nil
		}
		if accounts.HasFocus() {
			action.RemoveAccount(source)
			return nil
		}
	case tcell.KeyCtrlR:
		check(ifFormExist(pages))

		if tree.HasFocus() {
			action.FormRenameCategory(source)
			return nil
		}
		if accounts.HasFocus() {
			action.FormRenameAccount(source)
			return nil
		}
	case tcell.KeyF2:
		check(ifFormExist(pages))

		pageName := "Imports"
		pages.AddPage(pageName, action.FileExporer(source, ".pdf", pageName), true, true)
	case tcell.KeyF3:
		check(ifFormExist(pages))
		DrawStats(source)
	}
	return event
}

func ifFormExist(pages *tview.Pages) error {
	if pages.HasPage("Form") {
		return errors.New("Close 'Form' window")
	}
	return nil
}
