package cmd

import (
	//"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Form() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Transaction Information")

	return form
}

func FillForm(form *tview.Form, count int, row int, empty bool) *tview.Form {

	changed := func(text string, cell *tview.TableCell) {
		cell.SetText(text)
		UpdateTransaction(cell, text)
	}

	added := func(text string, cell *tview.TableCell, label string) {
		cell.SetText(text)
	}
	
	form.Clear(true)

	form.SetCancelFunc(func() {
		pages.RemovePage("Dialog")
	})

	for i := 0; i < count; i++ {
		cell := table.GetCell(row, i)
		if empty {		
			InsertRow(&row_settings{
				row: row,
				column: i,
				text: "",
				selectable: true,
				color: tcell.ColorWhite,
			})
			cell = table.GetCell(row, i)
			
			form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { added(text, cell, table.GetCell(0, i).Text) })
		} else {
			form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { changed(text, cell) })
		}
	}

	if empty {
		form.AddButton("Add", nil)
	}
	
	return form
}

func FillTreeAndListForm(node *tview.TreeNode, list *tview.List) {
	form.Clear(true)

	if node != nil {
		title := node.GetText()
		changed := func(text string, node *tview.TreeNode) {
			node.SetText(text)
		}
		form.AddInputField("Title: ", title, 0, nil, func(text string) { changed(text, node) })
	}
	if list != nil {
		title, _ := list.GetItemText(list.GetCurrentItem())
		changed := func(text string, list *tview.List) {
			list.SetItemText(list.GetCurrentItem(), text, "")
		}
		form.AddInputField("Title: ", title, 0, nil, func(text string) { changed(text, list) })
	}

}
