package cmd

import (
	"database/sql"
	"log"
	"fmt"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

var (
	accounts = tview.NewList()
)

type account_type struct {
	id        int64
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
		account_id := a.id
		accounts.AddItem(a.title, a.currency, 0, func() { SelectedAccount(account_id) })
	}
	
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

	result, err := db.Exec(query, a.title, a.currency, a.balance)
	check(err)

	created_id, _ := result.LastInsertId()
	
	db.Close()
	
	accounts.AddItem(a.title, a.currency, 0, func() { SelectedAccount(created_id) })
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

func SelectedAccount(id int64) {
	request := fmt.Sprintf(`
		SELECT 
		Transactions.id, transaction_type, date, amount, Transactions.balance, Accounts.title as account, Categories.title as category
		FROM Transactions
		INNER JOIN Categories ON Categories.id = Transactions.category_id
		INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE account_id = %v
	`, id)
	table.Clear()
	form = Table(request)
}