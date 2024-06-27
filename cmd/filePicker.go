package cmd

import (
	"os"
	"strings"
	"strconv"
	"errors"
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
		
		transactions := parser.ParsePdf(ref.path)
		insertIntoDb(transactions)
		
		pages.RemovePage("Files")
	})
	x, _, _, _ := file_table.GetRect()
	pages.AddPage("Files", Modal(file_table, 30, x), true, true)
}

func insertIntoDb(transactions []parser.Transaction) {

	
	for _, import_transaction := range transactions {
		newRow := table.GetRowCount()
		
		var transaction Transaction
		
		transaction.transaction_type = "debit"
		transaction.account_id = 2
		transaction.category_id = 2
		transaction.date = import_transaction.date
		amount, _ := strconv.ParseFloat(import_transaction.price, 64)
		transaction.amount = amount
		transaction.description = import_transaction.description
		
		AddTransaction(transaction, newRow)
	}
}