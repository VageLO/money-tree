package cmd

import (
	s "github.com/VageLO/money-tree/structs"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
    "os"
    "path/filepath"
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
    if err == nil {
        return
    }

    dir, e := os.UserConfigDir()
	if e != nil {
        log.Fatalln(e)
    }
	configPath := filepath.Join(dir, "money-tree")

    // Create money-tree directory in UserConfigDir
	if e = os.Mkdir(configPath, 0750); e != nil && !os.IsExist(e) {
		log.Fatalln(e)
	}
    
    // Create log file
    logFile, e := os.OpenFile(filepath.Join(configPath, "tree.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening log file: %v\n", e)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println(err)
    panic(err)
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
