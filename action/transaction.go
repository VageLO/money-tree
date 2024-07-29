package action

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	m "main/modal"
	s "main/structs"
	"os"
	"strconv"
	"strings"

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

	if t.ToAccountId.Valid {
		query := `Update Transactions SET account_id = ?, category_id = ?, 
	transaction_type = ?, date = ?, amount = ?, description = ?, to_account_id = ?, to_amount = ? WHERE id = ?`

		_, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			t.Amount,
			t.Description,
			t.ToAccountId.Int64,
			t.ToAmount.Float64,
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

	attachments := findAttachments(source, t.Id)

	if len(attachments) >= 0 {
		compareAttachments(source, attachments, source.Attachments, t.Id)
	} else {
		Attachments(source, t.Id, source.Attachments)
	}

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
	if t.ToAccountId.Valid {
		query := `INSERT INTO Transactions (account_id, category_id, transaction_type, date, amount, description, to_amount, to_account_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		result, err = db.Exec(
			query,
			t.AccountId,
			t.CategoryId,
			t.TransactionType,
			t.Date,
			t.Amount,
			t.Description,
			t.ToAmount.Float64,
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

	// Load Attachments if exist
	if len(source.Attachments) > 0 {
		Attachments(source, t.Id, source.Attachments)
	}

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

	deleteAttachments(source, findAttachments(source, transaction.Id))

	defer db.Close()
	LoadAccounts(source)
	table.RemoveRow(row)
}

func Attachments(source *s.Source, id int64, attachments []string) {
	defer m.ErrorModal(source.Pages, source.Modal)

	for index, att := range attachments {
		bytes, err := os.ReadFile(att)
		check(err)

		file, err := os.OpenFile(att, os.O_RDONLY, 0644)
		check(err)

		extension := strings.Split(file.Name(), ".")

		err = os.WriteFile(fmt.Sprintf("./attachments/%v_%v.%v", id, index, extension[len(extension)-1]), bytes, 0644)
		check(err)
	}
}

func deleteAttachments(source *s.Source, attachments []string) {
	defer m.ErrorModal(source.Pages, source.Modal)

	for _, attachment := range attachments {
		err := os.Remove(attachment)
		check(err)
	}
}

func findAttachments(source *s.Source, id int64) []string {
	defer m.ErrorModal(source.Pages, source.Modal)

	folder := "./attachments"
	files, err := os.ReadDir(folder)

	check(err)

	var attachments []string

	for _, file := range files {
		if file.IsDir() {
			continue
		} else if !strings.Contains(file.Name(), fmt.Sprintf("%v", id)) {
			continue
		}

		attachments = append(attachments, folder+"/"+file.Name())
	}

	return attachments
}

func compareAttachments(source *s.Source, existedAttachments []string, newAttachments []string, id int64) {
	defer m.ErrorModal(source.Pages, source.Modal)

	existedBuffs := makeBuffs(source, existedAttachments)
	newBuffs := makeBuffs(source, newAttachments)

	for i, existedBuff := range existedBuffs {
		for j, newBuff := range newBuffs {
			res := bytes.Compare(existedBuff, newBuff)
			if res == 0 {
				continue
			}
			deleteAttachments(source, []string{existedAttachments[i]})
			Attachments(source, id, []string{newAttachments[j]})
		}
	}
}

func makeBuffs(source *s.Source, attachments []string) [][]byte {
	defer m.ErrorModal(source.Pages, source.Modal)

	var buffs [][]byte

	for _, attachment := range attachments {
		buff, err := os.ReadFile(attachment)
		check(err)

		buffs = append(buffs, buff)
	}

	return buffs
}
