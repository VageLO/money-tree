package parser

import (
	"os"
	"strings"
	
	"github.com/rivo/tview"
)

func FilePicker(path string) (*tview.Table, error) {
	table := tview.NewTable().SetBorders(true)

	folder, err := os.Open(path)
	if err != nil {
		return table, err
	}
	
    files, err := folder.Readdir(0)
    if err != nil {
		return table, err
	}
	
	count := 0
    for _, v := range files {
		slice := strings.Split(v.Name(), ".")
		if v.IsDir() || slice[len(slice)-1] != "pdf" {
			continue
		}
        tableCell := tview.NewTableCell(v.Name())
		table.SetCell(count, 0, tableCell)
		count++
    }
	return table, nil
}