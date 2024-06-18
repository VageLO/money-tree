package cmd

import (
	"fmt"
	
	"github.com/rivo/tview"
)

func Dialog(p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, 35, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)
}

func Modal(p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, 35, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)
}

func ErrorModal() {
	if r := recover(); r != nil {
		//if pages.HasPage("Dialog") {
		//	pages.RemovePage("Dialog")
		//}
		err := fmt.Sprintf("Error: %v", r)
		modal.SetText(err)
		pages.AddPage("Modal", Modal(modal), true, true)
	}
}