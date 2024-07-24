package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	s "main/structs"
)

var (
	_      = InitDB()
	source = &s.Source{
		App:          tview.NewApplication(),
		AccountList:  tview.NewList(),
		CategoryTree: tview.NewTreeView(),
		Form:         tview.NewForm(),
		Table:        tview.NewTable().SetFixed(1, 1),
		FileTable:    tview.NewTable(),
		Modal:        tview.NewModal(),
		Pages:        tview.NewPages(),
		Columns:      []string{"Description", "Date", "Account", "Category", "Amount", "Transaction Type"},
	}
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Init() {

	source.Form.SetCancelFunc(func() {
		source.Pages.RemovePage("Form")
	})

	source.Pages.AddPage("Transactions", Transactions(), true, true)

	source.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return Shortcuts(event)
	})

	err := source.App.SetRoot(source.Pages, true).EnableMouse(true).EnablePaste(true).Run()
	check(err)

}
