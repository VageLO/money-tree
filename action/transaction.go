package action

import (
	"database/sql"
	"errors"
	m "main/modal"
	s "main/structs"
	"strconv"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func LoadTransactions(request string, source *s.Source) {
	source.Table.Clear()
	source.Table.SetBorder(true)
	source.Table.SetTitle("Transactions")

	SelectTransactions(request, source)

	source.Table.Select(1, 1).SetFixed(0, 1).SetSelectedFunc(func(row int, column int) {
		if row <= 0 {
			return
		}
		FillForm(len(source.Columns), row, false, source)

		source.Pages.AddPage("Form", m.Modal(source.Form, 30, 50), true, true)
	})

	source.Table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}

func SelectTransactions(request string, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	rows, err := db.Query(request)
	check(err)

	for i, columnTitle := range source.Columns {
		InsertCell(s.Cell{
			Row:        0,
			Column:     i,
			Text:       columnTitle,
			Selectable: false,
			Color:      tcell.ColorYellow,
		}, source.Table)
	}

	for i := 1; rows.Next(); i++ {
		var t s.Transaction

		err := rows.Scan(
			&t.Id,
			&t.AccountId,
			&t.ToAccountId,
			&t.CategoryId,
			&t.TransactionType,
			&t.Date,
			&t.Amount,
			&t.ToAmount,
			&t.Description,
			&t.Account,
			&t.ToAccount,
			&t.Category,
		)
		check(err)

		row := []string{t.Description, t.Date, t.Account, t.Category,
			strconv.FormatFloat(t.Amount, 'f', 2, 32), t.TransactionType}

		InsertRows(s.Row{
			Columns:     source.Columns,
			Index:       i,
			Data:        row,
			Transaction: t,
		}, source.Table)
	}

	defer rows.Close()
	defer db.Close()
}

func UpdateTransaction(t s.Transaction, row int, source *s.Source) {
	pages := source.Pages
	modal := source.Modal
	table := source.Table

	defer m.ErrorModal(pages, modal)

	check(t.IsEmpty())

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(s.Transaction)

	if t.ToAccountId.Valid && t.ToAmount.Valid {
		query := `Update Transactions SET account_id = ?, category_id = ?, 
	transaction_type = ?, date = ?, amount = ?, description = ?, to_account_id = ?, to_amount = ? WHERE id = ?`

		_, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			strconv.FormatFloat(t.Amount, 'f', 2, 32),
			t.Description,
			t.ToAccountId.Int64,
			strconv.FormatFloat(t.ToAmount.Float64, 'f', 2, 32),
			transaction.Id,
		)
	} else {
		query := `Update Transactions SET account_id = ?, category_id = ?, 
	transaction_type = ?, date = ?, amount = ?, description = ?, to_account_id = NULL, to_amount = NULL WHERE id = ?`

		_, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			strconv.FormatFloat(t.Amount, 'f', 2, 32),
			t.Description, transaction.Id,
		)
	}

	check(err)

	t.Id = transaction.Id

	data := []string{t.Description, t.Date, t.Account, t.Category,
		strconv.FormatFloat(t.Amount, 'f', 2, 32), t.TransactionType}

	UpdateRows(s.Row{
		Columns:     source.Columns,
		Index:       row,
		Data:        data,
		Transaction: t,
	}, source.Table)

	LoadAccounts(source)

	pages.RemovePage("Form")
	defer db.Close()
}

func AddTransaction(t s.Transaction, newRow int, source *s.Source) {
	pages := source.Pages
	modal := source.Modal

	defer m.ErrorModal(pages, modal)

	check(t.IsEmpty())

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	var result sql.Result
	if t.ToAccountId.Valid && t.ToAmount.Valid {
		query := `INSERT INTO Transactions (account_id, category_id, transaction_type, date, amount, description, to_amount, to_account_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			strconv.FormatFloat(t.Amount, 'f', 2, 32),
			t.Description,
			strconv.FormatFloat(t.ToAmount.Float64, 'f', 2, 32),
			t.ToAccountId.Int64,
		)
	} else {
		query := `INSERT INTO Transactions (account_id, category_id, transaction_type, date, amount, description) VALUES (?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			strconv.FormatFloat(t.Amount, 'f', 2, 32),
			t.Description,
		)
	}

	check(err)

	createdId, err := result.LastInsertId()
	check(err)

	t.Id = createdId

	row := []string{t.Description, t.Date, t.Account, t.Category,
		strconv.FormatFloat(t.Amount, 'f', 2, 32), t.TransactionType}

	InsertRows(s.Row{
		Columns:     source.Columns,
		Index:       newRow,
		Data:        row,
		Transaction: t,
	}, source.Table)

	LoadAccounts(source)
	pages.RemovePage("Form")

	defer db.Close()
}

func DeleteTransaction(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	table := source.Table

	row, _ := table.GetSelection()

	if table.GetRowCount() <= 1 {
		check(errors.New("Table Empty"))
	}

	cell := table.GetCell(row, 0)
	transaction := cell.GetReference().(s.Transaction)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	query := `DELETE FROM Transactions WHERE id = ?`

	_, err = db.Exec(query, transaction.Id)
	check(err)

	defer db.Close()
	LoadAccounts(source)
	table.RemoveRow(row)
}
