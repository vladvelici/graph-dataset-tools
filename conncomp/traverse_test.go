package main

import "testing"

var traversalType = [][]int{
	{1, 2}, // 0
	{0, 3}, // 1
	{0, 4}, // 2
	{1},    // 3
	{2},    // 4
}

var cyclic = [][]int{
	{1, 2},    // 0
	{0, 3, 4}, // 1
	{0, 4},    // 2
	{1},       // 3
	{2, 1},    // 4
}

func nodeSlice(nodes [][]int) []*Node {
	result := make([]*Node, len(nodes))
	for i, _ := range result {
		result[i] = NewNode(i)
	}
	for i, neigh := range nodes {
		node := result[i]
		for _, j := range neigh {
			node.Neighbours[j] = result[j]
		}
	}
	return result
}

func TestNodeSlice(t *testing.T) {
	raw := [][]int{
		{1},
		{0},
	}
	nodes := nodeSlice(raw)
	if len(nodes) != 2 {
		t.FailNow()
	}
	if nodes[0].Neighbours[1] != nodes[1] {
		t.Error("Node 0 does not have node 1 as neighbours")
	}
	if nodes[1].Neighbours[0] != nodes[0] {
		t.Error("Node 1 does not have node 0 as neighbours")
	}
	if len(nodes[1].Neighbours) != 1 || len(nodes[0].Neighbours) != 1 {
		t.Error("Wrong number of neighbours")
	}
}

// Make sure the traversal type in TestDfs is depth first
func TestDfsEdgeTraversalType(t *testing.T) {
	nodes := nodeSlice(traversalType)
	visits := make([]bool, len(nodes))

	next := []*Node{nodes[0]}

	visit := func(n *Node) {
		t.Logf("Visiting node %d.\n", n.Id)

		if visits[n.Id] == true {
			t.Errorf("Re-visiting node %d!", n.Id)
		}

		visits[n.Id] = true

		if len(next) == 0 && !(n == nodes[1] || n == nodes[2]) {
			t.Errorf("Bad state. No next and visiting %d", n.Id)
		}

		if len(next) > 0 && next[0] != n {
			t.Errorf("Wrong next expected visit. Expecting %d, got %d.", next[0].Id, n.Id)
		}

		if len(next) > 0 {
			next = next[1:]
		}

		if n == nodes[1] {
			next = append(next, nodes[3], nodes[2], nodes[4])
		} else if n == nodes[2] {
			next = append(next, nodes[4], nodes[1], nodes[3])
		}
	}

	visited := func(n *Node) bool {
		return visits[n.Id]
	}

	DfsEdge(nodes[0], visit, visited, func(a, b *Node) {})
}

func TestDfdEdgeParents(t *testing.T) {
	nodes := nodeSlice(traversalType)
	visits := make([]bool, len(nodes))

	visit := func(n *Node) {
		visits[n.Id] = true
	}

	visited := func(n *Node) bool {
		return visits[n.Id]
	}

	f := func(from, to *Node) {
		var expected *Node

		switch to {
		case nodes[0]:
			expected = nodes[0]
		case nodes[1]:
			expected = nodes[0]
		case nodes[2]:
			expected = nodes[0]
		case nodes[3]:
			expected = nodes[1]
		case nodes[4]:
			expected = nodes[2]
		}

		t.Logf("Seeing edge: %d -> %d", from.Id, to.Id)
		if from != expected {
			t.Errorf("Seeing edge: %d -> %d, where %d should have parent %d.", from.Id, to.Id, to.Id, expected.Id)
		}
	}

	DfsEdge(nodes[0], visit, visited, f)
}

func TestBfsTraversal(t *testing.T) {
	nodes := nodeSlice(traversalType)
	visits := make([]bool, len(nodes))

	next := []*Node{nodes[0]}

	visit := func(n *Node) {
		t.Logf("Visiting node %d.\n", n.Id)

		if visits[n.Id] == true {
			t.Errorf("Re-visiting node %d!", n.Id)
		}

		visits[n.Id] = true

		if len(next) == 0 && !(n == nodes[1] || n == nodes[2]) {
			t.Errorf("Bad state. No next and visiting %d", n.Id)
		}

		if len(next) > 0 && next[0] != n {
			t.Errorf("Wrong next expected visit. Expecting %d, got %d.", next[0].Id, n.Id)
		}

		if len(next) > 0 {
			next = next[1:]
		}

		if n == nodes[1] {
			next = append(next, nodes[2], nodes[3], nodes[4])
		} else if n == nodes[2] {
			next = append(next, nodes[1], nodes[4], nodes[3])
		}
	}

	visited := func(n *Node) bool {
		return visits[n.Id]
	}

	Bfs(nodes[0], visit, visited, func(n *Node) {})
}

func TestDfsCycles(t *testing.T) {
	nodes := nodeSlice(cyclic)
	visits := make([]bool, len(nodes))

	visit := func(n *Node) {
		t.Logf("Visiting node %d.\n", n.Id)

		if visits[n.Id] == true {
			t.Errorf("Re-visiting node %d!", n.Id)
		}

		visits[n.Id] = true
	}

	visited := func(n *Node) bool {
		return visits[n.Id]
	}

	DfsEdge(nodes[0], visit, visited, func(a, b *Node) {})
}

func TestBfsCycles(t *testing.T) {
	nodes := nodeSlice(cyclic)
	visits := make([]bool, len(nodes))

	visit := func(n *Node) {
		t.Logf("Visiting node %d.\n", n.Id)

		if visits[n.Id] == true {
			t.Errorf("Re-visiting node %d!", n.Id)
		}

		visits[n.Id] = true
	}

	visited := func(n *Node) bool {
		return visits[n.Id]
	}

	Bfs(nodes[0], visit, visited, func(n *Node) {})
}

// tr_index test
func TestTrIndex(t *testing.T) {
	index := make(tr_index)
	nodes := nodeSlice(cyclic)

	index.visit(nodes[1])
	index.visit(nodes[3])

	if index.visited(nodes[4]) {
		t.Error("Node 4 should not be visited.")
	}

	if !index.visited(nodes[1]) {
		t.Error("Node 1 should be visited.")
	}

}
