package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
var (
	bottom_flex = tview.NewFlex().
			SetDirection(tview.FlexColumn)
)
func TransactionsTable(nextSlide func()) (title string, content tview.Primitive) {
	
	//Table
	form := Table()
	
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

	buttons_flex := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(tview.NewButton("Create"), 0, 1, false).
			AddItem(tview.NewButton("Delete"), 0, 1, false)
			
	table_flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(table, 0, 3, true).
			AddItem(buttons_flex, 3, 1, false)
			
	bottom_flex.AddItem(table_flex, 0, 1, true)
	
	modal_flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(top_flex, 20, 1, false).
			AddItem(bottom_flex, 0, 1, true)
				
	flex.AddItem(modal_flex, 0, 2, true)

	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event.Buttons() == tcell.Button1 {
			if form.InRect(event.Position()) == false {
				bottom_flex.RemoveItem(form)
			}
		}
		return event, action
	})

	return "Transactions", flex
}
