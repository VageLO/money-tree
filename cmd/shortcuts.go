package cmd

import (
	"github.com/gdamore/tcell/v2"
	"main/action"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	table := source.Table
	tree := source.CategoryTree
	accounts := source.AccountList

	switch key := event.Key(); key {
	case tcell.KeyCtrlA, tcell.KeyInsert:
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
		if tree.HasFocus() {
			action.FormRenameCategory(source)
			return nil
		}
		if accounts.HasFocus() {
			action.FormRenameAccount(source)
			return nil
		}
	case tcell.KeyF2:
		FilePicker("./sql")
	}
	return event
}
