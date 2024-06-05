package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Form() *tview.Form {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Transaction Information")

	return form
}

func FillForm(form *tview.Form, count int, row int, empty bool) *tview.Form {

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

	form.SetCancelFunc(func() {
		pages.RemovePage("Dialog")
	})

	for i := 0; i < count; i++ {
		cell := table.GetCell(row, i)
		column := i
		if empty {
			form.AddInputField(table.GetCell(0, i).Text, "", 0, nil, func(text string) { changed(text, row, column) })
		} else {
			form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, func(text string) { changed(text, row, column) })
		}
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
