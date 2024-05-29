package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const tableData = `OrderDate|Region|Rep|Item|Units|UnitCost|Total
1/6/2017|East|Jones|Pencil|95|1.99|189.05
1/23/2017|Central|Kivell|Binder|50|19.99|999.50
2/9/2017|Central|Jardine|Pencil|36|4.99|179.64
2/26/2017|Central|Gill|Pen|27|19.99|539.73
3/15/2017|West|Sorvino|Pencil|56|2.99|167.44
4/1/2017|East|Jones|Binder|60|4.99|299.40
4/18/2017|Central|Andrews|Pencil|75|1.99|149.25
5/5/2017|Central|Jardine|Pencil|90|4.99|449.10
5/22/2017|West|Thompson|Pencil|32|1.99|63.68
6/8/2017|East|Jones|Binder|60|8.99|539.40
6/25/2017|Central|Morgan|Pencil|90|4.99|449.10
7/12/2017|East|Howard|Binder|29|1.99|57.71
7/29/2017|East|Parent|Binder|81|19.99|1,619.19
8/15/2017|East|Jones|Pencil|35|4.99|174.65
9/1/2017|Central|Smith|Desk|2|125.00|250.00
9/18/2017|East|Jones|Pen Set|16|15.99|255.84
10/5/2017|Central|Morgan|Binder|28|8.99|251.72
10/22/2017|East|Jones|Pen|64|8.99|575.36
11/8/2017|East|Parent|Pen|15|19.99|299.85
11/25/2017|Central|Kivell|Pen Set|96|4.99|479.04
12/12/2017|Central|Smith|Pencil|67|1.29|86.43
12/29/2017|East|Parent|Pen Set|74|15.99|1,183.26
1/15/2018|Central|Gill|Binder|46|8.99|413.54
2/1/2018|Central|Smith|Binder|87|15.00|1,305.00
2/18/2018|East|Jones|Binder|4|4.99|19.96
3/7/2018|West|Sorvino|Binder|7|19.99|139.93
3/24/2018|Central|Jardine|Pen Set|50|4.99|249.50
4/10/2018|Central|Andrews|Pencil|66|1.99|131.34
4/27/2018|East|Howard|Pen|96|4.99|479.04
5/14/2018|Central|Gill|Pencil|53|1.29|68.37
5/31/2018|Central|Gill|Binder|80|8.99|719.20
6/17/2018|Central|Kivell|Desk|5|125.00|625.00
7/4/2018|East|Jones|Pen Set|62|4.99|309.38
7/21/2018|Central|Morgan|Pen Set|55|12.49|686.95
8/7/2018|Central|Kivell|Pen Set|42|23.95|1,005.90
8/24/2018|West|Sorvino|Desk|3|275.00|825.00
9/10/2018|Central|Gill|Pencil|7|1.29|9.03
9/27/2018|West|Sorvino|Pen|76|1.99|151.24
10/14/2018|West|Thompson|Binder|57|19.99|1,139.43
10/31/2018|Central|Andrews|Pencil|14|1.29|18.06
11/17/2018|Central|Jardine|Binder|11|4.99|54.89
12/4/2018|Central|Jardine|Binder|94|19.99|1,879.06
12/21/2018|Central|Andrews|Binder|28|4.99|139.72`

func TransactionsTable(nextSlide func()) (title string, content tview.Primitive) {
	
	count := 0
	table := tview.NewTable().
		SetFixed(1, 1)
	for row, line := range strings.Split(tableData, "\n") {
		count = len(strings.Split(line, "|"))
		for column, cell := range strings.Split(line, "|") {
			color := tcell.ColorWhite
			if row == 0 {
				color = tcell.ColorYellow
			} else if column == 0 {
				color = tcell.ColorDarkCyan
			}
			align := tview.AlignLeft
			tableCell := tview.NewTableCell(cell).
				SetTextColor(color).
				SetAlign(align).
				SetSelectable(row != 0 && column != 0)
			if column >= 1 && column <= 3 {
				tableCell.SetExpansion(1)
			}
			table.SetCell(row, column, tableCell)
		}
	}
	table.SetBorder(true).SetTitle("Table")
	
	// Transaction Form
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Transaction Information")
	
	// List with accounts
	accounts := tview.NewList()
	accounts.ShowSecondaryText(false).
		AddItem("Alfa Bank", "123", '1', nil).
		AddItem("BNB", "123", '2', nil)
	accounts.SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Account List")
	
	// List with categories
	categories := tview.NewList()
	categories.ShowSecondaryText(false).
		AddItem("Work", "123", '1', nil).
		AddItem("Store", "123", '2', nil)
	categories.SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle("Category List")
		
	// Flex:
	flex := tview.NewFlex()
	
	top_flex := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(accounts, 0, 1, false).
			AddItem(categories, 0, 2, false)	
			
	bottom_flex := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(table, 0, 1, true)
			
	modal_flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(top_flex, 20, 1, false).
			AddItem(bottom_flex, 0, 1, true)
				
	flex.AddItem(modal_flex, 0, 2, true)
	
	// Form Buttons
	close := func() {
		bottom_flex.RemoveItem(form)
		app.SetFocus(table)
	}
	save := func() {
		bottom_flex.RemoveItem(form)
		app.SetFocus(table)
	}
	
	app.SetMouseCapture(func(event *tcell.EventMouse, action tview.MouseAction) (*tcell.EventMouse, tview.MouseAction) {
		if event.Buttons() == tcell.Button1 {
			if form.InRect(event.Position()) == false {
				bottom_flex.RemoveItem(form)
			}
		}
		return event, action
	})

	// Table action
	table.Select(0, 0).SetFixed(1, 1).SetSelectedFunc(func(row int, column int) {
		form.Clear(true)
		
		for i := 0; i < count; i++ {
			cell := table.GetCell(row, i)
			form.AddInputField(table.GetCell(0, i).Text, cell.Text, 0, nil, nil)
		}
		form.AddButton("Save", save).AddButton("Cancel", close)
		
		bottom_flex.AddItem(form, 0, 1, false)
		app.SetFocus(form)
	})
		
	selectRow := func() {
		table.SetBorders(false).
			SetSelectable(true, false).
			SetSeparator('|')
	}
	
	selectRow()

	return "Transactions", flex
}
