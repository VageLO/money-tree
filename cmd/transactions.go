package cmd

import (
	"database/sql"
	"strconv"
	"errors"
	"time"
	"fmt"
	"os"
	
	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

var transaction_query string

type Transaction struct {
	id               int64
	account_id 		 int64
	to_account_id 	 sql.NullInt64
	category_id 	 int64
	transaction_type string
	date             string
	amount           float64
	to_amount		 sql.NullFloat64
	account          string
	to_account 		 sql.NullString
	category         string
	description		 string
}

func Transactions() tview.Primitive {
	
	query, err := os.ReadFile("./sql/Select_On_Transactions.sql")
	transaction_query = string(query)
    check(err)
	
	FillTable(transaction_query)
	
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
				pages.RemovePage("Form")
			}
			if modal.InRect(event.Position()) == false {
				pages.RemovePage("Modal")
			}
			if file_table.InRect(event.Position()) == false {
				pages.RemovePage("Files")
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
		err := rows.Scan(&t.id, &t.account_id, &t.to_account_id, &t.category_id, &t.transaction_type, &t.date, &t.amount, &t.to_amount, &t.description, &t.account, &t.to_account, &t.category)
		
		check(err)
		
		row := []string{t.description, t.date, t.account, t.category,
			strconv.FormatFloat(t.amount, 'f', 2, 32), t.transaction_type}
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
	pages.RemovePage("Form")

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
	pages.RemovePage("Form")

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
		return errors.New(fmt.Sprintf("%+v", t))
	}
	_, err := time.Parse("2006-01-02", t.date)
	if err != nil {
		return errors.New("Allowed date format (YYYY-MM-DD)")
	}
	return nil
}

func IsTransfer(cell *tview.TableCell, t *Transaction) bool {
	transfer := false
	tran_type := "debit"
	
	if cell != nil {
		reference := cell.GetReference().(Transaction)
		transfer = reference.to_account_id.Valid
		tran_type = reference.transaction_type
	}
	
	form.AddCheckbox("transfer", transfer, func(checked bool) {
		if !checked {
			TranTypes(tran_type, t)
			
			to_account_index := form.GetFormItemIndex("to_account")
			form.RemoveFormItem(to_account_index)
			
			to_amount_index := form.GetFormItemIndex("to_amount")
			form.RemoveFormItem(to_amount_index)
			return
		}
		
		if cell != nil {
			ToAccount(cell, t)
		} else {
			ToAccount(nil, t)
		}

		type_index := form.GetFormItemIndex("transaction_type")
		if type_index != -1 {
			form.RemoveFormItem(type_index)
		}
	})
	return transfer
}

func TranTypes(label string, t *Transaction) {
	types := []string{ "debit", "credit" }
	initial := 0

	for idx, title := range types {
		if title == label {
			initial = idx
		}
	}
	
	form.AddDropDown("transaction_type", types, initial, func(option string, optionIndex int) { 
		if types[optionIndex] != option {
			return
		}
		t.transaction_type = types[optionIndex]
	})
}

func ToAccount(cell *tview.TableCell, t *Transaction) {
	var label string
	var reference Transaction
	var amount string
	initial := 0
	
	if cell != nil {
		reference = cell.GetReference().(Transaction)
		label = reference.to_account.String
		amount = strconv.FormatFloat(reference.to_amount.Float64, 'f', 2, 32)
	}
	
	accounts, a_types := SelectAccounts()
	
	for idx, title := range accounts {
		if title == label {
			initial = idx
		}
	}
	T_Selected(accounts[initial], initial, a_types, t)
	
	form.AddDropDown("to_account", accounts, initial, func(option string, optionIndex int) { T_Selected(option, optionIndex, a_types, t) })
	
	form.AddInputField("to_amount", amount, 0, nil, func(text string) { 
		amount, _ := strconv.ParseFloat(text, 64)
		t.to_amount.Scan(amount)
	})
}