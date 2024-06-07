package cmd

import (
	"log"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	form = Table()
	table = tview.NewTable().
		SetFixed(1, 1)
	count = 0
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FillTable(columns []string, row int, data []string, id int) {
	
	for idx, title := range columns {
		color := tcell.ColorYellow
		align := tview.AlignLeft
		tableCell := tview.NewTableCell(title).
			SetTextColor(color).
			SetAlign(align).
			SetSelectable(false)
		if idx >= 1 && idx <= 3 {
			tableCell.SetExpansion(1)
		}
		table.SetCell(0, idx, tableCell)
	}
	
	for idx, cell_data := range data {
		color := tcell.ColorWhite
		align := tview.AlignLeft
		tableCell := tview.NewTableCell(cell_data).
			SetTextColor(color).
			SetAlign(align).
			SetSelectable(true)
		if idx >= 1 && idx <= 3 {
			tableCell.SetExpansion(1)
		}

		tableCell.SetReference(struct{id int; field string}{id, columns[idx]})
		table.SetCell(row, idx, tableCell)
	}
}

func Table() *tview.Form {

	SelectTransactions(`
		SELECT 
		Transactions.id, transaction_type, data, amount, Transactions.balance, Accounts.title as account, Categories.title as category
		FROM Transactions
		INNER JOIN Categories ON Categories.id = Transactions.category_id
		INNER JOIN Accounts ON Accounts.id = Transactions.account_id
	`)

	table.SetBorder(true).SetTitle("Transactions")

	form := Form()

	// Table action
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		form = FillForm(form, count, row, false)

		pages.AddPage("Dialog", Dialog(form), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')

	return form
}

func AddToTable() {
	newRow := table.GetRowCount()
	table.InsertRow(newRow)

	form = FillForm(form, count, newRow, true)
	pages.AddPage("Dialog", Dialog(form), true, true)

	app.SetFocus(form)
}

func AddTransaction() {

}
