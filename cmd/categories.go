package cmd

import (
	"fmt"
	"main/action"
	m "main/modal"
	s "main/structs"

	"github.com/rivo/tview"
)

func MakeTree() *s.TreeNode {
	defer m.ErrorModal(source.Pages, source.Modal)

	_, _, category_nodes, err := action.SelectCategories(`SELECT * FROM Categories WHERE parent_id IS NULL`)
	check(err)

	for i, node := range category_nodes {
		query := fmt.Sprintf(`SELECT * FROM Categories WHERE parent_id = %v`, node.Reference.Id)
		_, _, children_nodes, err := action.SelectCategories(query)
		check(err)
		category_nodes[i].Children = children_nodes
	}

	var rootNode = &s.TreeNode{
		Text:     ".",
		Children: category_nodes,
	}
	return rootNode
}

func CategoryTree() {
	rootNode := MakeTree()
	source.CategoryTree.SetBorder(true).
		SetTitle("Category Tree")

	root := action.AddNode(rootNode, nil)
	source.CategoryTree.SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(n *tview.TreeNode) {
			original := n.GetReference().(*s.TreeNode)
			if original.Expand {
				n.SetExpanded(!n.IsExpanded())
			}
			if original.Selected != nil {
				original.Selected()
			}
		})

	source.CategoryTree.GetRoot().ExpandAll()
}