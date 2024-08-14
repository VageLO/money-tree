package cmd

import (
	"database/sql"
	"github.com/VageLO/money-tree/action"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func Statistics(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)

	var t s.Transaction
	var stats []s.Statistics

	table := tview.NewTable()
	table.
		SetBorders(true).
		SetBorder(true).
		SetTitle("Statistics")

    currentYear, currentMonth, _:= time.Now().Date()

	startDate := tview.NewInputField().SetLabel(" Start Date: ")
	endDate := tview.NewInputField().SetLabel(" End Date: ")
	yearField := tview.NewInputField().SetLabel(" Year: ").SetText(strconv.Itoa(currentYear))

    months := []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July,
        time.August, time.September, time.October, time.November, time.December}

    var strMonths []string
    for _, month := range months {
        strMonths = append(strMonths, month.String())
    }

	dropDown := tview.NewDropDown().SetLabel("Account: ")
	monthsDropDown:= tview.NewDropDown().SetLabel("Months: ")

	// Flex
	topFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(dropDown, 0, 1, false).
		AddItem(yearField, 0, 1, false).
		AddItem(monthsDropDown, 0, 1, false).
		AddItem(startDate, 0, 1, false).
		AddItem(endDate, 0, 1, false)
	topFlex.SetBorder(true)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, false).
		AddItem(table, 0, 10, true)

	accounts, a_types := action.SelectAccounts(source)

	reloadChart := func(text string, index int) {
		action.SelectedAccount(text, index, a_types, &t)
		stats = getStatistics(source, &t, startDate.GetText(), endDate.GetText())
		loadStatictisTable(table, stats)
	}

    // StartDate Field
	startDate.SetDoneFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
	startDate.SetFinishedFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
    // EndDate Field
	endDate.SetDoneFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
	endDate.SetFinishedFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
    // Year Field
    inputYear := func() {
	    defer m.ErrorModal(source.Pages, source.Modal)

        index, _ := monthsDropDown.GetCurrentOption()
        accountIndex, accountText := dropDown.GetCurrentOption()

        //TODO: Make year validation
        year, err := strconv.Atoi(yearField.GetText())
        check(err)

        firstDay, lastDay := getMonth(year, months[index])
        startDate.SetText(firstDay)
        endDate.SetText(lastDay)

		reloadChart(accountText, accountIndex)
    }

    yearField.SetDoneFunc(func(key tcell.Key) {
        inputYear()
	})
	yearField.SetFinishedFunc(func(key tcell.Key) {
        inputYear()
	})

    
    monthsDropDown.SetOptions(strMonths, func(text string, index int) {
        accountIndex, accountText := dropDown.GetCurrentOption()

        firstDay, lastDay := getMonth(currentYear, months[index])
        startDate.SetText(firstDay)
        endDate.SetText(lastDay)

		reloadChart(accountText, accountIndex)
	})
    
	dropDown.SetOptions(accounts, func(text string, index int) {
		reloadChart(text, index)
	})
	dropDown.SetCurrentOption(0)
    
    // Set current month in drop down
    intMonth := int(currentMonth)
	monthsDropDown.SetCurrentOption(intMonth - 1)

	source.Pages.AddPage("Statistics", flex, true, true)
}

func loadStatictisTable(table *tview.Table, data []s.Statistics) {
	table.Clear()

	columns := []string{"Category", "Debit", "Credit", "Total"}
	for i, columnTitle := range columns {
		action.InsertCell(s.Cell{
			Row:        0,
			Column:     i,
			Text:       columnTitle,
			Selectable: false,
			Color:      tcell.ColorYellow,
		}, table)
	}

	for i, value := range data {
		action.InsertRows(s.Row{
			Columns: columns,
			Index:   i + 1,
			Data: []string{
				value.Category,
				strconv.FormatFloat(value.Debit, 'f', 2, 32),
				strconv.FormatFloat(value.Credit, 'f', 2, 32),
				strconv.FormatFloat(value.Credit-value.Debit, 'f', 2, 32),
			},
			Reference: value,
		}, table)
	}
}

func getStatistics(source *s.Source, t *s.Transaction, firstDate, lastDate string) []s.Statistics {
	defer m.ErrorModal(source.Pages, source.Modal)

	db, err := sql.Open("sqlite3", "./database.db")
	check(err)

	str_id := strconv.FormatInt(t.AccountId, 10)
	request := s.StatisticsQuery
	request = strings.ReplaceAll(request, "?", str_id)
	request = strings.ReplaceAll(request, "FIRST", firstDate)
	request = strings.ReplaceAll(request, "LAST", lastDate)

	result, err := db.Query(request)
	check(err)

	var stats []s.Statistics

	for result.Next() {
		var _s s.Statistics

		err := result.Scan(&_s.Debit, &_s.Credit, &_s.Category)
		check(err)

		stats = append(stats, _s)
	}
	defer result.Close()
	defer db.Close()

	return stats
}

func getMonth(year int, month time.Month) (firstDay, lastDay string) {
	now := time.Now()

	currentLocation := now.Location()

	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth.Format("2006-01-02"), lastOfMonth.Format("2006-01-02")
}
