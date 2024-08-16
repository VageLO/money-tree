package cmd

import (
    "os"
    "log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	s "github.com/VageLO/money-tree/structs"
)

var (
    _      = readConfig()   
	_      = initDB()
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
        Config:       s.Config{},
	}
)

func check(err error) {
    file, e := os.OpenFile("./tree.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if e != nil {
        log.Fatalf("error opening log file: %v", e)
    }
    defer file.Close()
    log.SetOutput(file)

	if err != nil {
        log.Println(err)
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
