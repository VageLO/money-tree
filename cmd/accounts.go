package cmd

import (
	"github.com/rivo/tview"
)

var (
	accounts = tview.NewList()
)

func AccountsList() *tview.List {
	//var tableData []string

	accounts.ShowSecondaryText(false).
		AddItem("Alfa Bank", "123", '1', nil).
		AddItem("BNB", "123", '2', func() {
			//FillTable(tableData)
			table.SetBorder(true).SetTitle("Accounts")
		})
	accounts.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")

	return accounts
}

func RenameAccount() {
	FillTreeAndListForm(nil, accounts)
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func RemoveAccount() {
	accounts.RemoveItem(accounts.GetCurrentItem())
}

func AddAccount() {
	form.Clear(true)
	var input_text string
	form.AddInputField("Title: ", "", 0, nil, func(text string) {
		input_text = text
	})
	form.AddButton("Add", func() {
		accounts.AddItem(input_text, "", '3', nil)
		pages.RemovePage("Dialog")
	})
	pages.AddPage("Dialog", Dialog(form), true, true)
}
