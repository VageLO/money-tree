package action

import (
	"fmt"
	"github.com/rivo/tview"
)

func locateCheckboxes(form *tview.Form, index int) {
	if index > 0 {
		index = index + 1
	}
	form.AddCheckbox(fmt.Sprintf("check %v", index), false, func(checked bool) {
		switch checked {
		case true:
			item := form.GetFormItem(index)
			item.SetDisabled(false)
		case false:
			item := form.GetFormItem(index)
			item.SetDisabled(true)
		}
	})
	item := form.GetFormItem(index)
	item.SetDisabled(true)
}
