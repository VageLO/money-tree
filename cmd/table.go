package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

var (
	table = tview.NewTable().
		SetFixed(1, 1)
	count = 0
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FillTable(request string) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	rows, err := db.Query(request)
	check(err)

	columns, err := rows.Columns()
	check(err)

	types, err := rows.ColumnTypes()

	for _, typpe := range types {
		for i := 0; rows.Next(); i++ {

			if err := rows.Scan(*typpe); err != nil {
				log.Fatal(err)
			}
			log.Println(typpe)
			//	color := tcell.ColorDarkCyan
			//	align := tview.AlignLeft
			//	tableCell := tview.NewTableCell(title).
			//		SetTextColor(color).
			//		SetAlign(align).
			//		SetSelectable(true)
			//	table.SetCell(1, 0, tableCell)
		}

	}

	for _, column := range columns {
		fmt.Printf("%v\n", column)
	}

	for idx, title := range columns {
		color := tcell.ColorDarkCyan
		align := tview.AlignLeft
		tableCell := tview.NewTableCell(title).
			SetTextColor(color).
			SetAlign(align).
			SetSelectable(true)
		table.SetCell(0, idx, tableCell)
	}

	defer rows.Close()

}

func Table() *tview.Form {

	FillTable(`Select * From Transactions`)

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
