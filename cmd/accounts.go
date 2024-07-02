package cmd

import (
	"fmt"
	"strconv"
	"main/action"
	
	"github.com/rivo/tview"
)

var (
	accounts = tview.NewList()
)

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
		accounts.AddItem(a.title, second_title, 0, func() { action.WhereAccount(account_id) })
	}
	
	return accounts
}

