package cmd

import (
	"strings"
	"strconv"

	//"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func FormStyle(formTitle string) {
	form.SetBorder(true).SetTitle(formTitle)
}

func added(text string, label string, t *Transaction) {
	if text == "" {return}
	defer ErrorModal()

	switch field := label; field {
	case "date":
		t.date = text
	case "description":
		t.description = text
	case "amount":
		amount, err := strconv.ParseFloat(text, 64)
		check(err)
		t.amount = amount
	}
}
	
func FillForm(columns int, row int, IsEmptyForm bool) {
	form.Clear(true)
	FormStyle("Transaction Information")

	var transaction Transaction
	
	form.SetCancelFunc(func() {
		pages.RemovePage("Modal")
	})

	for i := 0; i < columns; i++ {
		if IsEmptyForm == false {	
			filledForm(i, row, &transaction)
			continue
		}
		emptyForm(i, &transaction)
	}
	if IsEmptyForm {
		form.AddButton("Add", func() {AddTransaction(transaction, row)})
	} else if !IsEmptyForm {
		form.AddButton("Save", func() {UpdateTransaction(transaction, row)})
	}
}

func emptyForm(index int, t *Transaction) {

	column_name := table.GetCell(0, index).Text
	
	if column_name == "transaction_type" {
		if exist := IsTransfer(nil, t); !exist {
			TranTypes("debit", t)
			return
		}
		ToAccount(nil, t)
		return
	}
	if column_name == "category" {
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`)
		
		form.AddDropDown(table.GetCell(0, index).Text, categories, 0, func(option string, optionIndex int) { C_Selected(option, optionIndex, c_types, t) })
		return
	}
	if column_name == "account" {
		accounts, a_types := SelectAccounts()

		form.AddDropDown(table.GetCell(0, index).Text, accounts, 0, func(option string, optionIndex int) { A_Selected(option, optionIndex, a_types, t) })
		return
	}
	
	form.AddInputField(table.GetCell(0, index).Text, "", 0, nil, func(text string) { added(text, column_name, t) })
}

func filledForm(index, row int, t *Transaction) {
	column_name := table.GetCell(0, index).Text
	cell := table.GetCell(row, index)

	if column_name == "transaction_type" {
		if exist := IsTransfer(cell, t); !exist {
			TranTypes(cell.Text, t)
			return
		}
		ToAccount(cell, t)
		return
	}
	if column_name == "category" {
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`)
		initial := 0

		for idx, title := range categories {
			if title == cell.Text {
				initial = idx
			}
		}
		C_Selected(categories[initial], initial, c_types, t)
		
		form.AddDropDown(column_name, categories, initial, func(option string, optionIndex int) { C_Selected(option, optionIndex, c_types, t) })
		return
	}
	if column_name == "account" {
		accounts, a_types := SelectAccounts()
		initial := 0
		
		for idx, title := range accounts {
			if title == cell.Text {
				initial = idx
			}
		}
		A_Selected(accounts[initial], initial, a_types, t)
		
		form.AddDropDown(column_name, accounts, initial, func(option string, optionIndex int) { A_Selected(option, optionIndex, a_types, t) })
		return
	}
	
	added(cell.Text, column_name, t)
	
	form.AddInputField(table.GetCell(0, index).Text, cell.Text, 0, nil, func(text string) { added(text, column_name, t) })
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
	if accounts.GetItemCount() <= 0 {
		return
	}
	FillTreeAndListForm(nil, accounts)
	pages.AddPage("Form", Modal(form, 20, 50), true, true)
}

func FormRenameNode() {
	node := tree.GetCurrentNode()
	if node == nil {
		return
	}
	FillTreeAndListForm(node, nil)
	pages.AddPage("Form", Modal(form, 20, 50), true, true)
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
	pages.AddPage("Form", Modal(form, 20, 50), true, true)
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
		pages.RemovePage("Form")
	})
	pages.AddPage("Form", Modal(form, 20, 50), true, true)
}

func A_Selected(option string, optionIndex int, a_types []account_type, t *Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.title != option {
		return
	}
	t.account_id = selected_a.id
	t.account = selected_a.title
}

func C_Selected(option string, optionIndex int, c_types []category_type, t *Transaction) {
	selected_c := c_types[optionIndex]
	if selected_c.title != option {
		return
	}
	t.category_id = selected_c.id
	t.category = selected_c.title
}

func T_Selected(option string, optionIndex int, a_types []account_type, t *Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.title != option {
		return
	}
	t.to_account_id.Scan(selected_a.id)
}