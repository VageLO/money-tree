package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

func Form() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Transaction Information")
	
	return form
}

func FillForm(form *tview.Form, count int, row int) *tview.Form {
	// Close Button
	close := func() {
		bottom_flex.RemoveItem(form)
		app.SetFocus(table)
	}
	
	// TODO: Description
	changed := func(text string, row int, column int) {
		color := tcell.ColorWhite
		if row == 0 {
			color = tcell.ColorYellow
		} else if column == 0 {
			color = tcell.ColorDarkCyan
		}
		tableCell := tview.NewTableCell(text).
			SetTextColor(color).
			SetAlign(tview.AlignLeft).
			SetSelectable(row != 0)
		if column >= 1 && column <= 3 {
			tableCell.SetExpansion(1)
		}
		table.SetCell(row, column, tableCell)
	}
	
	form.Clear(true)
	for i := 0; i < count; i++ {
		cell := table.GetCell(row, i)
		column := i
		form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) {changed(text, row, column)})
	}

	// Form Buttons
	form.AddButton("Cancel", close)
	form.AddButton("Delete", func(){
		table.RemoveRow(row)
		bottom_flex.RemoveItem(form)
	})
	
	return form
}
