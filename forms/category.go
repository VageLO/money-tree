package forms

import (
	s "main/structs"
	
	"github.com/rivo/tview"
)

func FillNodeForm(node *tview.TreeNode) {
	form.Clear(true)

	FormStyle("Category Information")
	title := node.GetText()
	form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameNode(text, node) })
}

func FormAddCategory() {
	form.Clear(true)
	formStyle("Add Category")
	root := tree.GetRoot()

	n := &TreeNode{
		text:   "",
		expand: true,
		reference: &s.Category{},
		children: []*TreeNode{},
	}
	new_node := add(n, root)

	form.AddInputField("Title: ", "", 0, nil, func(text string) { 
		new_node.SetText(text)
	})
	
	var selected_dropdown *tview.TreeNode
	var options []string
	options = append(options, root.GetText())

	for _, children := range root.GetChildren() {
		options = append(options, children.GetText())
	}

	initial := 0

	selected_node := tree.GetCurrentNode()
	if selected_node != nil {
		for idx, title := range options {
			if title == selected_node.GetText() {
				initial = idx
			}
		}
	}

	form.AddDropDown("Categories", options, initial, func(option string, optionIndex int) {
		if root.GetText() == option {
			selected_dropdown = root
			reference := new_node.GetReference().(*TreeNode)
			reference.parent = root
			new_node.SetReference(reference)
			return
		}
		
		for _, children := range root.GetChildren() {
			if children.GetText() == option {
				selected_dropdown = children
				reference := new_node.GetReference().(*TreeNode)
				reference.parent = children
				new_node.SetReference(reference)
				return
			}
		}
	})

	form.AddButton("Add", func() {
		AddCategory(new_node, selected_dropdown)
		pages.RemovePage("Form")
	})
	pages.AddPage("Form", Modal(form, 30, 50), true, true)
}

func FormRenameCategory() {
	node := tree.GetCurrentNode()
	if node == nil {
		return
	}
	FillNodeForm(node)
	pages.AddPage("Form", Modal(form, 30, 50), true, true)
}

func SelectedCategory(option string, optionIndex int, c_types []s.Category, t *s.Transaction) {
	selected_c := c_types[optionIndex]
	if selected_c.title != option {
		return
	}
	t.category_id = selected_c.id
	t.category = selected_c.title
}