package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type node struct {
	text     string
	expand   bool
	selected func()
	children []*node
	parent   *tview.TreeNode
}

var (
	tree = tview.NewTreeView().SetAlign(false).SetTopLevel(0).SetGraphics(true).SetPrefixes(nil)
)

func MakeTree() *node {
	//var tableData []string

	var rootNode = &node{
		text: ".",
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
					//FillTable(tableData)
					table.SetBorder(true).SetTitle("Categories")
				}, expand: true, children: []*node{{text: "test"}}},
			}},
		}}
	return rootNode
}

func TreeView() *tview.TreeView {
	rootNode := MakeTree()
	tree.SetBorder(true).
		SetTitle("Category Tree")

		// Add nodes
	var add func(target *node, parent *tview.TreeNode) *tview.TreeNode

	add = func(target *node, parent *tview.TreeNode) *tview.TreeNode {

		node := tview.NewTreeNode(target.text).
			SetSelectable(target.expand || target.selected != nil).
			SetExpanded(target == rootNode).
			SetReference(target)

		if target.expand {
			node.SetColor(tcell.ColorPurple)
		} else if target.selected != nil {
			node.SetColor(tcell.ColorGreen)
		}
		if parent != nil {
			target.parent = parent
		}
		for _, child := range target.children {
			node.AddChild(add(child, node))
		}
		return node
	}

	root := add(rootNode, nil)
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

func RenameNode() {
	//node := tree.GetCurrentNode()
	//FillTreeAndListForm(node, nil)
	pages.AddPage("Dialog", Dialog(form), true, true)
}

func RemoveNode() {
	selected_node := tree.GetCurrentNode()
	if selected_node == nil {
		return
	}
	selected_node.ClearChildren()
	reference := selected_node.GetReference().(*node)
	reference.parent.RemoveChild(selected_node)
}
