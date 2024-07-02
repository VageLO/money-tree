package cmd

import (
	"os"
	m "main/modal"
	"main/action"
	"main/forms"
	
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var transaction_query string

func Transactions() tview.Primitive {
	defer m.ErrorModal(pages, modal)
	
	query, err := os.ReadFile("./sql/Select_On_Transactions.sql")
	transaction_query = string(query)
    check(err)
	
	action.FillTable(transaction_query)
	
	// List with accounts
	accounts := AccountsList()

	// Tree with categories
	categories := TreeView()

	//Flex
	flex := tview.NewFlex()

	top_flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(accounts, 0, 1, false).
		AddItem(categories, 0, 1, false)

	bottom_flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(table, 0, 2, true)

	modal_flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(top_flex, 20, 1, false).
		AddItem(bottom_flex, 0, 1, true)

	flex.AddItem(modal_flex, 0, 2, true)

	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event.Buttons() == tcell.Button1 {
			if form.InRect(event.Position()) == false {
				pages.RemovePage("Form")
			}
			if modal.InRect(event.Position()) == false {
				pages.RemovePage("Modal")
			}
			if file_table.InRect(event.Position()) == false {
				pages.RemovePage("Files")
			}
		}
		return event, action
	})

	return flex
}

func FillTable(request string) {
	table.Clear()
	
	table.SetTitle("Transactions")
	
	action.SelectTransactions(request)
	
	table.Select(1, 1).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		forms.Fill(column_count, row, false)

		pages.AddPage("Form", Modal(form, 30, 50), true, true)
	})

	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator('|')
}