package cmd

import (
	"log"
	"reflect"

	"github.com/rivo/tview"
)

var (
	table = tview.NewTable().
		SetFixed(1, 1)
	count = 0
)

type Transaction struct {
	id               int
	account_id       int
	category_id      int
	account_title    string
	category_title   string
	transaction_type string
	data             string
	amount           float32
	balance          float32
	currency         string
}

func FillTable(data []string, inter interface{}) {
	table.Clear()
	t := reflect.TypeOf(inter)
	log.Println(t)

	//	for idx, title := range data {
	//		color := tcell.ColorDarkCyan
	//		align := tview.AlignLeft
	//		tableCell := tview.NewTableCell(title).
	//			SetTextColor(color).
	//			SetAlign(align).
	//			SetSelectable(true)
	//		table.SetCell(0, idx, tableCell)
	//	}
	//
	// db, err := sql.Open("sqlite3", "./database.db")
	// rows, err := db.Query(`Select * From Transactions`)
	// //rows, err := db.Query(`SELECT Transactions.*, Transactions.type as transaction_type, Categories.title as category_title, Accounts.title as account_title, Accounts.currency FROM Transactions
	// //INNER JOIN Accounts ON Transactions.account_id = Accounts.id
	// //INNER JOIN Categories ON Transactions.category_id = Categories.id;`)
	//
	//	if err != nil {
	//		log.Printf("select %v", err)
	//	}
	//
	// defer rows.Close()
	//
	//	for i := 0; rows.Next(); i++ {
	//		var transaction Transaction
	//
	//		if err := rows.Scan(&transaction.id,
	//			&transaction.account_id, &transaction.category_id,
	//			&transaction.account_title, &transaction.category_title,
	//			&transaction.transaction_type, &transaction.data,
	//			&transaction.amount, &transaction.balance,
	//			&transaction.currency); err != nil {
	//			log.Fatal(err)
	//		}
	//		log.Println(transaction)
	//		//	color := tcell.ColorDarkCyan
	//		//	align := tview.AlignLeft
	//		//	tableCell := tview.NewTableCell(title).
	//		//		SetTextColor(color).
	//		//		SetAlign(align).
	//		//		SetSelectable(true)
	//		//	table.SetCell(1, 0, tableCell)
	//	}
}

func Table() *tview.Form {

	tran := &Transaction{
		id:         3,
		account_id: 34,
	}
	column_title := []string{"data", "transaction_type"}
	FillTable(column_title, tran)

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
