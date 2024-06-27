package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)



var (
	_ = InitDB()
	app   = tview.NewApplication()
	pages = tview.NewPages()
	modal = tview.NewModal()
	form = tview.NewForm()
	table = tview.NewTable().SetFixed(1, 1)
	file_table = tview.NewTable().SetBorders(false)
)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Init() {
	
	pages.AddPage("Transactions", Transactions(), true, true)

	// Shortcuts
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return Shortcuts(event)
	})

	// Start the application.
	err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run() 
	check(err)
}
