package action

import (
	"errors"
	s "main/structs"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	columns = []string{"description", "date", "account", "category", "amount", "transaction_type"}
	column_count = len(columns)
)

func InsertCell(c *s.Cell) {
	align := tview.AlignLeft
	tableCell := tview.NewTableCell(c.Text).
		SetTextColor(c.Color).
		SetAlign(align).
		SetSelectable(c.Selectable)
	tableCell.SetExpansion(1)
	if c.Reference != nil {
		tableCell.SetReference(c.Reference)
	}
	table.SetCell(c.Row, c.Column, tableCell)
}

func UpdateCell (c *s.Cell) {
	tableCell := table.GetCell(c.Row, c.Column)
	tableCell.SetText(c.Text)
	if c.Reference != nil {
		tableCell.SetReference(c.Reference)
	}
}

func InsertRows(column_row []string, row int, data_row []string, transaction s.Transaction) {

	for i, data := range data_row {
		InsertCell(&s.Cell{
			Row: row,
			Column: i,
			Text: data,
			Selectable: true,
			Color: tcell.ColorWhite,
			Reference: transaction,
		})
	}
}

func UpdateRows(column_row []string, row int, data_row []string, transaction s.Transaction) {

	for i, data := range data_row {
		UpdateCell(&s.Cell{
			Row: row,
			Column: i,
			Text: data,
			Reference: transaction,
		})
	}
}

func AddToTable() {
	defer ErrorModal()

	newRow := table.GetRowCount()
	if tree.GetRowCount() <= 0 || accounts.GetItemCount() <= 0 {
		check(errors.New("Account and category must be created"))
	}
	FillForm(column_count, newRow, true)
	pages.AddPage("Form", Modal(form, 30, 50), true, true)

	app.SetFocus(form)
}
