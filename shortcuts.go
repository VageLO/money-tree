package main

import (
	"github.com/gdamore/tcell/v2"
	//"github.com/rivo/tview"
)

func Shortcuts(event *tcell.EventKey) *tcell.EventKey {
	// Shortcut: creating new row in table and open dialog window with table attributes
	if event.Key() == tcell.KeyCtrlA {
		if table.HasFocus() == false {
			return nil
		}
		
		newRow := table.GetRowCount()
		table.InsertRow(newRow)
		
		form = FillForm(form, count, newRow, true)
		pages.AddPage("Dialog", Dialog(form), true, true)
		
		app.SetFocus(form)
		return nil
	}
	// Shortcut: deleting table row
	if event.Key() == tcell.KeyCtrlD {
		if table.HasFocus() == false {
			return nil
		}
		
		row, _ := table.GetSelection()
		table.RemoveRow(row)
		
		app.SetFocus(table)
		return nil
	} 
	if event.Key() == tcell.KeyCtrlR {
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
