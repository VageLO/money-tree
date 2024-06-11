package cmd

import (
	//"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type add_transaction struct{
	transaction_type, date, amount, balance, account, category string
}
	
func Form() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Transaction Information")

	return form
}

func FillForm(form *tview.Form, columns int, row int, IsEmptyForm bool) *tview.Form {
	
	form.Clear(true)
	var t add_transaction
	
	changed := func(text string, cell *tview.TableCell) {
		cell.SetText(text)
		UpdateTransaction(cell, text)
	}
	
	added := func(text string, cell *tview.TableCell, label string) {
		cell.SetText(text)
		switch field := label; field {
		case "transaction_type":
			t.transaction_type = text
		case "date":
			t.date = text
		case "amount":
			t.amount = text
		case "balance":
			t.balance = text
		case "account":
			t.account = text
		case "category":
			t.category = text
		}
	}

	form.SetCancelFunc(func() {
		pages.RemovePage("Dialog")
	})

	for i := 0; i < columns; i++ {
		cell := table.GetCell(row, i)
		if IsEmptyForm == false {		
			form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { changed(text, cell) })
			continue
		} 
		InsertRow(&row_settings{
			row: row,
			column: i,
			text: "",
			selectable: true,
			color: tcell.ColorWhite,
		})
		cell = table.GetCell(row, i)
		
		column_name := table.GetCell(0, i).Text
		
		if column_name == "category" {
			categories, _ := SelectCategories()
			form.AddDropDown(table.GetCell(0, i).Text, categories, 0, nil)
			continue
		}
		if column_name == "account" {
			accounts, _ := SelectAccounts()
			form.AddDropDown(table.GetCell(0, i).Text, accounts, 0, nil)
			continue
		}
		
		form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { added(text, cell, column_name) })
	}

	if IsEmptyForm {
		form.AddButton("Add", func() {AddTransaction(&t)})
	}
	
	return form
}

func FillTreeAndListForm(node *tview.TreeNode, list *tview.List) {
	form.Clear(true)

	if node != nil {
		title := node.GetText()
		changed := func(text string, node *tview.TreeNode) {
			node.SetText(text)
		}
		form.AddInputField("Title: ", title, 0, nil, func(text string) { changed(text, node) })
	}
	if list != nil {
		title, _ := list.GetItemText(list.GetCurrentItem())
		changed := func(text string, list *tview.List) {
			list.SetItemText(list.GetCurrentItem(), text, "")
		}
		form.AddInputField("Title: ", title, 0, nil, func(text string) { changed(text, list) })
	}
}
