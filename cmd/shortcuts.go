package cmd

import (
	"errors"
	"fmt"
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
	case tcell.KeyEscape:
		if name, _ := source.Pages.GetFrontPage(); name != "" && name != "Transactions" {
			source.Pages.RemovePage(name)
			return nil
		}
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
		if tree.GetRowCount() <= 1 || accounts.GetItemCount() <= 1 {
			check(errors.New("Account and category must be created"))
		}
		action.FileExporer(source, ".pdf", "Imports")
		return nil
	case tcell.KeyF3:
		check(ifFormExist(pages))
		if tree.GetRowCount() <= 1 || accounts.GetItemCount() <= 1 {
			check(errors.New("Account and category must be created"))
		}
		Statistics(source)
		return nil
	}
	return event
}

func ifFormExist(pages *tview.Pages) error {
	if pages.HasPage("Form") {
		return errors.New(fmt.Sprintf("Close '%v' by pressing Esc", source.Form.GetTitle()))
	}
	return nil
}
