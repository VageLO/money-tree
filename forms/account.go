package forms

import (
	"strconv"
	"strings"
	s "main/structs"
	"main/action"
	
	"github.com/rivo/tview"
)

func FillListForm(list *tview.List, form *tview.Form) {
	form.Clear(true)

	formStyle("Account Information")
	title, second := list.GetItemText(list.GetCurrentItem())
	split := strings.Split(second, " ")
	balance := split[0]
	currency := split[1]
	
	form.AddInputField("Title: ", title, 0, nil, func(text string) { action.RenameAccount(text, "title", list) })
	
	form.AddInputField("Currency: ", currency, 0, nil, func(text string) { action.RenameAccount(text, "currency", list) })
	
	form.AddInputField("Balance: ", balance, 0, nil, func(text string) { action.RenameAccount(text, "balance", list) })
}

func FormAddAccount(form *tview.Form, pages *tview.Pages) {
	form.Clear(true)
	formStyle("Add Account")
	
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
		defer ErrorModal()

		balance, err := strconv.ParseFloat(text, 64);
		check(err)
		a.Balance = balance
	})
	form.AddButton("Add", func() { AddAccount(&a) })
	pages.AddPage("Form", Modal(form, 30, 50), true, true)
}

func FormRenameAccount(accounts *tview.List, form *tview.Form, pages *tview.Pages) {
	if accounts.GetItemCount() <= 0 {
		return
	}
	FillListForm(accounts, form)
	pages.AddPage("Form", Modal(form, 30, 50), true, true)
}

func SelectedAccount(option string, optionIndex int, a_types []s.Account, t *s.Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.title != option {
		return
	}
	t.account_id = selected_a.id
	t.account = selected_a.title
}