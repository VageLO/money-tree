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
	currency string
	balance float64
}

func AccountsList() *tview.List {

	_, account_types := SelectAccounts()
	
	accounts.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")
		
	for _, a := range account_types {
		accounts.AddItem(a.title, a.currency, '1', nil)
	}
	//accounts.ShowSecondaryText(false)
	
	return accounts
}

func RenameAccount(text string, list *tview.List) {
	list.SetItemText(list.GetCurrentItem(), text, "")
}

func RemoveAccount() {
	accounts.RemoveItem(accounts.GetCurrentItem())
}

func AddAccount(a *account_type) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	query := `
	INSERT INTO Accounts (title, currency, balance) VALUES (?, ?, ?)`

	_, err = db.Exec(query, a.title, a.currency, a.balance)
	check(err)

	db.Close()
	
	accounts.AddItem(a.title, a.currency, '3', nil)
	pages.RemovePage("Dialog")
}

func SelectAccounts() ([]string, []account_type) {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	root_accounts, err := db.Query(`SELECT * FROM Accounts`)
	check(err)

	var account_titles []string
	var account_types []account_type
	
	for root_accounts.Next() {
		var a account_type
		if err := root_accounts.Scan(&a.id, &a.title, &a.currency, &a.balance); err != nil {
			log.Fatal(err)
		}
		account_titles = append(account_titles, a.title)
		account_types = append(account_types, a)
	}

	defer root_accounts.Close()
	db.Close()
	return account_titles, account_types
}