package cmd

import (
	s "main/structs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO: Diagram
func DrawStats(source *s.Source) {
	screen, err := tcell.NewScreen()
	check(err)

	err = screen.Init()
	check(err)

	defer screen.Fini()

	width, height := screen.Size()

	prim := tview.NewBox()
	prim.SetRect(0, 0, width, height)
	prim.SetBorder(true)

	prim.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		width -= 5
		height -= 5
		y += 5
		x += 5

		drawElement(screen, x, y, width, height, 20)
		drawElement(screen, x+5, y, width, height, 50)
		drawElement(screen, x+10, y, width, height, 10)
		drawElement(screen, x+15, y, width, height, 30)
		
		x, y, width, height = drawGraph(screen, x, y, width, height)

		return x, y, width, height
	})
	source.Pages.AddAndSwitchToPage("Stats", prim, false)
}

func drawGraph(screen tcell.Screen, xx, yy, width, height int) (int, int, int, int) {

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

	return xx, yy, width, height
}

func drawElement(screen tcell.Screen, xx, yy, width, height int, percentage int) {
	percentageFloat := float64(percentage) / 100.0
	fullSize := float64(height - yy)
	percentageSize := int(fullSize * percentageFloat)
	size := height - percentageSize

	for y := height; y >= size; y-- {
		screen.SetContent(xx + 5, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlue))
	}
}
