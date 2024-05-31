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
		return Shortcuts(event)
	})

	// Start the application.
	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
