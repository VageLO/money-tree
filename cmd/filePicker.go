package cmd

import (
	"os"
	"strings"
	"strconv"
	"errors"
	"reflect"
	"main/parser"
	
	"github.com/rivo/tview"
)

func FilePicker(path string) {
	defer ErrorModal()
	
	file_table.Clear()
	file_table.SetTitle("Pdf files")
	file_table.SetSelectable(true, false)
	
	folder, err := os.Open(path)
    check(err)
	
    files, err := folder.Readdir(0)
    check(err)
	
	count := 0
    for _, v := range files {
		slice := strings.Split(v.Name(), ".")
		if v.IsDir() || slice[len(slice)-1] != "pdf" {
			continue
		}
        
		tableCell := tview.NewTableCell(v.Name())
		tableCell.SetReference(struct {
			path string
		} {path + "/" + v.Name()})
		tableCell.SetSelectable(true)
		
		file_table.SetCell(count, 0, tableCell)
		count++
    }
	
	file_table.SetSelectedFunc(func(row int, column int) {
		if tree.GetRowCount() <= 0 || accounts.GetItemCount() <= 0 {
			check(errors.New("Account and category must be created"))
		}
		
		cell := file_table.GetCell(row, column)
		
		ref := cell.GetReference().(struct { 
			path string
		})
		
		var transaction Transaction
		selectForm(ref.path, &transaction)
		
		pages.RemovePage("Files")
	})
	x, _, _, _ := file_table.GetRect()
	pages.AddPage("Files", Modal(file_table, 30, x), true, true)
}

func insertIntoDb(path string, t *Transaction) {
	transactions := parser.ParsePdf(path)

	for _, import_transaction := range transactions {
		newRow := table.GetRowCount()
		var transaction Transaction
		
		debit := []string{"Безналичная операция", "Отправление средств"}
		credit := []string{"Поступление"}
		
		transaction.account_id = t.account_id
		transaction.category_id = t.category_id
		transaction.account = t.account
		transaction.category = t.category
		
		r := reflect.ValueOf(&import_transaction).Elem()
		rt := r.Type()
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			rv := reflect.ValueOf(&import_transaction)
			value := reflect.Indirect(rv).FieldByName(field.Name)
			
			if field.Name == "date" {
				transaction.date = value.String()
			} 
			if field.Name == "price" {
				amount, _ := strconv.ParseFloat(value.String(), 64)
				transaction.amount = amount
			} 
			if field.Name == "description" {
				transaction.description = value.String()
			}
			if field.Name == "typeof" && contains(debit, value.String()){
				transaction.transaction_type = "debit"
			} else if field.Name == "typeof" && contains(credit, value.String()) {
				transaction.transaction_type = "credit"
			}
		}
		AddTransaction(transaction, newRow)
	}
}

func selectForm(path string, t *Transaction) {
	form.Clear(true)
	FormStyle("Select account and category")

	categories, c_types, _ := SelectCategories(`SELECT * FROM Categories`)
	accounts, a_types := SelectAccounts()
	
	C_Selected(categories[0], 0, c_types, t)
	A_Selected(accounts[0], 0, a_types, t)

	form.AddDropDown("Category", categories, 0, func(option string, optionIndex int) { C_Selected(option, optionIndex, c_types, t) })
	
	form.AddDropDown("Account", accounts, 0, func(option string, optionIndex int) { A_Selected(option, optionIndex, a_types, t) })
	
	form.AddButton("Import", func() {insertIntoDb(path, t)})
	pages.AddPage("Form", Modal(form, 20, 50), true, true)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}