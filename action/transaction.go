package action

import (
	"database/sql"
	"errors"
	"fmt"
	s "main/structs"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func SelectTransactions(request string) error {
	if db, err := sql.Open("sqlite3", "./database.db"); err != nil {
		return err
	}

	if rows, err := db.Query(request); err != nil {
		return err
	}

	for i, column_title := range columns {
		InsertCell(&s.Cell{
			row:        0,
			column:     i,
			text:       column_title,
			selectable: false,
			color:      tcell.ColorYellow,
		})
	}

	for i := 1; rows.Next(); i++ {
		var t s.Transaction

		err := rows.Scan(&t.id, &t.account_id, &t.to_account_id, &t.category_id, &t.transaction_type, &t.date, &t.amount, &t.to_amount, &t.description, &t.account, &t.to_account, &t.category)

		if err != nil {
			return err
		}

		row := []string{t.description, t.date, t.account, t.category,
			strconv.FormatFloat(t.amount, 'f', 2, 32), t.transaction_type}
		InsertRows(columns, i, row, t)
	}

	defer rows.Close()
	defer db.Close()

	return nil
}

func UpdateTransaction(t s.Transaction, row int) error {
	if err := t.isEmpty(); err != nil {
		return err
	}

	if db, err := sql.Open("sqlite3", "./database.db"); err != nil {
		return err
	}

	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(Transaction)

	query := `Update Transactions SET account_id = ?, category_id = ?, 
	transaction_type = ?, date = ?, amount = ?, description = ? WHERE id = ?`

	_, err = db.Exec(query, t.account_id, t.category_id, t.transaction_type, t.date, strconv.FormatFloat(t.amount, 'f', 2, 32), t.description, transaction.id)

	if err != nil {
		return err
	}

	t.id = transaction.id

	data := []string{t.date, t.transaction_type, t.account, t.category,
		strconv.FormatFloat(t.amount, 'f', 2, 32), t.description}

	UpdateRows(columns, row, data, t)
	AccountsList()

	pages.RemovePage("Form")
	defer db.Close()

	return nil
}

func AddTransaction(t s.Transaction, newRow int) {
	if err := t.isEmpty(); err != nil {
		return err
	}

	if db, err := sql.Open("sqlite3", "./database.db"); err != nil {
		return err
	}

	var result sql.Result
	if t.to_account_id.Valid && t.to_amount.Valid {
		query := `INSERT INTO Transactions (account_id, category_id, transaction_type, date, amount, description, to_amount, to_account_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(query, t.account_id, t.category_id, t.transaction_type, t.date, strconv.FormatFloat(t.amount, 'f', 2, 32), t.description, strconv.FormatFloat(t.to_amount.Float64, 'f', 2, 32), t.to_account_id.Int64)
	} else {
		query := `INSERT INTO Transactions (account_id, category_id, transaction_type, date, amount, description) VALUES (?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(query, t.account_id, t.category_id, t.transaction_type, t.date, strconv.FormatFloat(t.amount, 'f', 2, 32), t.description)
	}

	if err != nil {
		return err
	}

	if created_id, err := result.LastInsertId(); err != nil {
		return err
	}

	t.id = created_id

	row := []string{t.description, t.date, t.account, t.category,
		strconv.FormatFloat(t.amount, 'f', 2, 32), t.transaction_type}

	InsertRows(columns, newRow, row, t)
	AccountsList()
	pages.RemovePage("Form")

	defer db.Close()
	return nil
}

func DeleteTransaction() error {
	row, _ := table.GetSelection()

	if table.GetRowCount() <= 1 {
		return errors.New("Table Empty")
	}

	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(s.Transaction)

	if db, err := sql.Open("sqlite3", "./database.db"); err != nil {
		return err
	}

	query := `DELETE FROM Transactions WHERE id = ?`

	if _, err = db.Exec(query, transaction.id); err != nil {
		return err
	}

	defer db.Close()
	AccountsList()
	table.RemoveRow(row)
	return nil
}

func (t s.Transaction) isEmpty() error {
	if t.account_id == 0 || t.category_id == 0 || t.transaction_type == "" || t.date == "" || t.amount == 0 {
		return errors.New(fmt.Sprintf("%+v", t))
	}
	_, err := time.Parse("2006-01-02", t.date)
	if err != nil {
		return errors.New("Allowed date format (YYYY-MM-DD)")
	}
	return nil
}

func IsTransfer(form *tview.Form, cell *tview.TableCell, t *s.Transaction) bool {
	transfer := false
	tran_type := "debit"

	if cell != nil {
		reference := cell.GetReference().(s.Transaction)
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

func TranTypes(form *tview.Form, label string, t *s.Transaction) {
	types := []string{"debit", "credit"}
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

func ToAccount(form *tview.Form, cell *tview.TableCell, t *s.Transaction) {
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
		added(text, "to_amount", t)
	})
}

func LoadTransactions(request string, source *s.Source) {
	source.Table.Clear()

	source.Table.SetTitle("Transactions")

	SelectTransactions(request)

	source.Table.Select(1, 1).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		forms.Fill(len(source.Columns), row, false)

		source.Pages.AddPage("Form", m.Modal(source.Form, 30, 50), true, true)
	})

	source.Table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}

