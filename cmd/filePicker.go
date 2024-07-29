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
)

func FilePicker(path string) {
	pages := source.Pages
	tree := source.CategoryTree
	accounts := source.AccountList

	defer m.ErrorModal(pages, source.Modal)

	folder, err := os.Open(path)
	check(err)

	files, err := folder.Readdir(0)
	check(err)

	var pdfFiles []string
	for _, v := range files {
		slice := strings.Split(v.Name(), ".")
		if v.IsDir() || slice[len(slice)-1] != "pdf" {
			continue
		}

		pdfFiles = append(pdfFiles, path+"/"+v.Name())
	}

	selected := func(path string, source *s.Source) {
		defer m.ErrorModal(source.Pages, source.Modal)
		if tree.GetRowCount() <= 0 || accounts.GetItemCount() <= 0 {
			check(errors.New("Account and category must be created"))
		}
		var transaction s.Transaction
		selectForm(path, &transaction)

		pages.RemovePage("Files")
	}

	m.FileTable(source, "Files", pdfFiles, selected)
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
