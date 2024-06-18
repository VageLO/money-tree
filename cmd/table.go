package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	column_count = 0
)

type cell_type struct {
	row int
	column int
	text string
	selectable bool
	color tcell.Color
	reference interface{}
}

func InsertCell(s *cell_type) {
	align := tview.AlignLeft
	tableCell := tview.NewTableCell(s.text).
		SetTextColor(s.color).
		SetAlign(align).
		SetSelectable(s.selectable)
	if s.column >= 1 && s.column <= 3 {
		tableCell.SetExpansion(1)
	}
	if s.reference != nil {
		tableCell.SetReference(s.reference)
	}
	table.SetCell(s.row, s.column, tableCell)
}

func InsertRows(column_row []string, row int, data_row []string, id int) {

	for i, data := range data_row {
		InsertCell(&cell_type{
			row: row,
			column: i,
			text: data,
			selectable: true,
			color: tcell.ColorWhite,
			reference: struct {
					id    int
					field string
				}{id, column_row[i]},
		})
	}
}

func FillTable(request string) {
	table.Clear()
	
	table.SetTitle("Transactions")
	
	SelectTransactions(request)

	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		FillForm(column_count, row, false)

		pages.AddPage("Dialog", Dialog(form), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}

func AddToTable() {
	newRow := table.GetRowCount()

	FillForm(column_count, newRow, true)
	pages.AddPage("Dialog", Dialog(form), true, true)

	app.SetFocus(form)
}
