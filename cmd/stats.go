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

		for y := y; y <= height; y++ {
			if y == height {
				screen.SetContent(x, y, '\u2514', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			} else {
				screen.SetContent(x, y, '\u2502', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
		}

		// Draw the bottom side of the cube
		y = height
		for x := x; x < width; x++ {
			screen.SetContent(x, y, '\u2500', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}

		return x, y, width, height
	})
	source.Pages.AddAndSwitchToPage("Stats", prim, false)
}
