package action

import (
	s "main/structs"
	m "main/modal"
	
	"github.com/rivo/tview"
)

func FillNodeForm(node *tview.TreeNode, source *s.Source) {
	form := source.Form
	form.Clear(true)

	formStyle("Category Information", form)
	title := node.GetText()
	form.AddInputField("Title: ", title, 0, nil, func(text string) { RenameNode(text, node) })
}

func FormAddCategory(source *s.Source) {
	form := source.Form
	pages := source.Pages
	tree := source.CategoryTree
	
	form.Clear(true)
	formStyle("Add Category", form)
	
	root := tree.GetRoot()

	n := &s.TreeNode{
		Text:   "",
		Expand: true,
		Reference: &s.Category{},
		Children: []*s.TreeNode{},
	}
	new_node := AddNode(n, root)

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
			reference := new_node.GetReference().(*s.TreeNode)
			reference.Parent = root
			new_node.SetReference(reference)
			return
		}
		
		for _, children := range root.GetChildren() {
			if children.GetText() == option {
				selected_dropdown = children
				reference := new_node.GetReference().(*s.TreeNode)
				reference.Parent = children
				new_node.SetReference(reference)
				return
			}
		}
	})

	form.AddButton("Add", func() {
		AddCategory(new_node, selected_dropdown)
		pages.RemovePage("Form")
	})
	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func FormRenameCategory(source *s.Source) {
	form := source.Form
	pages := source.Pages
	tree := source.CategoryTree
	
	node := tree.GetCurrentNode()
	if node == nil {
		return
	}
	FillNodeForm(node, source)
	pages.AddPage("Form", m.Modal(form, 30, 50), true, true)
}

func SelectedCategory(option string, optionIndex int, c_types []s.Category, t *s.Transaction) {
	selected_c := c_types[optionIndex]
	if selected_c.Title != option {
		return
	}
	t.CategoryId = selected_c.Id
	t.Category = selected_c.Title
}