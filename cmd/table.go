package cmd

import (
	"database/sql"
	"log"
	"strconv"
	
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
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

func FillTable(request string) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	rows, err := db.Query(request)
	check(err)
	
	var m = map[string]interface{}{
		"id": 0,
	}
	
	columns, err := rows.Columns()
	check(err)
	
	for idx, title := range columns {
		color := tcell.ColorDarkCyan
		align := tview.AlignLeft
		tableCell := tview.NewTableCell(title).
			SetTextColor(color).
			SetAlign(align).
			SetSelectable(true)
		table.SetCell(0, idx, tableCell)
	}
	
	//types, err := rows.ColumnTypes()

	for i := 1; rows.Next(); i++ {
		
		
		if err := rows.Scan(&m["id"]); err != nil {
			log.Fatal(err)
		}
		
		color := tcell.ColorDarkCyan
		align := tview.AlignLeft
		tableCell := tview.NewTableCell(strconv.Itoa(m["id"])).
			SetTextColor(color).
			SetAlign(align).
			SetSelectable(true)
		table.SetCell(i, 0, tableCell)
	}

	defer rows.Close()
}

func Table() *tview.Form {

	FillTable(`Select id From Transactions`)

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
