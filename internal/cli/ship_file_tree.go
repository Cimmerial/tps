package cli

import (
	"sort"
	"strings"
)

const (
	shipTreeBranch = "├─ "
	shipTreeLast   = "└─ "
	shipTreeVert   = "│  "
	shipTreeSpace  = "   "
)

type shipFileTreeEntry struct {
	Number int
	Path   string
}

type shipFileTreeNode struct {
	name     string
	dir      bool
	number   int
	relPath  string
	children []*shipFileTreeNode
}

func buildShipFileTree(entries []shipFileTreeEntry) []*shipFileTreeNode {
	root := &shipFileTreeNode{name: "", dir: true}
	for _, e := range entries {
		p := strings.TrimSpace(strings.ReplaceAll(e.Path, "\\", "/"))
		if p == "" {
			continue
		}
		parts := strings.Split(p, "/")
		cur := root
		for i, part := range parts {
			last := i == len(parts)-1
			child := findShipTreeChild(cur, part)
			if child == nil {
				child = &shipFileTreeNode{name: part, dir: !last}
				cur.children = append(cur.children, child)
			}
			if last {
				child.dir = false
				child.number = e.Number
				child.relPath = p
			} else {
				child.dir = true
			}
			cur = child
		}
	}
	sortShipTreeChildren(root)
	return root.children
}

func findShipTreeChild(n *shipFileTreeNode, name string) *shipFileTreeNode {
	for _, c := range n.children {
		if c.name == name {
			return c
		}
	}
	return nil
}

func sortShipTreeChildren(n *shipFileTreeNode) {
	if len(n.children) == 0 {
		return
	}
	sort.Slice(n.children, func(i, j int) bool {
		a, b := n.children[i], n.children[j]
		if a.dir != b.dir {
			return a.dir
		}
		return a.name < b.name
	})
	for _, c := range n.children {
		sortShipTreeChildren(c)
	}
}

func shipFileEntriesFromPaths(paths []string) []shipFileTreeEntry {
	entries := make([]shipFileTreeEntry, len(paths))
	for i, p := range paths {
		entries[i] = shipFileTreeEntry{Path: p}
	}
	return entries
}

// shipFileDisplayNumbers assigns 1..N in depth-first tree order (dirs then files, alphabetical).
func shipFileDisplayNumbers(paths []string) map[string]int {
	tree := buildShipFileTree(shipFileEntriesFromPaths(paths))
	nums := make(map[string]int, len(paths))
	n := 0
	var walk func([]*shipFileTreeNode)
	walk = func(nodes []*shipFileTreeNode) {
		for _, node := range nodes {
			if node.dir {
				walk(node.children)
			} else {
				n++
				nums[node.relPath] = n
			}
		}
	}
	walk(tree)
	return nums
}

func shipFilePathForDisplayNumber(paths []string, displayNum int) (string, bool) {
	tree := buildShipFileTree(shipFileEntriesFromPaths(paths))
	cur := 0
	var found string
	var walk func([]*shipFileTreeNode) bool
	walk = func(nodes []*shipFileTreeNode) bool {
		for _, node := range nodes {
			if node.dir {
				if walk(node.children) {
					return true
				}
			} else {
				cur++
				if cur == displayNum {
					found = node.relPath
					return true
				}
			}
		}
		return false
	}
	walk(tree)
	return found, found != ""
}

func shipFileTreeLines(nodes []*shipFileTreeNode, dirLine func(name string) string, leafLine func(number int, fileName, relPath string) string) []string {
	var out []string
	var walk func([]*shipFileTreeNode, string)
	walk = func(list []*shipFileTreeNode, prefix string) {
		for i, n := range list {
			last := i == len(list)-1
			branch := shipTreeBranch
			if last {
				branch = shipTreeLast
			}
			if n.dir {
				out = append(out, prefix+branch+dirLine(n.name))
				childPrefix := prefix + shipTreeVert
				if last {
					childPrefix = prefix + shipTreeSpace
				}
				walk(n.children, childPrefix)
			} else {
				out = append(out, prefix+branch+leafLine(n.number, n.name, n.relPath))
			}
		}
	}
	walk(nodes, "")
	return out
}
