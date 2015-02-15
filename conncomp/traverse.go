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

		for _, ngh := range node.Neighbours {
			todo.PushBack(ngh)
		}
	}
}

// Depth first search keeping track of parents of nodes... Written to implement the minimum spanning tree.
func DfsEdge(root *Node, visit func(*Node), visited func(*Node) bool, f func(from, to *Node)) {
	todo := list.New()
	parents := list.New()
	todo.PushBack(root)
	parents.PushBack(root)

	for todo.Len() > 0 {
		el := todo.Back()
		par := parents.Back()
		node := el.Value.(*Node)
		var parNode *Node
		parNode = par.Value.(*Node)
		todo.Remove(el)
		parents.Remove(par)

		if visited(node) {
			continue
		}

		visit(node)
		f(parNode, node)

		for _, ngh := range node.Neighbours {
			if visited(ngh) {
				continue
			}
			todo.PushBack(ngh)
			parents.PushBack(node)
		}
	}
}
