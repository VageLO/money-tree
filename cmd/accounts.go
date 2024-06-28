package cmd

import (
	"database/sql"
	"log"
	"fmt"
	"strings"
	"errors"
	"strconv"
	
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

	accounts.Clear()
	
	_, account_types := SelectAccounts()
	
	accounts.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")

	for _, a := range account_types {
		account_id := a.id
		second_title := fmt.Sprintf("%v %v", strconv.FormatFloat(a.balance, 'f', 2, 32), a.currency)
		accounts.AddItem(a.title, second_title, 0, func() { SelectedAccount(account_id) })
	}
	
	return accounts
}

func RenameAccount(value, field string, list *tview.List) {
	defer ErrorModal()
	if value == "" {
		check(errors.New("Can't be empty"))
	}
	selected_item := list.GetCurrentItem()
	
	title, second := list.GetItemText(selected_item)
	split := strings.Split(second, " ")
	balance := split[0]
	currency := split[1]
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	query := fmt.Sprintf(`UPDATE Accounts SET %v = ? WHERE title = ?`, field)

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
	
	defer db.Close()
	list.SetItemText(selected_item, title, balance + " " + currency)
}

func RemoveAccount() {
	defer ErrorModal()
	
	if accounts.GetItemCount() <= 0 {
		return
	}
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	selected_account := accounts.GetCurrentItem()
	title, _ := accounts.GetItemText(selected_account)
	
	query := `
	DELETE FROM Accounts WHERE title = ?`

	_, err = db.Exec(query, title)
	check(err)
	
	accounts.RemoveItem(selected_account)
	defer db.Close()
}

func AddAccount(a *account_type) {
	defer ErrorModal()
	check(a.isEmpty())
	
	db, err := sql.Open("sqlite3", "./database.db")
	check(err)
	
	query := `
	INSERT INTO Accounts (title, currency, balance) VALUES (?, ?, ?)`
	
	result, err := db.Exec(query, a.title, a.currency, a.balance)
	check(err)
	
	created_id, _ := result.LastInsertId()
	balance := fmt.Sprintf("%v %v", a.balance, a.currency)
	accounts.AddItem(a.title, balance, 0, func() { SelectedAccount(created_id) })
	pages.RemovePage("Form")
	defer db.Close()
}

func SelectAccounts() ([]string, []account_type) {
	defer ErrorModal()
	
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

	defer db.Close()
	defer root_accounts.Close()
	return account_titles, account_types
}

func SelectedAccount(id int64) {
	request := fmt.Sprintf(`SELECT Transactions.*, Accounts.title, Categories.title FROM Transactions INNER JOIN Categories ON Categories.id = Transactions.category_id INNER JOIN Accounts ON Accounts.id = Transactions.account_id WHERE account_id = %v`, id)
	FillTable(request)
}

func (a account_type) isEmpty() error {
	if a.title == "" || a.currency == "" {
		return errors.New("Empty field or can't be zero")
	}
	return nil
}