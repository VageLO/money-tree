package cmd

import (
	"main/parser"

	"github.com/gdamore/tcell/v2"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	defer ErrorModal()

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
			DeleteTransaction()
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
	case tcell.KeyF2:
		// TODO: deal with exit focus
		file_table, err := parser.FilePicker("./sql")
		check(err)
		x, y, _, _ := file_table.GetRect()
		pages.AddPage("Modal", Modal(file_table, x, y), true, true)
	}
	return event
}
