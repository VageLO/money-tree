package modal

import (
	"fmt"

	"github.com/rivo/tview"
)

func Modal(p tview.Primitive, hight, width int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, hight, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func ErrorModal(pages *tview.Pages, modal *tview.Modal) {
	if r := recover(); r != nil {
		//if pages.HasPage("Modal") {
		//	pages.RemovePage("Modal")
		//}
		err := fmt.Sprintf("Error: %v", r)
		modal.SetText(err)
		pages.AddPage("Modal", Modal(modal, 20, 40), true, true)
	}
}
