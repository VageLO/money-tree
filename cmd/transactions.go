package cmd

import (
	"database/sql"
	"strconv"
	"errors"
	"time"
	
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

type Transaction struct {
	id               int64
	account_id int64
	category_id int64
	transaction_type string
	date             string
	amount           float64
	account          string
	category         string
	description		 string
}

func Transactions() tview.Primitive {

	FillTable(`SELECT Transactions.*, Accounts.title, Categories.title FROM Transactions INNER JOIN Categories ON Categories.id = Transactions.category_id
INNER JOIN Accounts ON Accounts.id = Transactions.account_id`)
	
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
			if modal.InRect(event.Position()) == false {
				pages.RemovePage("Modal")
			}
		}
		return event, action
	})

	return flex
}

func SelectTransactions(request string) {
	defer ErrorModal()
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	rows, err := db.Query(request)
	check(err)
	
	for i, column_title := range columns {
		InsertCell(&cell_type{
			row: 0,
			column: i,
			text: column_title,
			selectable: false,
			color: tcell.ColorYellow,
		})
	}
	
	for i := 1; rows.Next(); i++ {
		var t Transaction
		err := rows.Scan(&t.id, &t.account_id, &t.category_id, &t.transaction_type, &t.date, &t.amount, &t.description, &t.account, &t.category)
		
		check(err)
		
		row := []string{t.date, t.transaction_type, t.account, t.category,
			strconv.FormatFloat(t.amount, 'f', 2, 32), t.description}
		InsertRows(columns, i, row, t)
	}

	defer rows.Close()
	defer db.Close()
}

func UpdateTransaction(t Transaction, row int) {
	defer ErrorModal()
	check(t.isEmpty())
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(Transaction)

	query := `Update Transactions SET account_id = ?, category_id = ?, 
	transaction_type = ?, date = ?, amount = ?, description = ? WHERE id = ?`
	_, err = db.Exec(query, t.account_id, t.category_id, t.transaction_type, t.date, strconv.FormatFloat(t.amount, 'f', 2, 32), t.description, transaction.id)
	check(err)
	
	t.id = transaction.id
	
	data := []string{t.date, t.transaction_type, t.account, t.category,
	strconv.FormatFloat(t.amount, 'f', 2, 32), t.description}
	
	UpdateRows(columns, row, data, t)
	AccountsList()
	pages.RemovePage("Dialog")

	defer db.Close()
}

func AddTransaction(t Transaction, newRow int) {
	defer ErrorModal()
	check(t.isEmpty())
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	query := `INSERT INTO Transactions (account_id, category_id, transaction_type,
	date, amount, description) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, t.account_id, t.category_id, t.transaction_type, t.date, strconv.FormatFloat(t.amount, 'f', 2, 32), t.description)
	check(err)
	
	created_id, err := result.LastInsertId()
	check(err)

	t.id = created_id
	
	row := []string{t.date, t.transaction_type, t.account, t.category,
	strconv.FormatFloat(t.amount, 'f', 2, 32), t.description}
	
	InsertRows(columns, newRow, row, t)
	AccountsList()
	pages.RemovePage("Dialog")

	defer db.Close()
}

func DeleteTransaction() {
	defer ErrorModal()
	row, _ := table.GetSelection()

	if table.GetRowCount() <= 1 {return}
	
	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(Transaction)
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	query := `DELETE FROM Transactions WHERE id = ?`

	_, err = db.Exec(query, transaction.id)
	check(err)
	
	defer db.Close()
	AccountsList()
	table.RemoveRow(row)
}

func (t Transaction) isEmpty() error {
	if t.account_id == 0 || t.category_id == 0 || t.transaction_type == "" || t.date == "" || t.amount == 0 {
		return errors.New("Empty field or can't be zero")
	}
	_, err := time.Parse("2006-01-02", t.date)
	if err != nil {
		return errors.New("Allowed date format (YYYY-MM-DD)")
	}
	return nil
}