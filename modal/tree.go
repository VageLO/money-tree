package modal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	s "main/structs"
	"os"
	"path/filepath"
)

type nodeReference struct {
	path   string
	isDir  bool
	parent *tview.TreeNode
}

type tree struct {
	*tview.TreeView
	rootNode *tview.TreeNode
}

var rootDir, _ = os.UserHomeDir()

func NewTree(source *s.Source, pattern string) *tree {
	defer ErrorModal(source.Pages, source.Modal)

	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed).
		SetReference(newNodeReference(rootDir, true, nil))

	tree := &tree{
		TreeView: tview.NewTreeView().
			SetRoot(root).
			SetCurrentNode(root),
		rootNode: root,
	}

	tree.addNode(root, rootDir, pattern)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		tree.expandOrAddNode(node, source, pattern)
	})

	return tree
}

func (tree *tree) addNode(directoryNode *tview.TreeNode, path string, pattern string) {
	files, err := os.ReadDir(path)
	check(err)

	for _, file := range files {
		if filepath.Ext(file.Name()) != pattern {
			continue
		}
		node := createTreeNode(file.Name(), file.IsDir(), directoryNode)
		directoryNode.AddChild(node)
	}
}

func (tree tree) expandOrAddNode(node *tview.TreeNode, source *s.Source, pattern string) {
	defer ErrorModal(source.Pages, source.Modal)

	reference := node.GetReference()
	if reference == nil {
		return
	}

	nodeReference := reference.(*nodeReference)
	if !nodeReference.isDir && pattern == "" {
		source.Attachments = append(source.Attachments, nodeReference.path)
		source.Pages.RemovePage("FileExplorer")
		return

	} else if !nodeReference.isDir && pattern != "" {
		source.Imports = append(source.Imports, nodeReference.path)
		source.Pages.RemovePage("FileExplorer")
		return

	}

	children := node.GetChildren()
	if len(children) == 0 {
		path := nodeReference.path
		tree.addNode(node, path, pattern)
	} else {
		node.SetExpanded(!node.IsExpanded())
	}
}

func createTreeNode(fileName string, isDir bool, parent *tview.TreeNode) *tview.TreeNode {
	var parentPath string

	if parent == nil {
		parentPath = rootDir
	} else {
		reference, ok := parent.GetReference().(*nodeReference)
		if !ok {
			parentPath = rootDir
		} else {
			parentPath = reference.path
		}
	}

	var color tcell.Color
	if isDir {
		color = tcell.ColorGreen
	} else {
		color = tview.Styles.PrimaryTextColor
	}

	return tview.NewTreeNode(fileName).
		SetReference(
			newNodeReference(
				filepath.Join(parentPath, fileName),
				isDir,
				parent,
			),
		).
		SetSelectable(true).
		SetColor(color)
}

func newNodeReference(path string, isDir bool, parent *tview.TreeNode) *nodeReference {
	return &nodeReference{
		path:   path,
		isDir:  isDir,
		parent: parent,
	}
}
