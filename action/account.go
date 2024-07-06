package action

import (
	"database/sql"
	"errors"
	"fmt"
	m "main/modal"
	s "main/structs"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func LoadAccounts(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	source.AccountList.Clear()

	_, account_types := SelectAccounts(source)

	source.AccountList.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")

	for _, a := range account_types {
		account_id := a.Id
		second_title := fmt.Sprintf("%v %v", strconv.FormatFloat(a.Balance, 'f', 2, 32), a.Currency)
		source.AccountList.AddItem(a.Title, second_title, 0, func() { WhereAccount(account_id, source) })
	}
}

func RenameAccount(value, field string, list *tview.List) error {
	if value == "" {
		return errors.New("Can't be empty")
	}
	selected_item := list.GetCurrentItem()

	title, second := list.GetItemText(selected_item)
	split := strings.Split(second, " ")
	balance := split[0]
	currency := split[1]

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`UPDATE Accounts SET %v = ? WHERE title = ?`, field)

	if _, err = db.Exec(query, value, title); err != nil {
		return err
	}

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
	list.SetItemText(selected_item, title, balance+" "+currency)

	return nil
}

func RemoveAccount(accounts *tview.List) error {
	if accounts.GetItemCount() <= 0 {
		return errors.New("GetItemCount on account list is <= 0")
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	selected_account := accounts.GetCurrentItem()
	title, _ := accounts.GetItemText(selected_account)

	query := `
	DELETE FROM Accounts WHERE title = ?`

	if _, err = db.Exec(query, title); err != nil {
		return err
	}

	accounts.RemoveItem(selected_account)

	defer db.Close()
	return nil
}

func AddAccount(a *s.Account, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	check(a.IsEmpty())

	accounts := source.AccountList

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	query := `
	INSERT INTO Accounts (title, currency, balance) VALUES (?, ?, ?)`

	result, err := db.Exec(query, a.Title, a.Currency, a.Balance)
	check(err)

	created_id, _ := result.LastInsertId()
	balance := fmt.Sprintf("%v %v", a.Balance, a.Currency)

	accounts.AddItem(a.Title, balance, 0, func() { WhereAccount(created_id, source) })

	defer db.Close()
}

func SelectAccounts(source *s.Source) ([]string, []s.Account) {
	defer m.ErrorModal(source.Pages, source.Modal)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	root_accounts, err := db.Query(`SELECT * FROM Accounts`)
	check(err)

	var account_titles []string
	var account_types []s.Account

	for root_accounts.Next() {
		var a s.Account
		err := root_accounts.Scan(&a.Id, &a.Title, &a.Currency, &a.Balance)
		check(err)

		account_titles = append(account_titles, a.Title)
		account_types = append(account_types, a)
	}

	defer db.Close()
	defer root_accounts.Close()
	return account_titles, account_types
}

func WhereAccount(id int64, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	query, err := os.ReadFile("./sql/Select_On_Transactions_Where_AccountID.sql")
	check(err)

	str_id := strconv.FormatInt(id, 10)

	request := string(query)
	request = strings.ReplaceAll(request, "?", str_id)
	LoadTransactions(request, source)
}
