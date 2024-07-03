package forms

import (
	s "main/structs"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

func formStyle(formTitle string) {
	form.SetBorder(true).SetTitle(formTitle)
}

func added(text string, label string, t *s.Transaction) {
	if text == "" {
		return
	}
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
	case "to_amount":
		amount, err := strconv.ParseFloat(text, 64)
		check(err)
		t.to_amount.Scan(amount)
	}
}

func Fill(columns int, row int, IsEmptyForm bool) {
	form := tview.NewForm()
	formStyle("Transaction Information")

	var transaction s.Transaction

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
		form.AddButton("Add", func() { AddTransaction(transaction, row) })
	} else if !IsEmptyForm {
		form.AddButton("Save", func() { UpdateTransaction(transaction, row) })
	}
}

func Empty(index int, t *s.Transaction) {

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

		form.AddDropDown(table.GetCell(0, index).Text, categories, 0, func(option string, optionIndex int) { SelectedCategory(option, optionIndex, c_types, t) })
		return
	}
	if column_name == "account" {
		accounts, a_types := SelectAccounts()

		form.AddDropDown(table.GetCell(0, index).Text, accounts, 0, func(option string, optionIndex int) { SelectedAccount(option, optionIndex, a_types, t) })
		return
	}

	form.AddInputField(table.GetCell(0, index).Text, "", 0, nil, func(text string) { added(text, column_name, t) })
}

func Filled(index, row int, t *s.Transaction) {
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
		SelectedCategory(categories[initial], initial, c_types, t)

		form.AddDropDown(column_name, categories, initial, func(option string, optionIndex int) { SelectedCategory(option, optionIndex, c_types, t) })
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
		SelectedAccount(accounts[initial], initial, a_types, t)

		form.AddDropDown(column_name, accounts, initial, func(option string, optionIndex int) { SelectedAccount(option, optionIndex, a_types, t) })
		return
	}

	added(cell.Text, column_name, t)

	form.AddInputField(table.GetCell(0, index).Text, cell.Text, 0, nil, func(text string) { added(text, column_name, t) })
}

func SelectedTransfer(option string, optionIndex int, a_types []s.Account, t *s.Transaction) {
	selected_a := a_types[optionIndex]
	if selected_a.title != option {
		return
	}
	t.to_account_id.Scan(selected_a.id)
}

