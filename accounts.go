package main

import (
	"github.com/rivo/tview"
)

var (
	accounts = tview.NewList()
)

func AccountsList() *tview.List {
	tableData := `OrderDate|Region|Rep|Item|Units|UnitCost|Total
1/6/2017|East|Jones|Pencil|95|1.99|189.05`
		
	accounts.ShowSecondaryText(false).
		AddItem("Alfa Bank", "123", '1', nil).
		AddItem("BNB", "123", '2', func() {
			_ = FillTable(tableData)
			table.SetBorder(true).SetTitle("Accounts")
			app.SetFocus(table)
		})
	accounts.
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")
		
	return accounts
}