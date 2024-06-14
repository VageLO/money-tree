package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	form  = Table(`
		SELECT 
		Transactions.id, transaction_type, date, amount, Transactions.balance, Accounts.title as account, Categories.title as category
		FROM Transactions
		INNER JOIN Categories ON Categories.id = Transactions.category_id
		INNER JOIN Accounts ON Accounts.id = Transactions.account_id
	`)
	table = tview.NewTable().
		SetFixed(1, 1)
	column_count = 0
)

type row_settings struct {
	row int
	column int
	text string
	selectable bool
	color tcell.Color
	reference interface{}
}

func InsertRow(s *row_settings) {
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

func FillTable(columns []string, row int, data []string, id int) {

	for idx, title := range columns {
		
		InsertRow(&row_settings{
			row: 0,
			column: idx,
			text: title,
			selectable: false,
			color: tcell.ColorYellow,
		})
	}

	for idx, cell_data := range data {
		InsertRow(&row_settings{
			row: row,
			column: idx,
			text: cell_data,
			selectable: true,
			color: tcell.ColorWhite,
			reference: struct {
					id    int
					field string
				}{id, columns[idx]},
		})
	}
}

func Table(request string) *tview.Form {

	SelectTransactions(request)

	table.SetBorder(true).SetTitle("Transactions")

	form := Form()

	// Table action
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		form = FillForm(form, column_count, row, false)

		pages.AddPage("Dialog", Dialog(form), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')

	return form
}

func AddToTable() {
	newRow := table.GetRowCount()

	form = FillForm(form, column_count, newRow, true)
	pages.AddPage("Dialog", Dialog(form), true, true)

	app.SetFocus(form)
}
