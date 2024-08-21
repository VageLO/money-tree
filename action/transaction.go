package action

import (
	"database/sql"
	"errors"
	"fmt"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"
)

func check(err error) {
	if err == nil {
		return
	}

	dir, e := os.UserConfigDir()
	if e != nil {
		log.Fatalln(e)
	}
	configPath := filepath.Join(dir, "money-tree")

	// Create money-tree directory in UserConfigDir
	if e = os.Mkdir(configPath, 0750); e != nil && !os.IsExist(e) {
		log.Fatalln(e)
	}

	// Create log file
	logFile, e := os.OpenFile(filepath.Join(configPath, "tree.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening log file: %v\n", e)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println(err)
	panic(err)
}

var SelectedRows []int

func LoadTransactions(request string, source *s.Source) {
	table := source.Table
	table.Clear()
	table.SetBorder(true)
	table.SetTitle("Transactions")

	SelectTransactions(request, source)

	table.Select(0, 0).SetFixed(0, 1).SetSelectedFunc(func(row int, column int) {
		if table.GetRowCount() <= 1 {
			return
		}
		FilledForm(row, source)

		source.Pages.AddPage("Form", m.Modal(source.Form, 30, 50), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}

func SelectTransactions(request string, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	db, err := sql.Open("sqlite3", source.Config.Database)
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

		amount := strconv.FormatFloat(t.Amount, 'f', 2, 32)
		row := []string{t.Description, t.Date, t.Account, t.Category,
			amount, t.TransactionType}

		if t.ToAccountId.Valid {
			row[2] = fmt.Sprintf("%v > %v", t.Account, t.ToAccount.String)
			if t.ToAmount.Valid && t.ToAmount.Float64 != 0 {
				toAmount := strconv.FormatFloat(t.ToAmount.Float64, 'f', 2, 32)
				row[4] = amount + " > " + toAmount
			}
		}

		InsertRows(s.Row{
			Columns:   source.Columns,
			Index:     i,
			Data:      row,
			Reference: t,
		}, source.Table)
	}

	defer rows.Close()
	defer db.Close()
}

func UpdateTransaction(t s.Transaction, row int, source *s.Source, removePreviousAttachments bool) {
	pages := source.Pages
	modal := source.Modal

	defer m.ErrorModal(pages, modal)

	check(t.IsEmpty())

	db, err := sql.Open("sqlite3", source.Config.Database)
	check(err)

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
			t.Id,
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
			t.Amount,
			t.Description,
			t.Id,
		)
	}

	check(err)

	updateAttachments(source, source.Attachments, t.Id, removePreviousAttachments)

	amount := strconv.FormatFloat(t.Amount, 'f', 2, 32)
	data := []string{t.Description, t.Date, t.Account, t.Category,
		amount, t.TransactionType}

	if t.ToAccountId.Valid {
		data[2] = fmt.Sprintf("%v > %v", t.Account, t.ToAccount.String)
		if t.ToAmount.Valid && t.ToAmount.Float64 != 0 {
			toAmount := strconv.FormatFloat(t.ToAmount.Float64, 'f', 2, 32)
			data[4] = amount + " > " + toAmount
		}
	}

	UpdateRows(s.Row{
		Columns:   source.Columns,
		Index:     row,
		Data:      data,
		Reference: t,
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

	db, err := sql.Open("sqlite3", source.Config.Database)
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
			t.Amount,
			t.Description,
		)
	}

	check(err)

	createdId, err := result.LastInsertId()
	check(err)

	t.Id = createdId

	// Load Attachments if exist
	if len(source.Attachments) > 0 {
		addAttachments(source, t.Id, source.Attachments)
	}

	row := []string{t.Description, t.Date, t.Account, t.Category,
		strconv.FormatFloat(t.Amount, 'f', 2, 32), t.TransactionType}

	InsertRows(s.Row{
		Columns:   source.Columns,
		Index:     newRow,
		Data:      row,
		Reference: t,
	}, source.Table)

	LoadAccounts(source)
	pages.RemovePage("Form")

	defer db.Close()
}

func DeleteTransaction(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	table := source.Table

	if table.GetRowCount() <= 1 {
		check(errors.New("Table Empty"))
	}

	remove := func(row int) {
		cell := table.GetCell(row, 0)
		transaction := cell.GetReference().(s.Transaction)

		db, err := sql.Open("sqlite3", source.Config.Database)
		check(err)

		query := `DELETE FROM Transactions WHERE id = ?`

		_, err = db.Exec(query, transaction.Id)
		check(err)

		deleteAttachments(source, findAttachments(source, transaction.Id))

		defer db.Close()
		LoadAccounts(source)
	}

	if len(SelectedRows) <= 0 {
		row, _ := table.GetSelection()
		remove(row)
		table.RemoveRow(row)
		return
	}

	for _, row := range SelectedRows {
		remove(row)
	}

	slices.Sort(SelectedRows)
	for i := len(SelectedRows) - 1; i >= 0; i-- {
		table.RemoveRow(SelectedRows[i])
	}
	SelectedRows = []int{}
}

func addAttachments(source *s.Source, id int64, attachments []string) {
	defer m.ErrorModal(source.Pages, source.Modal)
	currentIndex := len(findAttachments(source, id))

	for index, att := range attachments {
		index += currentIndex
		bytes, err := os.ReadFile(att)
		check(err)

		file, err := os.OpenFile(att, os.O_RDONLY, 0644)
		check(err)

		extension := filepath.Ext(file.Name())

		path := filepath.Join(source.Config.Attachments, fmt.Sprintf("%v_%v%v", id, index, extension))
		err = os.WriteFile(path, bytes, 0644)
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

	folder := source.Config.Attachments
	files, err := os.ReadDir(folder)

	check(err)

	var attachments []string

	for _, file := range files {
		value, err := strconv.ParseInt(strings.Split(file.Name(), "_")[0], 10, 64)
		check(err)

		if file.IsDir() {
			continue
		} else if value != id {
			continue
		}

		attachments = append(attachments, filepath.Join(folder, file.Name()))
	}

	return attachments
}

func updateAttachments(source *s.Source, newAttachments []string, id int64, remove bool) {
	currentAttachments := findAttachments(source, id)

	var addArray []string
	var deleteArray []string

	if remove {
		// Check deleted attachments and delete them
		for _, currentAttachment := range currentAttachments {
			if exist, _ := Contains(newAttachments, currentAttachment); !exist {
				deleteArray = append(deleteArray, currentAttachment)
			}
		}
		deleteAttachments(source, deleteArray)
	}

	renameAttachments(source, id)

	for _, newAttachment := range newAttachments {
		if exist, _ := Contains(currentAttachments, newAttachment); !exist {
			addArray = append(addArray, newAttachment)
		}
	}
	addAttachments(source, id, addArray)
}

func renameAttachments(source *s.Source, id int64) {
	defer m.ErrorModal(source.Pages, source.Modal)

	attachments := findAttachments(source, id)

	for index, attachment := range attachments {
		extension := filepath.Ext(attachment)
		err := os.Rename(
			attachment,
			filepath.Join(
				filepath.Dir(attachment),
				fmt.Sprintf("%v_%v%v", id, index, extension),
			),
		)
		check(err)
	}
}
