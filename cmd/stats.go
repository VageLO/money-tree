package cmd

import (
	"database/sql"
	"main/action"
	m "main/modal"
	s "main/structs"
	"os"
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

	firstDay, lastDay := getCurrentMonth()
	startDate := tview.NewInputField().SetLabel(" Start Date: ").SetText(firstDay)
	endDate := tview.NewInputField().SetLabel(" End Date: ").SetText(lastDay)

	dropDown := tview.NewDropDown().SetLabel("Account: ")

	// Flex
	topFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(dropDown, 0, 1, false).
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

	// InputFields exit functions
	startDate.SetDoneFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
	startDate.SetFinishedFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
	endDate.SetDoneFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})
	endDate.SetFinishedFunc(func(key tcell.Key) {
		index, text := dropDown.GetCurrentOption()
		reloadChart(text, index)
	})

	dropDown.SetOptions(accounts, func(text string, index int) {
		reloadChart(text, index)
	})
	dropDown.SetCurrentOption(0)

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

	query, err := os.ReadFile("./sql/Select_Sum_Of_Account.sql")
	check(err)

	str_id := strconv.FormatInt(t.AccountId, 10)
	request := string(query)
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

func getCurrentMonth() (firstDay, lastDay string) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth.Format("2006-01-02"), lastOfMonth.Format("2006-01-02")
}

