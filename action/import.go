package action

import (
	m "main/modal"
	"main/parser"
	s "main/structs"
	"strconv"
)

func insertIntoDb(source *s.Source, path string, t *s.Transaction) {
	transactions := parser.ParsePdf(path)
	table := source.Table

	for _, importTransaction := range transactions {
		newRow := table.GetRowCount()
		var transaction s.Transaction

		debit := []string{"Безналичная операция", "Отправление средств", "Банкомат"}
		credit := []string{"Поступление", "Получение средств", "Поступление (Credit)", "Внесение наличных"}

		transaction.AccountId = t.AccountId
		transaction.CategoryId = t.CategoryId
		transaction.Account = t.Account
		transaction.Category = t.Category

		err := false

		if importTransaction.Status != "Завершено успешно" {
			err = true
		}

		transaction.Date = importTransaction.Date

		amount, _ := strconv.ParseFloat(importTransaction.Price, 64)
		transaction.Amount = amount

		transaction.Description = importTransaction.Description

		if contains(debit, importTransaction.Typeof) {
			transaction.TransactionType = "debit"
		} else if contains(credit, importTransaction.Typeof) {
			transaction.TransactionType = "credit"
		}

		if err {
			continue
		}

		AddTransaction(transaction, newRow, source)
	}
}

func ImportForm(source *s.Source, path string) {
	form := source.Form

	var transaction s.Transaction

	form.Clear(true)
	FormStyle("Select account and category", form)

	categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`, source)
	accounts, a_types := SelectAccounts(source)

	SelectedCategory(categories[0], 0, c_types, &transaction)
	SelectedAccount(accounts[0], 0, a_types, &transaction)

	form.AddDropDown("Category", categories, 0, func(option string, optionIndex int) {
		SelectedCategory(option, optionIndex, c_types, &transaction)
	})

	form.AddDropDown("Account", accounts, 0, func(option string, optionIndex int) {
		SelectedAccount(option, optionIndex, a_types, &transaction)
	})

	form.AddButton("Import", func() { insertIntoDb(source, path, &transaction) })
	source.Pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func FileExporer(source *s.Source, pattern, pageName string) {
	tree := newTree(source, pattern, pageName)
	source.Pages.AddPage(pageName, tree, true, true)
}
