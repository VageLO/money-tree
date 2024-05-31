package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ( 
	form = Table()
	app = tview.NewApplication()
	pages = tview.NewPages()
)

func main() {

	pages.AddPage("Transactions", TransactionsTable(), true, true)
	
	// Shortcuts
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
		} else if event.Key() == tcell.KeyCtrlD {
			// Shortcut: deleting table row
			if table.HasFocus() == false {
				return nil
			}
			
			row, _ := table.GetSelection()
			table.RemoveRow(row)
			
			app.SetFocus(table)
			return nil
		}
		return event
	})

	// Start the application.
	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
