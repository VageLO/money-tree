package cmd

import (
	"database/sql"
	"log"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

var (
	accounts = tview.NewList()
)

type account_type struct {
	id        int
	title     string
}

func AccountsList() *tview.List {

	account_titles, _ := SelectAccounts()
	
	accounts.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")
		
	for _, title := range account_titles {
		accounts.AddItem(title, "", '1', nil)
	}
	accounts.ShowSecondaryText(false)
	
	return accounts
}

func RenameAccount(text string, list *tview.List) {
	list.SetItemText(list.GetCurrentItem(), text, "")
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

func SelectAccounts() ([]string, []account_type) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	root_accounts, err := db.Query(`SELECT id, title FROM Accounts`)
	check(err)

	var account_titles []string
	var account_types []account_type
	
	for root_accounts.Next() {
		var a account_type
		if err := root_accounts.Scan(&a.id, &a.title); err != nil {
			log.Fatal(err)
		}
		account_titles = append(account_titles, a.title)
		account_types = append(account_types, a)
	}

	defer root_accounts.Close()
	db.Close()
	return account_titles, account_types
}