package action

import (
	"errors"
	//"fmt"
	m "main/modal"
	s "main/structs"
	"strconv"
	"time"

	"github.com/rivo/tview"
)

func FormStyle(formTitle string, form *tview.Form) {
	form.SetBorder(true).SetTitle(formTitle)
}

func FillForm(columnsLen int, row int, IsEmptyForm bool, source *s.Source) {

	form := source.Form
	form.Clear(true)

	FormStyle("Transaction Details", form)

	var transaction s.Transaction

	for i := 0; i < columnsLen; i++ {
		if !IsEmptyForm {
			table := source.Table
			cell := table.GetCell(row, i)
			transaction = cell.GetReference().(s.Transaction)

			FilledForm(i, row, &transaction, source)
			continue
		}
		EmptyForm(i, &transaction, source)
	}

	initTransaction := transaction

	if IsEmptyForm {
		form.AddButton("➕", func() { AddTransaction(transaction, row, source) })

	} else if !IsEmptyForm {
		form.AddButton("💾", func() {
			defer m.ErrorModal(source.Pages, source.Modal)
			if initTransaction == transaction {
				check(errors.New("Nothing Changed"))
			}
			UpdateTransaction(transaction, row, source)
		})
	}

	form.AddButton("Add Attacments", func() {
		defer m.ErrorModal(source.Pages, source.Modal)

		source.Attachments = []string{}
		pageName := "FileExporer"
		source.Pages.AddPage(
			pageName,
			FileExporer(source, "", pageName), true, true,
		)
	})

	// TODO: FileExplorer
	form.AddButton("📎", func() {
		attachments := findAttachments(source, transaction.Id)
		m.FileTable(source, "Attachments", attachments, m.OpenFiles)
		//source.Pages.AddPage("FileExplorer", m.NewTree(source), true, true)
	})

	form.AddButton("❌", func() {
		source.Pages.RemovePage("Form")
	})
}

func EmptyForm(index int, t *s.Transaction, source *s.Source) {

	table := source.Table
	columns := source.Columns
	form := source.Form
	accountsList := source.AccountList
	tree := source.CategoryTree

	columnName := table.GetCell(0, index).Text

	// Transaction Type field
	if columnName == columns[5] {
		TransactionTypes(t, source)
		return
	}

	// Category field
	if columnName == columns[3] {
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

		form.AddDropDown(
			table.GetCell(0, index).Text,
			categories,
			initial,
			func(option string, optionIndex int) {
				SelectedCategory(option, optionIndex, c_types, t)
			})
		return
	}

	// Account field
	if columnName == columns[2] {
		accounts, a_types := SelectAccounts(source)

		initial := 0
		title, _ := accountsList.GetItemText(accountsList.GetCurrentItem())
		for idx, account := range accounts {
			if account == title {
				initial = idx
			}
		}

		form.AddDropDown(
			table.GetCell(0, index).Text,
			accounts,
			initial,
			func(option string, optionIndex int) {
				SelectedAccount(option, optionIndex, a_types, t)
			})
		return
	}

	// Date field
	if columnName == columns[1] {
		date := time.Now().Format("2006-01-02")
		added(date, columnName, t, source)

		form.AddInputField(table.GetCell(0, index).Text, date, 0, nil, func(text string) {
			added(text, columnName, t, source)
		})
		return
	}

	form.AddInputField(table.GetCell(0, index).Text, "", 0, nil, func(text string) {
		added(text, columnName, t, source)
	})
}

func FilledForm(index, row int, t *s.Transaction, source *s.Source) {

	defer m.ErrorModal(source.Pages, source.Modal)
	table := source.Table
	columns := source.Columns
	form := source.Form

	columnName := table.GetCell(0, index).Text
	cell := table.GetCell(row, index)

	// Transaction Type field
	if columnName == columns[5] {
		TransactionTypes(t, source)
		return
	}

	// Category field
	if columnName == columns[3] {
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`, source)
		initial := 0

		for idx, title := range categories {
			if title == cell.Text {
				initial = idx
			}
		}
		SelectedCategory(categories[initial], initial, c_types, t)

		form.AddDropDown(columnName, categories, initial, func(option string, optionIndex int) {
			SelectedCategory(option, optionIndex, c_types, t)
		})
		return
	}

	// Account field
	if columnName == columns[2] {
		accounts, a_types := SelectAccounts(source)
		initial := 0

		for idx, title := range accounts {
			if title == cell.Text {
				initial = idx
			}
		}
		SelectedAccount(accounts[initial], initial, a_types, t)

		form.AddDropDown(
			columnName,
			accounts,
			initial,
			func(option string, optionIndex int) {
				SelectedAccount(option, optionIndex, a_types, t)
			})
		return
	}

	added(cell.Text, columnName, t, source)

	form.AddInputField(
		columnName,
		cell.Text,
		0,
		nil,
		func(text string) {
			added(text, columnName, t, source)
		})
}

func SelectedTransfer(option string, optionIndex int, a_types []s.Account, t *s.Transaction) {

	selected_a := a_types[optionIndex]
	if selected_a.Title != option {
		return
	}
	t.ToAccountId.Scan(selected_a.Id)
	t.ToAccount.Scan(selected_a.Title)

}

func TransactionTypes(t *s.Transaction, source *s.Source) {

	form := source.Form
	types := []string{"debit", "credit", "transfer"}
	initial := 0

	for idx, title := range types {
		if title == t.TransactionType {
			initial = idx
		}
	}

	form.AddDropDown(
		"Transaction Type",
		types,
		initial,
		func(option string, optionIndex int) {
			ToAccountIndex := form.GetFormItemIndex("To Account")
			ToAmountIndex := form.GetFormItemIndex("To Amount")

			if option == "transfer" && ToAccountIndex == -1 && ToAmountIndex == -1 {
				Transfer(source, t)
				t.TransactionType = "transfer"
				return
			}

			if option != "transfer" && ToAccountIndex != -1 && ToAmountIndex != -1 {
				form.RemoveFormItem(ToAmountIndex)
				form.RemoveFormItem(ToAccountIndex)
				t.ToAccountId.Valid = false
				t.ToAccountId.Scan(nil)
				t.ToAccount.Scan(nil)
				t.ToAmount.Valid = false
				t.ToAmount.Scan(nil)
			}
			t.TransactionType = option
		},
	)
}

func Transfer(source *s.Source, t *s.Transaction) {

	form := source.Form
	initial := 0
	var label string
	var amount string

	label = t.ToAccount.String
	amount = strconv.FormatFloat(t.ToAmount.Float64, 'f', 2, 32)

	accounts, a_types := SelectAccounts(source)

	for idx, title := range accounts {
		if title == label {
			initial = idx
		}
	}

	SelectedTransfer(accounts[initial], initial, a_types, t)

	form.AddDropDown(
		"To Account",
		accounts,
		initial,
		func(option string, optionIndex int) {
			SelectedTransfer(option, optionIndex, a_types, t)
		})

	form.AddInputField(
		"To Amount",
		amount,
		0,
		nil,
		func(text string) {
			added(text, "To Amount", t, source)
		})
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
	case "To Amount":
		amount, err := strconv.ParseFloat(text, 64)
		check(err)
		t.ToAmount.Scan(amount)
	}
}
