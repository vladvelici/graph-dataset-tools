package main

import "container/list"

// Breath first search graph traversal.
// visit - sets a node as visited.
// visited - checks if a node was visited.
// f - function that is called at each visit, for convenience.
func Bfs(root *Node, visit func(*Node), visited func(*Node) bool, f func(*Node)) {
	todo := list.New()
	todo.PushBack(root)

	for todo.Len() > 0 {
		el := todo.Front()
		node := el.Value.(*Node)
		todo.Remove(el)

		if visited(node) {
			continue
		}

		visit(node)
		f(node)

		for _, ngh := range node.Neibourghs {
			todo.PushBack(ngh)
		}
	}
}

// IsPath with restrictions. Performs BFS but does not follow the restricted edge.
func IsPathRestricted(from, to *Node, restriction *Edge) bool {
	index := make(tr_index)

	todo := list.New()
	todo.PushBack(from)

	for todo.Len() > 0 {
		el := todo.Front()
		node := el.Value.(*Node)
		todo.Remove(el)

		if node == to {
			return true
		}

		if index.visited(node) {
			continue
		}

		index.visit(node)

		for _, ngh := range node.Neibourghs {
			if (node.Id == restriction.From && ngh.Id == restriction.To) ||
				(node.Id == restriction.To && ngh.Id == restriction.From) {
				// skip restricted edge
				continue
			}
			todo.PushBack(ngh)
		}
	}

	return false
}

// Depth first search keeping track of parents of nodes... Written to implement the minimum spanning tree.
func DfsEdge(root *Node, visit func(*Node), visited func(*Node) bool, f func(from, to *Node)) {
	todo := list.New()
	parents := list.New()
	todo.PushBack(root)
	parents.PushBack(nil)

	for todo.Len() > 0 {
		el := todo.Back()
		par := parents.Back()
		node := el.Value.(*Node)
		var parNode *Node
		if par.Value != nil {
			parNode = par.Value.(*Node)
		}
		todo.Remove(el)
		parents.Remove(par)

		if visited(node) {
			continue
		}

		visit(node)
		if parNode != nil {
			f(parNode, node)
		}

		for _, ngh := range node.Neibourghs {
			todo.PushBack(ngh)
			parents.PushBack(node)
		}
	}
}
