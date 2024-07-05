package action

import (
	"errors"
	m "main/modal"
	s "main/structs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InsertCell(c s.Cell, table *tview.Table) {
	align := tview.AlignLeft
	tableCell := tview.NewTableCell(c.Text).
		SetTextColor(c.Color).
		SetAlign(align).
		SetSelectable(c.Selectable)
	tableCell.SetExpansion(1)
	if c.Reference != nil {
		tableCell.SetReference(c.Reference)
	}
	table.SetCell(c.Row, c.Column, tableCell)
}

func UpdateCell(c *s.Cell, table *tview.Table) {
	tableCell := table.GetCell(c.Row, c.Column)
	tableCell.SetText(c.Text)
	if c.Reference != nil {
		tableCell.SetReference(c.Reference)
	}
}

func InsertRows(row s.Row, table *tview.Table) {

	for i, data := range row.Data {
		InsertCell(s.Cell{
			Row:        row.Index,
			Column:     i,
			Text:       data,
			Selectable: true,
			Color:      tcell.ColorWhite,
			Reference:  row.Transaction,
		}, table)
	}
}

func UpdateRows(row s.Row, table *tview.Table) {

	for i, data := range row.Data {
		UpdateCell(&s.Cell{
			Row:       row.Index,
			Column:    i,
			Text:      data,
			Reference: row.Transaction,
		}, table)
	}
}

func AddToTable(source *s.Source) {
	form := source.Form
	pages := source.Pages
	tree := source.CategoryTree
	
	defer m.ErrorModal(pages, source.Modal)

	newRow := source.Table.GetRowCount()
	if tree.GetRowCount() <= 0 || source.AccountList.GetItemCount() <= 0 {
		check(errors.New("Account and category must be created"))
	}
	Fill(len(source.Columns), newRow, true, source)
	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)

	source.App.SetFocus(form)
}
