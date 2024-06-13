package cmd

import (
	"database/sql"
	"log"
	"fmt"
	"strings"
	
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
		second_title := fmt.Sprintf("%v %v", a.balance, a.currency)
		accounts.AddItem(a.title, second_title, 0, func() { SelectedAccount(account_id) })
	}
	
	return accounts
}

func RenameAccount(value, field string, list *tview.List) {
	selected_item := list.GetCurrentItem()
	
	title, second := list.GetItemText(selected_item)
	split := strings.Split(second, " ")
	balance := split[0]
	currency := split[1]
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	query := fmt.Sprintf(`UPDATE Accounts SET %v = ? WHERE title = ?`, field)
	log.Println(query)
	_, err = db.Exec(query, value, title)
	check(err)
	
	// Shame
	if field == "currency" {
		currency = value
	}
	if field == "balance" {
		balance = value
	} 
	if field == "title" {
		title = value
	}
	
	db.Close()
	list.SetItemText(selected_item, title, balance + " " + currency)
}

func RemoveAccount() {
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	selected_account := accounts.GetCurrentItem()
	title, _ := accounts.GetItemText(selected_account)
	
	query := `
	DELETE FROM Accounts WHERE title = ?`

	_, err = db.Exec(query, title)
	check(err)
	
	db.Close()
	accounts.RemoveItem(selected_account)
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