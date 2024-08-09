package cmd

import (
	s "main/structs"
	m "main/modal"
	"main/action"
	"database/sql"
	"strings"
	"strconv"
	"errors"
	"time"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: Diagram
func DrawStats(source *s.Source) {
	defer m.ErrorModal(source.Pages, source.Modal)
	
	var t s.Transaction
	var stats []s.Statistics
	
	box := tview.NewBox()
	box.SetBorder(true)
	
	firstDay, lastDay := getCurrentMonth()

	startDate := tview.NewInputField().SetLabel("Start Date").SetText(firstDay)
	endDate := tview.NewInputField().SetLabel("End Date").SetText(lastDay)

	dropDown := tview.NewDropDown().SetLabel("Account")
	
	// Flex
	topFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(dropDown, 0, 1, false).
		AddItem(startDate, 0, 1, false).
		AddItem(endDate, 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(topFlex, 0, 1, false).
		AddItem(box, 0, 15, true)
	
	accounts, a_types := action.SelectAccounts(source)
	
	if len(accounts) <= 0 {
		check(errors.New("Accounts len = 0"))
	}
	
	reloadChart := func(text string, index int) {
		action.SelectedAccount(text, index, a_types, &t)
		stats = getStatistics(source, &t, startDate.GetText(), endDate.GetText())
		flex.RemoveItem(flex.GetItem(1))
		flex.AddItem(box, 0, 15, true)
	}
	
	// InputFields
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

	box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		width -= 5
		height -= 5
		y += 5
		x += 5
		count := 0
		
		for _, stat := range stats {
			drawElement(screen, x+count, y, width, height/2, int(stat.Debit), tcell.ColorRed)
			drawElement(screen, x+count, y, width, height, int(stat.Credit), tcell.ColorGreen)
			count += 5
		}
		
		drawGraph(screen, x, y, width, height)
		drawGraph(screen, x, y, width, height/2)
		
		return x, y, width, height
	})
	
	source.Pages.AddPage("Stats", flex, true, true)
}

func drawGraph(screen tcell.Screen, xx, yy, width, height int) {

	for y := yy; y <= height; y++ {
		if y == height {
			screen.SetContent(xx, y, '\u2514', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		} else {
			screen.SetContent(xx, y, '\u2502', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}

	// Draw the bottom side of the cube
	for x := xx+1; x < width; x++ {
		screen.SetContent(x, height, '\u2500', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}

}

func drawElement(screen tcell.Screen, xx, yy, width, height int, percentage int, color tcell.Color) {
	percentageFloat := float64(percentage) / 100.0
	fullSize := float64(height - yy)
	percentageSize := int(fullSize * percentageFloat)
	size := height - percentageSize

	for y := height; y >= size; y-- {
		screen.SetContent(xx + 5, y, ' ', nil, tcell.StyleDefault.Background(color))
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