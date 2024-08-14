package cmd

import (
	"github.com/VageLO/money-tree/action"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var transaction_query string

func Transactions() tview.Primitive {
	defer m.ErrorModal(source.Pages, source.Modal)

	action.LoadTransactions(s.Transactions, source)
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
		AddItem(topFlex, 0, 1, false).
		AddItem(bottomFlex, 0, 2, true)

	flex.AddItem(modalFlex, 0, 2, true)

	source.App.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		defer m.ErrorModal(source.Pages, source.Modal)
		pages := source.Pages
		if event.Buttons() == tcell.Button1 {
			if source.Modal.InRect(event.Position()) == false {
				pages.RemovePage("Modal")
			}
		}
		return event, action
	})

	return flex
}
