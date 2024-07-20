package action

import (
	m "main/modal"
	s "main/structs"
	"strconv"
	"strings"
)

func FillListForm(source *s.Source) {
	form := source.Form
	list := source.AccountList

	form.Clear(true)
	FormStyle("Account Information", form)

	title, second := list.GetItemText(list.GetCurrentItem())
	split := strings.Split(second, " ")
	balance := split[0]
	currency := split[1]

	form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameAccount(text, "title", list) })

	form.AddInputField("Currency: ", currency, 0, nil, func(text string) { RenameAccount(text, "currency", list) })

	form.AddInputField("Balance: ", balance, 0, nil, func(text string) { RenameAccount(text, "balance", list) })
}

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
	accounts := source.AccountList
	form := source.Form
	pages := source.Pages

	if accounts.GetItemCount() <= 0 {
		return
	}
	FillListForm(source)
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
