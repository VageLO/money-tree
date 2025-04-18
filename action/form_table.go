package action

import (
	"errors"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"reflect"
	"strconv"
	"time"

	// TODO: Change to slices go1.22
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

func FormStyle(formTitle string, form *tview.Form) {
	form.SetBorder(true).SetTitle(formTitle)
}

func MultiSelectionForm(source *s.Source) {

	form := source.Form
	form.Clear(true)

	FormStyle("Multi Transaction Details", form)

	columnsLen := len(source.Columns)

	var transaction s.Transaction

	for i := 0; i < columnsLen; i++ {
		emptyFormFields(i, &transaction, source)
		locateCheckboxes(form, i)
	}

	form.AddButton("Save", func() {
		defer m.ErrorModal(source.Pages, source.Modal)

		for _, row := range SelectedRows {
			cell := source.Table.GetCell(row, 0)
			reference := cell.GetReference().(s.Transaction)
			transaction.Id = reference.Id
			t := transaction

			CompTransactions(&t, &reference)

			check(t.IsEmpty())
			UpdateTransaction(t, row, source, false)
		}
		length := len(SelectedRows)
		for i := 0; i < length; i++ {
			SelectMultipleTransactions(SelectedRows[0], source)
		}
		source.Pages.RemovePage("Form")
	})

	// Clear attachments array
	source.Attachments = []string{}

	form.AddButton("📎", func() {
		m.FileTable(source, "Attachments", source.Attachments, m.OpenFiles)
	})

	source.Pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func EmptyForm(row int, source *s.Source) {

	form := source.Form
	form.Clear(true)

	FormStyle("Fill Transaction Details", form)

	columnsLen := len(source.Columns)

	var transaction s.Transaction

	for i := 0; i < columnsLen; i++ {
		emptyFormFields(i, &transaction, source)
	}

	form.AddButton("Add", func() { AddTransaction(transaction, row, source) })

	// Clear attachments array
	source.Attachments = []string{}

	form.AddButton("📎", func() {
		m.FileTable(source, "Attachments", source.Attachments, m.OpenFiles)
	})
}

func emptyFormFields(index int, t *s.Transaction, source *s.Source) {

	table := source.Table
	columns := source.Columns
	form := source.Form
	accountsList := source.AccountList
	tree := source.CategoryTree

	columnName := table.GetCell(0, index).Text

	switch columnName {

	case columns[5]:
		// Transaction Type field
		TransactionTypes(t, source)

	case columns[3]:
		// Category field
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

	case columns[2]:
		// Account field
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

	case columns[1]:
		// Date field
		date := time.Now().Format("2006-01-02")
		added(date, columnName, t, source)

		form.AddInputField(table.GetCell(0, index).Text, date, 0, nil, func(text string) {
			added(text, columnName, t, source)
		})

	default:
		form.AddInputField(table.GetCell(0, index).Text, "", 0, nil, func(text string) {
			added(text, columnName, t, source)
		})
	}
}

func FilledForm(row int, source *s.Source) {

	form := source.Form
	form.Clear(true)

	FormStyle("Transaction Details", form)

	columnsLen := len(source.Columns)

	var transaction s.Transaction

	for i := 0; i < columnsLen; i++ {
		table := source.Table
		cell := table.GetCell(row, i)
		transaction = cell.GetReference().(s.Transaction)

		filledFormFields(i, row, &transaction, source)
	}

	initTransaction := transaction

	// Find attachments by trasaction ID
	source.Attachments = findAttachments(source, transaction.Id)

	initAttachments := source.Attachments

	form.AddButton("Save", func() {
		defer m.ErrorModal(source.Pages, source.Modal)

		if initTransaction != transaction {
			UpdateTransaction(transaction, row, source, true)
		} else if slices.Compare(initAttachments, source.Attachments) != 0 {
			updateAttachments(source, source.Attachments, transaction.Id, true)
			source.Pages.RemovePage("Form")
		} else {
			check(errors.New("Change something"))
		}
	})

	form.AddButton("📎", func() {
		m.FileTable(source, "Attachments", source.Attachments, m.OpenFiles)
	})
}

func filledFormFields(index, row int, t *s.Transaction, source *s.Source) {

	defer m.ErrorModal(source.Pages, source.Modal)
	table := source.Table
	columns := source.Columns
	form := source.Form

	columnName := table.GetCell(0, index).Text
	cell := table.GetCell(row, index)

	switch columnName {

	case columns[5]:
		// Transaction Type field
		TransactionTypes(t, source)

	case columns[3]:
		// Category field
		categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`, source)
		initial := 0

		for idx, title := range categories {
			if title == cell.GetReference().(s.Transaction).Category {
				initial = idx
			}
		}
		SelectedCategory(categories[initial], initial, c_types, t)

		form.AddDropDown(columnName, categories, initial, func(option string, optionIndex int) {
			SelectedCategory(option, optionIndex, c_types, t)
		})

	case columns[2]:
		// Account field
		accounts, a_types := SelectAccounts(source)
		initial := 0

		for idx, title := range accounts {
			if title == cell.GetReference().(s.Transaction).Account {
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

	case columns[4]:
		// Amount field
		amount := strconv.FormatFloat(cell.GetReference().(s.Transaction).Amount, 'f', 2, 32)
		added(amount, columnName, t, source)

		form.AddInputField(
			columnName,
			amount,
			0,
			nil,
			func(text string) {
				added(text, columnName, t, source)
			})

	default:
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
	types := []string{"Withdrawal", "Deposit", "Transfer"}
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

			if option == "Transfer" && ToAccountIndex == -1 && ToAmountIndex == -1 {
				Transfer(source, t)
				t.TransactionType = "Transfer"
				return
			}

			if option != "Transfer" && ToAccountIndex != -1 && ToAmountIndex != -1 {
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
	if text == "" && label != "Description" {
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

func CompTransactions(newT, oldT *s.Transaction) {

	nT := reflect.ValueOf(newT).Elem()
	oT := reflect.ValueOf(oldT).Elem()

	tN := nT.Type()

	numFields := tN.NumField()

	for i := 0; i < numFields; i++ {
		newFieldType := tN.Field(i).Type

		newFieldValue := nT.Field(i)
		oldFieldValue := oT.Field(i)

		switch newFieldType.Kind() {

		case reflect.Int64:
			newInt := newFieldValue.Int()
			oldInt := oldFieldValue.Int()
			if newInt != oldInt {
				if newFieldValue.CanSet() && !newFieldValue.OverflowInt(newInt) {
					newFieldValue.SetInt(newInt)
				}
			}

		case reflect.Float64:
			newFloat := newFieldValue.Float()
			oldFloat := oldFieldValue.Float()
			if newFloat == oldFloat {
				continue
			} else if tN.Field(i).Name == "Amount" && newFloat == 0 {
				newFieldValue.SetFloat(oldFloat)
			}

		case reflect.String:
			newString := newFieldValue.String()
			oldString := oldFieldValue.String()
			if newString == oldString {
				continue
			} else if newString == "" {
				newFieldValue.SetString(oldString)
			} else if tN.Field(i).Name == "Date" {
				newFieldValue.SetString(oldString)
			}
		}
	}
}
