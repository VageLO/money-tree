package action

import (
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"strings"

	"github.com/rivo/tview"
)

func FormAddCategory(source *s.Source) {
	form := source.Form
	pages := source.Pages
	tree := source.CategoryTree

	form.Clear(true)
	FormStyle("Add Category", form)

	root := tree.GetRoot()

	n := &s.TreeNode{
		Text:      "",
		Expand:    true,
		Reference: &s.Category{},
		Children:  []*s.TreeNode{},
	}
	newNode := AddNode(n, root)

	form.AddInputField("Title: ", "", 0, nil, func(text string) {
		if text == "" {
			return
		}
		newNode.SetText(strings.TrimSpace(text))
	})

	var selectedDropdown *tview.TreeNode
	var options []string
	options = append(options, root.GetText())

	for _, children := range root.GetChildren() {
		options = append(options, children.GetText())
	}

	initial := 0

	selectedNode := tree.GetCurrentNode()
	if selectedNode != nil {
		for idx, title := range options {
			if title == selectedNode.GetText() {
				initial = idx
			}
		}
	}

	form.AddDropDown("Categories", options, initial, func(option string, optionIndex int) {
		if root.GetText() == option {
			selectedDropdown = root
			reference := newNode.GetReference().(*s.TreeNode)
			reference.Parent = root
			newNode.SetReference(reference)
			return
		}

		for _, children := range root.GetChildren() {
			if children.GetText() == option {
				selectedDropdown = children
				reference := newNode.GetReference().(*s.TreeNode)
				reference.Parent = children
				newNode.SetReference(reference)
				return
			}
		}
	})

	form.AddButton("Add", func() {
		AddCategory(newNode, selectedDropdown, source)
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

	form.Clear(true)
	FormStyle("Category Details", form)

	title := node.GetText()
	form.AddInputField("Title: ", title, 0, nil, func(text string) {
		if text == "" {
			return
		}
		title = strings.TrimSpace(text)
	})

	form.AddButton("Save", func() {
		RenameNode(title, node, source)
	})

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
