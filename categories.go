package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expand   bool
	selected func()
	children []*node
}

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(1).SetGraphics(true).SetPrefixes(nil)
)

func MakeTree() *node{
	tableData := `OrderDate|Region|Rep|Item|Units|UnitCost|Total
1/6/2017|East|Jones|Pencil|95|1.99|189.05
1/23/2017|Central|Kivell|Binder|50|19.99|999.50`

	var rootNode = &node{
	text: "Root",
	children: []*node{
		{text: "Expand all", selected: func() { tree.GetRoot().ExpandAll() }},
		{text: "Collapse all", selected: func() {
			for _, child := range tree.GetRoot().GetChildren() {
				child.CollapseAll()
			}
		}},
		{text: "Root node", expand: true, children: []*node{
			{text: "Child node"},
			{text: "Child node"},
			{text: "Selected child node", selected: func() {
				// Updating table on selected node
				table.Clear()
				_ = FillTable(tableData)
				table.SetBorder(true).SetTitle("Categories")
			}},
		}},
	}}
	return rootNode
}

func TreeView() *tview.TreeView {
	rootNode := MakeTree()
	tree.SetBorder(true).
		SetTitle("Category Tree")

	// Add nodes.
	var add func(target *node) *tview.TreeNode
	add = func(target *node) *tview.TreeNode {
		node := tview.NewTreeNode(target.text).
			SetSelectable(target.expand || target.selected != nil).
			SetExpanded(target == rootNode).
			SetReference(target)
		if target.expand {
			node.SetColor(tcell.ColorPurple)
		} else if target.selected != nil {
			node.SetColor(tcell.ColorGreen)
		}
		for _, child := range target.children {
			node.AddChild(add(child))
		}
		return node
	}
	root := add(rootNode)
	tree.SetRoot(root).
		SetCurrentNode(root).
		SetSelectedFunc(func(n *tview.TreeNode) {
			original := n.GetReference().(*node)
			if original.expand {
				n.SetExpanded(!n.IsExpanded())
			} else if original.selected != nil {
				original.selected()
			}
		})

	return tree
}