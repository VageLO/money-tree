package action

import (
	m "github.com/VageLO/money-tree/modal"
	"github.com/VageLO/money-tree/parser"
	s "github.com/VageLO/money-tree/structs"
	"strconv"
)

func insertIntoDb(source *s.Source, path string, t *s.Transaction) {
    defer m.ErrorModal(source.Pages, source.Modal)
	err, transactions := parser.ParsePdf(path)
    if err != nil {
        source.Pages.RemovePage("Form")
    }
    check(err)

	table := source.Table

	for _, importTransaction := range transactions {
		newRow := table.GetRowCount()
		var transaction s.Transaction

		Withdrawal := []string{"Безналичная операция", "Отправление средств", "Банкомат"}
		Deposit := []string{"Поступление", "Получение средств", "Поступление (Credit)", "Внесение наличных"}

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

		if exist, _ := Contains(Withdrawal, importTransaction.Typeof); exist {
			transaction.TransactionType = "Withdrawal"
		} else if exist, _ := Contains(Deposit, importTransaction.Typeof); exist {
			transaction.TransactionType = "Deposit"
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

func Contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}
	return false, -1
}

func FileExporer(source *s.Source, pattern, pageName string) {
	tree := newTree(source, pattern, pageName)
	source.Pages.AddPage(pageName, tree, true, true)
}
