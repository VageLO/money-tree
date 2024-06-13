package cmd

import (
	"strings"
	"strconv"

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
			categories, _, _ := SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
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
		form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameNode(text, node) })
	}
	if list != nil {
		title, second := list.GetItemText(list.GetCurrentItem())
		split := strings.Split(second, " ")
		balance := split[0]
		currency := split[1]
		
		form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameAccount(text, "title", list) })
		
		form.AddInputField("Currency: ", currency, 0, nil, func(text string) { RenameAccount(text, "currency", list) })
		
		form.AddInputField("Balance: ", balance, 0, nil, func(text string) { RenameAccount(text, "balance", list) })
	}
}

func FormRenameAccount() {
	FillTreeAndListForm(nil, accounts)
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func FormRenameNode() {
	node := tree.GetCurrentNode()
	FillTreeAndListForm(node, nil)
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func FormAddAccount() {
	form.Clear(true)
	var a account_type
	form.AddInputField("Title: ", "", 0, nil, func(text string) {
		a.title = text
	})
	form.AddInputField("Currency: ", "", 0, nil, func(text string) {
		a.currency = text
	})
	form.AddInputField("Balance: ", "", 0, nil, func(text string) {
		a.balance, _ = strconv.ParseFloat(text, 64); 
	})
	form.AddButton("Add", func() { AddAccount(&a) })
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func FormAddNode() {
	form.Clear(true)
	root := tree.GetRoot()

	n := &node{
		text:   "",
		expand: true,
		reference: &category_type{},
		children: []*node{},
	}
	new_node := add(n, root)

	form.AddInputField("Title: ", "", 0, nil, func(text string) { 
		new_node.SetText(text)
	})
	
	var selected_dropdown *tview.TreeNode
	var options []string
	options = append(options, root.GetText())

	for _, children := range root.GetChildren() {
		options = append(options, children.GetText())
	}

	initial := 0

	selected_node := tree.GetCurrentNode()
	if selected_node != nil {
		for idx, title := range options {
			if title == selected_node.GetText() {
				initial = idx
			}
		}
	}

	form.AddDropDown("Categories", options, initial, func(option string, optionIndex int) {
		if root.GetText() == option {
			selected_dropdown = root
			reference := new_node.GetReference().(*node)
			reference.parent = root
			new_node.SetReference(reference)
			return
		}
		
		for _, children := range root.GetChildren() {
			if children.GetText() == option {
				selected_dropdown = children
				reference := new_node.GetReference().(*node)
				reference.parent = children
				new_node.SetReference(reference)
				return
			}
		}
	})

	form.AddButton("Add", func() {
		AddCategory(new_node, selected_dropdown)
		pages.RemovePage("Dialog")
	})

	pages.AddPage("Dialog", Dialog(form), true, true)
}