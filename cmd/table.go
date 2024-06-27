package cmd

import (
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	columns = []string{"date", "transaction_type", "account", "category", "amount", "description"}
	column_count = len(columns)
)

type cell_type struct {
	row int
	column int
	text string
	selectable bool
	color tcell.Color
	reference interface{}
}

func InsertCell(c *cell_type) {
	align := tview.AlignLeft
	tableCell := tview.NewTableCell(c.text).
		SetTextColor(c.color).
		SetAlign(align).
		SetSelectable(c.selectable)
	if c.column >= 1 && c.column <= 3 {
		tableCell.SetExpansion(1)
	}
	if c.reference != nil {
		tableCell.SetReference(c.reference)
	}
	table.SetCell(c.row, c.column, tableCell)
}

func UpdateCell (c *cell_type) {
	tableCell := table.GetCell(c.row, c.column)
	tableCell.SetText(c.text)
	if c.reference != nil {
		tableCell.SetReference(c.reference)
	}
}

func InsertRows(column_row []string, row int, data_row []string, transaction Transaction) {

	for i, data := range data_row {
		InsertCell(&cell_type{
			row: row,
			column: i,
			text: data,
			selectable: true,
			color: tcell.ColorWhite,
			reference: transaction,
		})
	}
}

func UpdateRows(column_row []string, row int, data_row []string, transaction Transaction) {

	for i, data := range data_row {
		UpdateCell(&cell_type{
			row: row,
			column: i,
			text: data,
			reference: transaction,
		})
	}
}

func FillTable(request string) {
	table.Clear()
	
	table.SetTitle("Transactions")
	
	SelectTransactions(request)
	
	table.Select(1, 1).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		FillForm(column_count, row, false)

		pages.AddPage("Form", Modal(form, 20, 50), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}

func AddToTable() {
	defer ErrorModal()

	newRow := table.GetRowCount()
	if tree.GetRowCount() <= 0 || accounts.GetItemCount() <= 0 {
		check(errors.New("Account and category must be created"))
	}
	FillForm(column_count, newRow, true)
	pages.AddPage("Form", Modal(form, 20, 50), true, true)

	app.SetFocus(form)
}
