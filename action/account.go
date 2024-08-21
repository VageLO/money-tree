package action

import (
	"database/sql"
	"fmt"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func LoadAccounts(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	source.AccountList.Clear()

	_, accountTypes := SelectAccounts(source)

	source.AccountList.
		SetBorderPadding(0, 0, 2, 0).
		SetBorder(true).
		SetTitle("Accounts")

	source.AccountList.AddItem("All Transactions", "----------------", 0, func() { LoadTransactions(s.Transactions, source) })

	for _, a := range accountTypes {
		accountId := a.Id
		secondTitle := fmt.Sprintf("%v %v", strconv.FormatFloat(a.Balance, 'f', 2, 32), a.Currency)
		source.AccountList.AddItem(a.Title, secondTitle, 0, func() { WhereAccount(accountId, source) })
	}
}

func RenameAccount(a s.Account, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	pages := source.Pages
	list := source.AccountList

	selectedItem := list.GetCurrentItem()
	title, _ := list.GetItemText(selectedItem)

	db, err := sql.Open("sqlite3", source.Config.Database)
	check(err)

	query := `UPDATE Accounts SET title = ?, currency = ?, balance = ? WHERE title = ?`

	_, err = db.Exec(query, a.Title, a.Currency, a.Balance, title)
	check(err)

	list.SetItemText(selectedItem, a.Title, strconv.FormatFloat(a.Balance, 'f', 2, 32)+" "+a.Currency)

	pages.RemovePage("Form")
	source.App.SetFocus(source.AccountList)

	defer db.Close()
}

func RemoveAccount(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	accounts := source.AccountList

	selectedAccount := accounts.GetCurrentItem()
	title, _ := accounts.GetItemText(selectedAccount)

	if accounts.GetItemCount() <= 1 || title == "All Transactions" {
		return
	}

	db, err := sql.Open("sqlite3", source.Config.Database+"?_foreign_keys=on")
	check(err)

	query := `DELETE FROM Accounts WHERE title = ?`

	_, err = db.Exec(query, title)
	check(err)

	accounts.RemoveItem(selectedAccount)

	newSelectedAccount := accounts.GetCurrentItem()
	title, _ = accounts.GetItemText(newSelectedAccount)
	_, accountTypes := SelectAccounts(source)

	for _, account := range accountTypes {
		if account.Title == title {
			WhereAccount(account.Id, source)
			return
		}
	}

	LoadTransactions(s.Transactions, source)
	defer db.Close()
}

func AddAccount(a *s.Account, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	check(a.IsEmpty())

	pages := source.Pages
	accounts := source.AccountList

	db, err := sql.Open("sqlite3", source.Config.Database)
	check(err)

	query := `INSERT INTO Accounts (title, currency, balance) VALUES (?, ?, ?)`

	result, err := db.Exec(query, a.Title, a.Currency, a.Balance)
	check(err)

	createdId, _ := result.LastInsertId()
	balance := fmt.Sprintf("%v %v", strconv.FormatFloat(a.Balance, 'f', 2, 32), a.Currency)

	accounts.AddItem(a.Title, balance, 0, func() { WhereAccount(createdId, source) })

	pages.RemovePage("Form")
	source.App.SetFocus(source.AccountList)

	defer db.Close()
}

func SelectAccounts(source *s.Source) ([]string, []s.Account) {
	defer m.ErrorModal(source.Pages, source.Modal)

	db, err := sql.Open("sqlite3", source.Config.Database)
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
	db.Close()
	defer root_accounts.Close()
	return account_titles, account_types
}

func WhereAccount(id int64, source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	str_id := strconv.FormatInt(id, 10)

	request := strings.ReplaceAll(s.TransactionsWhereAccountId, "?", str_id)
	LoadTransactions(request, source)
}
