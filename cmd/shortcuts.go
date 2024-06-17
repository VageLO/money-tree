package cmd

import (
	"github.com/gdamore/tcell/v2"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
	case tcell.KeyCtrlA, tcell.KeyInsert:
		if table.HasFocus() {
			AddToTable()
			return nil
		}
		if tree.HasFocus() {
			FormAddNode()
			return nil
		}
		if accounts.HasFocus() {
			FormAddAccount()
			return nil
		}
	case tcell.KeyCtrlD, tcell.KeyDelete, tcell.KeyBackspace:
		if table.HasFocus() {
			row, _ := table.GetSelection()
			table.RemoveRow(row)

			app.SetFocus(table)
			return nil
		}
		if tree.HasFocus() {
			RemoveNode()
			return nil
		}
		if accounts.HasFocus() {
			RemoveAccount()
			return nil
		}
	case tcell.KeyCtrlR:
		if tree.HasFocus() {
			FormRenameNode()
			return nil
		}
		if accounts.HasFocus() {
			FormRenameAccount()
			return nil
		}
	}
	return event
}
