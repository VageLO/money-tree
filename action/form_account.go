package action

import (
	m "main/modal"
	s "main/structs"
	"strconv"
)

func FormAddAccount(source *s.Source) {
	form := source.Form
	pages := source.Pages

	form.Clear(true)
	FormStyle("Add Account", form)

	var a = s.Account{}

	form.AddInputField("Title: ", "", 0, nil, func(text string) {
		a.Title = text
	})
	form.AddInputField("Currency: ", "", 0, nil, func(text string) {
		a.Currency = text
	})
	form.AddInputField("Balance: ", "0", 0, nil, func(text string) {
		if text == "" {
			a.Balance = 0
			return
		}
		defer m.ErrorModal(pages, source.Modal)

		balance, err := strconv.ParseFloat(text, 64)
		check(err)
		a.Balance = balance
	})
	form.AddButton("Add", func() { AddAccount(&a, source) })

	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func FormRenameAccount(source *s.Source) {

	accountList := source.AccountList
	form := source.Form
	pages := source.Pages

	title, _ := accountList.GetItemText(accountList.GetCurrentItem())

	if accountList.GetItemCount() <= 1 || title == "All Transactions" {
		return
	}

	form.Clear(true)
	FormStyle("Account Details", form)

	var a s.Account

	_, accounts := SelectAccounts(source)
	for _, account := range accounts {
		if account.Title == title {
			a = account
		}
	}
	form.AddInputField("Title: ", a.Title, 0, nil, func(text string) {
		a.Title = text
	})

	form.AddInputField("Currency: ", a.Currency, 0, nil, func(text string) {
		a.Currency = text
	})

	form.AddInputField("Balance: ", strconv.FormatFloat(a.Balance, 'f', 2, 32), 0, nil, func(text string) {
		if text == "" {
			a.Balance = 0
			return
		}
		defer m.ErrorModal(pages, source.Modal)

		balance, err := strconv.ParseFloat(text, 64)
		check(err)
		a.Balance = balance
	})

	form.AddButton("Save", func() { RenameAccount(a, source) })

	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func SelectedAccount(option string, optionIndex int, a_types []s.Account, t *s.Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.Title != option {
		return
	}
	t.AccountId = selected_a.Id
	t.Account = selected_a.Title
}
