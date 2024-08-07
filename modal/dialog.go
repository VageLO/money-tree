package modal

import (
	"fmt"
	s "main/structs"
	"os/exec"
	"path/filepath"

	"github.com/rivo/tview"
)

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
	table.SetTitle(pageName).SetBorder(true)
	table.SetBorders(false).SetSelectable(true, false)

	count := 0
	for _, file := range files {
		fileName := filepath.Base(file)

		tableCell := tview.NewTableCell(fileName)
		tableCell.SetReference(struct {
			path string
		}{file})
		tableCell.SetSelectable(true)

		table.SetCell(count, 0, tableCell)
		count++
	}

	table.SetSelectedFunc(func(row, column int) {
		cell := table.GetCell(row, column)

		reference := cell.GetReference().(struct {
			path string
		})

		selected(reference.path, source)
	})

	button := tview.NewButton("Cancel")
	button.SetSelectedFunc(func() {
		source.Pages.RemovePage(pageName)
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 15, true).
		AddItem(button, 0, 1, false)

		//	x, _, _, _ := flex.GetRect()
	source.Pages.AddPage(pageName, flex, true, true)
}

func OpenFiles(filePath string, source *s.Source) {
	defer ErrorModal(source.Pages, source.Modal)

	err := exec.Command("xdg-open", filePath).Start()
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
