package action

import (
	m "main/modal"
	s "main/structs"
	"strconv"

	"github.com/rivo/tview"
)

func FormStyle(formTitle string, form *tview.Form) {
	form.SetBorder(true).SetTitle(formTitle)
}

func added(text string, label string, t *s.Transaction, source *s.Source) {
	if text == "" {
		return
	}
	defer m.ErrorModal(source.Pages, source.Modal)

	switch field := label; field {
	case "Date":
		t.Date = text
	case "Description":
		t.Description = text
	case "Amount":
		amount, err := strconv.ParseFloat(text, 64)
		check(err)
		t.Amount = amount
	case "to_amount":
		amount, err := strconv.ParseFloat(text, 64)
		check(err)
		t.ToAmount.Scan(amount)
	}
}

func Fill(columnsLen int, row int, IsEmptyForm bool, source *s.Source) {
	form := source.Form
	form.Clear(true)
	FormStyle("Transaction Information", form)

	var transaction s.Transaction

	form.SetCancelFunc(func() {
		source.Pages.RemovePage("Form")
	})

	for i := 0; i < columnsLen; i++ {
		if IsEmptyForm == false {
			Filled(i, row, &transaction, source)
			continue
		}
		Empty(i, &transaction, source)
	}
	if IsEmptyForm {
		form.AddButton("Add", func() { AddTransaction(transaction, row, source) })
	} else if !IsEmptyForm {
		form.AddButton("Save", func() { UpdateTransaction(transaction, row, source) })
	}
}

func Empty(index int, t *s.Transaction, source *s.Source) {
	table := source.Table
	columns := source.Columns
	form := source.Form
	accountsList := source.AccountList
	tree := source.CategoryTree
	column_name := table.GetCell(0, index).Text

	if column_name == columns[5] {
		if exist := IsTransfer(source, nil, t); !exist {
			TranTypes(form, "debit", t)
			return
		}
		ToAccount(source, nil, t)
		return
	}
	if column_name == columns[3] {
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`, source)
		
		initial := 0
		
		selectedNode := tree.GetCurrentNode()
		if selectedNode != nil {
			for idx, title := range categories {
				if title == selectedNode.GetText() {
					initial = idx
				}
			}
		}
		
		form.AddDropDown(table.GetCell(0, index).Text, categories, initial, func(option string, optionIndex int) { SelectedCategory(option, optionIndex, c_types, t) })
		return
	}
	if column_name == columns[2] {
		accounts, a_types := SelectAccounts(source)
		
		form.AddDropDown(table.GetCell(0, index).Text, accounts, accountsList.GetCurrentItem(), func(option string, optionIndex int) { SelectedAccount(option, optionIndex, a_types, t) })
		return
	}

	form.AddInputField(table.GetCell(0, index).Text, "", 0, nil, func(text string) { added(text, column_name, t, source) })
}

func Filled(index, row int, t *s.Transaction, source *s.Source) {
	table := source.Table
	columns := source.Columns
	form := source.Form

	column_name := table.GetCell(0, index).Text
	cell := table.GetCell(row, index)

	if column_name == columns[5] {
		if exist := IsTransfer(source, cell, t); !exist {
			TranTypes(form, cell.Text, t)
			return
		}
		ToAccount(source, cell, t)
		return
	}
	if column_name == columns[3] {
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`, source)
		initial := 0

		for idx, title := range categories {
			if title == cell.Text {
				initial = idx
			}
		}
		SelectedCategory(categories[initial], initial, c_types, t)

		form.AddDropDown(column_name, categories, initial, func(option string, optionIndex int) { SelectedCategory(option, optionIndex, c_types, t) })
		return
	}
	if column_name == columns[2] {
		accounts, a_types := SelectAccounts(source)
		initial := 0

		for idx, title := range accounts {
			if title == cell.Text {
				initial = idx
			}
		}
		SelectedAccount(accounts[initial], initial, a_types, t)

		form.AddDropDown(column_name, accounts, initial, func(option string, optionIndex int) { SelectedAccount(option, optionIndex, a_types, t) })
		return
	}

	added(cell.Text, column_name, t, source)

	form.AddInputField(table.GetCell(0, index).Text, cell.Text, 0, nil, func(text string) { added(text, column_name, t, source) })
}

func SelectedTransfer(option string, optionIndex int, a_types []s.Account, t *s.Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.Title != option {
		return
	}
	t.ToAccountId.Scan(selected_a.Id)
}
