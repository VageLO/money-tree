package cmd

import (
	"fmt"
	"main/action"
	m "main/modal"
	"strconv"
)

func AccountsList() {
	defer m.ErrorModal(source.Pages, source.Modal)
	source.AccountList.Clear()

	_, account_types, err := action.SelectAccounts()
	check(err)

	source.AccountList.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")

	for _, a := range account_types {
		account_id := a.Id
		second_title := fmt.Sprintf("%v %v", strconv.FormatFloat(a.Balance, 'f', 2, 32), a.Currency)
		source.AccountList.AddItem(a.Title, second_title, 0, func() { action.WhereAccount(account_id) })
	}
}
