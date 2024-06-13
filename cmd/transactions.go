package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

type Transaction struct {
	id               int
	transaction_type string
	date             string
	amount           float64
	balance          float64
	account          string
	category         string
}

func TransactionsTable() tview.Primitive {

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

	columns := []string{"transaction_type", "date", "amount", "balance", "account", "category"}

	column_count = len(columns)

	for i := 1; rows.Next(); i++ {
		var t Transaction
		if err := rows.Scan(&t.id, &t.transaction_type, &t.date,
			&t.amount, &t.balance, &t.account, &t.category); err != nil {
			log.Fatal(err)
		}
		row := []string{t.transaction_type, t.date,
			strconv.FormatFloat(t.amount, 'f', 2, 32), strconv.FormatFloat(t.balance, 'f', 2, 32), t.account, t.category}
		FillTable(columns, i, row, t.id)
	}

	defer rows.Close()
	db.Close()
}

func UpdateTransaction(cell *tview.TableCell, text string) {
	if text == "" {
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	t := cell.GetReference().(struct {
		id    int
		field string
	})

	str := fmt.Sprintf(`Update Transactions SET %v = ? WHERE id = ?`, t.field)

	_, err = db.Exec(str, text, t.id)
	check(err)

	db.Close()
}

func AddTransaction(t *add_transaction) {

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	query := `
	INSERT INTO Transactions (account_id, category_id, transaction_type,
	date, amount, balance) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, t.account, t.category, t.transaction_type, t.date, t.amount, t.balance)
	check(err)

	db.Close()
}
