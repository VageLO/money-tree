package cmd

import (
	"strings"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func FormStyle(formTitle string) {
	form.SetBorder(true).SetTitle(formTitle)
}

func FillForm( columns int, row int, IsEmptyForm bool) {
	form.Clear(true)
	FormStyle("Transaction Information")

	var t Transaction
	
	changed := func(text string, cell *tview.TableCell) {
		cell.SetText(text)
		UpdateTransaction(cell, text)
	}
	
	added := func(text string, cell *tview.TableCell, label string) {
		if text == "" {return}
		defer ErrorModal()
		
		cell.SetText(text)
		switch field := label; field {
		case "transaction_type":
			t.transaction_type = text
		case "date":
			t.date = text
		case "amount":
			amount, err := strconv.ParseFloat(text, 64)
			check(err)
			t.amount = amount
		case "balance":
			balance, err := strconv.ParseFloat(text, 64)
			check(err)
			t.balance = balance
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
		InsertCell(&cell_type{
			row: row,
			column: i,
			text: "",
			selectable: true,
			color: tcell.ColorWhite,
		})
		cell = table.GetCell(row, i)
		
		column_name := table.GetCell(0, i).Text
		
		if column_name == "category" {
			categories, c_types, _ := SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
			form.AddDropDown(table.GetCell(0, i).Text, categories, 0, func(option string, optionIndex int) { C_Selected(option, optionIndex, c_types, &t) })
			continue
		}
		if column_name == "account" {
			accounts, a_types := SelectAccounts()
			form.AddDropDown(table.GetCell(0, i).Text, accounts, 0, func(option string, optionIndex int) { A_Selected(option, optionIndex, a_types, &t) })
			continue
		}
		
		form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { added(text, cell, column_name) })
	}

	if IsEmptyForm {
		form.AddButton("Add", func() {AddTransaction(&t)})
	}
}

func FillTreeAndListForm(node *tview.TreeNode, list *tview.List) {
	form.Clear(true)

	if node != nil {
		FormStyle("Category Information")
		title := node.GetText()
		form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameNode(text, node) })
	}
	if list != nil {
		FormStyle("Account Information")
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
	FormStyle("Add Account")
	
	var a = account_type{}
	
	form.AddInputField("Title: ", "", 0, nil, func(text string) {
		a.title = text
	})
	form.AddInputField("Currency: ", "", 0, nil, func(text string) {
		a.currency = text
	})
	form.AddInputField("Balance: ", "0", 0, nil, func(text string) {
		if text == "" { 
			a.balance = 0
			return
		}
		defer ErrorModal()

		balance, err := strconv.ParseFloat(text, 64);
		check(err)
		a.balance = balance
	})
	form.AddButton("Add", func() { AddAccount(&a) })
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func FormAddNode() {
	form.Clear(true)
	FormStyle("Add Category")
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

func A_Selected(option string, optionIndex int, a_types []account_type, t *Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.title != option {
		return
	}
	t.account_id = selected_a.id
}

func C_Selected(option string, optionIndex int, c_types []category_type, t *Transaction) {
	selected_c := c_types[optionIndex]
	if selected_c.title != option {
		return
	}
	t.category_id = selected_c.id
}