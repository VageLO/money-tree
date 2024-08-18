package modal

import (
	"fmt"
	s "github.com/VageLO/money-tree/structs"
	"log"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

type Reference struct {
	Path string
}

func Modal(p tview.Primitive, hight, width int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, hight, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func ErrorModal(pages *tview.Pages, modal *tview.Modal) {
	if r := recover(); r != nil {
		//if pages.HasPage("Modal") {
		//	pages.RemovePage("Modal")
		//}
		err := fmt.Sprintf("Error: %v", r)
		modal.SetText(err)
		pages.AddPage("Modal", Modal(modal, 20, 40), true, true)
	}
}

func FileTable(source *s.Source, pageName string, files []string,
	selected func(path string, source *s.Source)) {
	defer ErrorModal(source.Pages, source.Modal)

	table := source.FileTable
	table.Clear()
	table.SetTitle(pageName).SetBorder(true)
	table.SetBorders(false).SetSelectable(true, false)

	count := 0
	for _, file := range files {
		fileName := filepath.Base(file)

		tableCell := tview.NewTableCell(fileName)
		tableCell.SetReference(Reference{file})
		tableCell.SetSelectable(true)

		table.SetCell(count, 0, tableCell)
		count++
	}

	table.SetSelectedFunc(func(row, column int) {
		cell := table.GetCell(row, column)

		reference := cell.GetReference().(Reference)

		selected(reference.Path, source)
	})

	source.Pages.AddPage(pageName, table, true, true)
}

func check(err error) {
	configPath, e := os.UserConfigDir()
	if e != nil {
		log.Fatalln(e)
	}

	file, e := os.OpenFile(filepath.Join(configPath, "money-tree", "tree.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
