package action

import (
	"errors"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"

    "golang.org/x/exp/slices"
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
			Reference:  row.Reference,
		}, table)
	}
}

func UpdateRows(row s.Row, table *tview.Table) {

	for i, data := range row.Data {
		UpdateCell(&s.Cell{
			Row:       row.Index,
			Column:    i,
			Text:      data,
			Reference: row.Reference,
		}, table)
	}
}

func AddToTable(source *s.Source) {
	form := source.Form
	pages := source.Pages
	tree := source.CategoryTree

	defer m.ErrorModal(pages, source.Modal)

	newRow := source.Table.GetRowCount()
	if tree.GetRowCount() <= 1 || source.AccountList.GetItemCount() <= 1 {
		check(errors.New("Account and category must be created"))
	}
	EmptyForm(newRow, source)
	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)

	source.App.SetFocus(form)
}

func SelectMultipleTransactions(row int, source *s.Source) {

	defer m.ErrorModal(source.Pages, source.Modal)
    table := source.Table
    if table.GetRowCount() <= 1 {
	    return
    }

    if row == -1 {
        row, _ = table.GetSelection()
    }
    isTrue := false

    for column := 0; column <= len(source.Columns); column++ {
        cell := table.GetCell(row, column)
        fg, _, _ := cell.Style.Decompose()
        if fg == tcell.ColorRed {
            cell.SetTextColor(tcell.ColorWhite)
        } else if fg == tcell.ColorWhite {
            cell.SetTextColor(tcell.ColorRed)
            isTrue = true
        }
    }
    if isTrue {
        SelectedRows = append(SelectedRows, row)
        slices.Sort(SelectedRows)
        return
    } 
    if value := slices.Index(SelectedRows, row); value != -1 {
        SelectedRows = slices.Delete(SelectedRows, value, value + 1) 
        slices.Sort(SelectedRows)
    }
}
