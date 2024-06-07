package cmd

import (
	"database/sql"
	"log"
	"fmt"
	"strconv"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Transaction struct {
	id int
	transaction_type string
	data string
	amount float64
	balance float64
	account string
	category string
}

func TransactionsTable() tview.Primitive {
	
	form = Table()
	
	// List with accounts
	accounts := AccountsList()

	// Tree with categories
	categories := TreeView()

	//Flex
	flex := tview.NewFlex()

	top_flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(accounts, 0, 1, false).
		AddItem(categories, 0, 1, false)

	bottom_flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(table, 0, 2, true)

	modal_flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(top_flex, 20, 1, false).
		AddItem(bottom_flex, 0, 1, true)

	flex.AddItem(modal_flex, 0, 2, true)

	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event.Buttons() == tcell.Button1 {
			if form.InRect(event.Position()) == false {
				pages.RemovePage("Dialog")
			}
		}
		return event, action
	})

	return flex
}

func SelectTransactions(request string) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	rows, err := db.Query(request)
	check(err)
	
	columns, err := rows.Columns()
	check(err)
	columns = columns[1:]
	
	count = len(columns)
	
	for i := 1; rows.Next(); i++ {
		var t Transaction
		if err := rows.Scan(&t.id, &t.transaction_type, &t.data,
		&t.amount, &t.balance, &t.account, &t.category); err != nil {
			log.Fatal(err)
		}
		row := []string{t.transaction_type, t.data, 
		strconv.FormatFloat(t.amount, 'f', 2, 32), strconv.FormatFloat(t.balance, 'f', 2, 32), t.account, t.category}
		FillTable(columns, i, row, t.id)
	}

	defer rows.Close()
	db.Close()
}

func UpdateTransaction(cell *tview.TableCell, text string) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	t := cell.GetReference().(struct{id int; field string})
	str := fmt.Sprintf(`Update Transactions SET %v = ? WHERE id = ?`, t.field)
	
	_, err = db.Exec(str, text, t.id)
	check(err)
	
	db.Close()
}
