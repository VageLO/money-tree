package cmd

import (
	"github.com/gdamore/tcell/v2"
	//"github.com/rivo/tview"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	switch key := event.Key(); key {
		case tcell.KeyCtrlA, tcell.KeyInsert:
			if table.HasFocus() {
				newRow := table.GetRowCount()
				table.InsertRow(newRow)
				
				form = FillForm(form, count, newRow, true)
				pages.AddPage("Dialog", Dialog(form), true, true)
				
				app.SetFocus(form)
				return nil
			}
			if tree.HasFocus() {
				AddNode()
			}
			if accounts.HasFocus() {
				AddAccount()
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
			}
			if accounts.HasFocus() {
				RemoveAccount()
			}
		case tcell.KeyCtrlR:
			if tree.HasFocus() {
				RenameNode()
				return nil
			}
			if accounts.HasFocus() {
				RenameAccount()
				return nil
			}
	}
	return event
}
