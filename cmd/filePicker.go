package cmd

import (
	"errors"
	"main/action"
	m "main/modal"
	"main/parser"
	s "main/structs"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

func FilePicker(path string) {
	pages := source.Pages
	file_table := source.FileTable
	tree := source.CategoryTree
	accounts := source.AccountList

	defer m.ErrorModal(source.Pages, source.Modal)

	file_table.Clear()
	file_table.SetTitle("Pdf files")
	file_table.SetSelectable(true, false)

	folder, err := os.Open(path)
	check(err)

	files, err := folder.Readdir(0)
	check(err)

	count := 0
	for _, v := range files {
		slice := strings.Split(v.Name(), ".")
		if v.IsDir() || slice[len(slice)-1] != "pdf" {
			continue
		}

		tableCell := tview.NewTableCell(v.Name())
		tableCell.SetReference(struct {
			path string
		}{path + "/" + v.Name()})
		tableCell.SetSelectable(true)

		file_table.SetCell(count, 0, tableCell)
		count++
	}

	file_table.SetSelectedFunc(func(row int, column int) {
		if tree.GetRowCount() <= 0 || accounts.GetItemCount() <= 0 {
			check(errors.New("Account and category must be created"))
		}

		cell := file_table.GetCell(row, column)

		ref := cell.GetReference().(struct {
			path string
		})

		var transaction s.Transaction
		selectForm(ref.path, &transaction)

		pages.RemovePage("Files")
	})
	x, _, _, _ := file_table.GetRect()
	pages.AddPage("Files", m.Modal(file_table, 30, x), true, true)
}

func insertIntoDb(path string, t *s.Transaction) {
	transactions := parser.ParsePdf(path)
	table := source.Table

	for _, import_transaction := range transactions {
		newRow := table.GetRowCount()
		var transaction s.Transaction

		debit := []string{"Безналичная операция", "Отправление средств", "Банкомат"}
		credit := []string{"Поступление", "Получение средств", "Поступление (Credit)", "Внесение наличных"}

		transaction.AccountId = t.AccountId
		transaction.CategoryId = t.CategoryId
		transaction.Account = t.Account
		transaction.Category = t.Category

		r := reflect.ValueOf(&import_transaction).Elem()
		rt := r.Type()

		err_status := false

		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			rv := reflect.ValueOf(&import_transaction)
			value := reflect.Indirect(rv).FieldByName(field.Name)

			if field.Name == "status" && value.String() != "Завершено успешно" {
				err_status = true
			}
			if field.Name == "date" {
				transaction.Date = value.String()
			}
			if field.Name == "price" {
				amount, _ := strconv.ParseFloat(value.String(), 64)
				transaction.Amount = amount
			}
			if field.Name == "description" {
				transaction.Description = value.String()
			}
			if field.Name == "typeof" && contains(debit, value.String()) {
				transaction.TransactionType = "debit"
			} else if field.Name == "typeof" && contains(credit, value.String()) {
				transaction.TransactionType = "credit"
			}
		}

		if err_status {
			continue
		}

		action.AddTransaction(transaction, newRow, source)
	}
}

func selectForm(path string, t *s.Transaction) {
	form := source.Form
	pages := source.Pages

	form.Clear(true)
	action.FormStyle("Select account and category", form)

	categories, c_types, _ := action.SelectCategories(`SELECT * FROM Categories`, source)
	accounts, a_types := action.SelectAccounts(source)

	action.SelectedCategory(categories[0], 0, c_types, t)
	action.SelectedAccount(accounts[0], 0, a_types, t)

	form.AddDropDown("Category", categories, 0, func(option string, optionIndex int) {
		action.SelectedCategory(option, optionIndex, c_types, t)
	})

	form.AddDropDown("Account", accounts, 0, func(option string, optionIndex int) {
		action.SelectedAccount(option, optionIndex, a_types, t)
	})

	form.AddButton("Import", func() { insertIntoDb(path, t) })
	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
