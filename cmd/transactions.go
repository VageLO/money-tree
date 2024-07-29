package cmd

import (
	"main/action"
	m "main/modal"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var transaction_query string

func Transactions() tview.Primitive {
	defer m.ErrorModal(source.Pages, source.Modal)

	query, err := os.ReadFile("./sql/Select_On_Transactions.sql")
	transaction_query = string(query)
	check(err)

	action.LoadTransactions(transaction_query, source)
	action.LoadAccounts(source)
	CategoryTree()

	flex := tview.NewFlex()

	topFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(source.AccountList, 0, 1, false).
		AddItem(source.CategoryTree, 0, 1, false)

	bottomFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(source.Table, 0, 2, true)

	modalFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topFlex, 20, 1, false).
		AddItem(bottomFlex, 0, 1, true)

	flex.AddItem(modalFlex, 0, 2, true)

	source.App.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event.Buttons() == tcell.Button1 {
			if source.Form.InRect(event.Position()) == false &&
				!source.Pages.HasPage("FileExplorer") {
				source.Pages.RemovePage("Form")
			}
			if source.Modal.InRect(event.Position()) == false {
				source.Pages.RemovePage("Modal")
			}
			if source.FileTable.InRect(event.Position()) == false &&
				source.Pages.HasPage("Files") {
				source.Pages.RemovePage("Files")
			}
			if source.FileTable.InRect(event.Position()) == false &&
				source.Pages.HasPage("Attachments") {
				source.Pages.RemovePage("Attachments")
			}
		}
		return event, action
	})

	return flex
}
