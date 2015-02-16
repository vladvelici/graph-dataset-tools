package main

import "testing"

func mkgraph(arr [][]int) *Graph {
	g := NewGraph()
	for from, toList := range arr {
		for _, to := range toList {
			g.AddEdge(from, to)
		}
	}
	return g
}

func TestCreateGraph(t *testing.T) {
	g := NewGraph()

	edgeIndex := make(map[Edge]bool)

	totalEdges := 0
	for from, toList := range traversalType {
		for _, to := range toList {
			edgeIndex[Edge{from, to}] = false
			g.AddEdge(from, to)
			totalEdges++
		}
	}

	if len(traversalType) != len(g.Nodes) {
		t.Errorf("Expected %d nodes, only have %d.", len(traversalType), len(g.Nodes))
	}

	graphEdges := g.EdgeList()
	if totalEdges != len(graphEdges) {
		t.Errorf("Expected %d edges, but only found %d.", totalEdges, len(graphEdges))
	}

	for _, edge := range graphEdges {
		if found, ok := edgeIndex[*edge]; found || !ok {
			t.Errorf("Couldn't find edge or edge repeated %#v.", *edge)
		}
		edgeIndex[*edge] = true
	}
}

func TestConnectedGraphsSmall(t *testing.T) {
	g := mkgraph(traversalType)
	conn := g.ConnectedGraphs()
	if len(conn) != 1 {
		t.Error("More then one connected graph, when obviously there is only one.")
	}
}

var threeConnectedGraphs = [][]int{
	{1, 2},    // 0
	{0, 3, 2}, // 1
	{0, 1},    // 2
	{1},       // 3
	{5, 6},    // 4
	{4, 6},    // 5
	{5, 4},    // 6
	{8},       // 7
	{7},       // 8

}

func TestConnectedGraphsHard(t *testing.T) {
	g := mkgraph(threeConnectedGraphs)
	conn := g.ConnectedGraphs()
	if len(conn) != 3 {
		t.Error("Should have three connected graphs. Found %d", len(conn))
	}

	funcs := map[int]func(*Node) bool{
		4: func(node *Node) bool { return node.Id <= 3 },
		3: func(node *Node) bool { return node.Id >= 4 && node.Id <= 6 },
		2: func(node *Node) bool { return node.Id >= 7 && node.Id <= 8 },
	}

	gn := func(g *Graph) *Node {
		for _, n := range g.Nodes {
			return n
		}
		return nil
	}

	for _, c := range conn {
		index := make(tr_index)
		root := gn(c)
		Bfs(root, index.visit, index.visited, func(node *Node) {
			if c.Nodes[node.Id] != node {
				t.Errorf("CONN %d-el: Reached node %d outside its own graph.", len(c.Nodes), node.Id)
			}
			if funcs[len(c.Nodes)](node) == false {
				t.Errorf("CONN %d-el: Node %d reached outside its connected component.", len(c.Nodes), node.Id)
			}
		})
	}
}

func TestIsUndirected(t *testing.T) {
	g := mkgraph([][]int{
		{1},
		{0},
	})
	if !g.IsUndirected() {
		t.Error("Undirected graph is said not to be.")
	}

	// force an undirected graph.
	g = NewGraph()
	g.Nodes[0] = NewNode(0)
	g.Nodes[1] = NewNode(1)
	g.directedEdge(g.Nodes[0], g.Nodes[1])

	if g.IsUndirected() {
		t.Error("Directed graph is said to be undirected.")
	}
}

func TestIsConnected(t *testing.T) {
	g := mkgraph(threeConnectedGraphs)
	if g.IsConnected() {
		t.Error("Three connected graphs is not connected.")
	}
	conn := g.ConnectedGraphs()
	for _, c := range conn {
		if c.IsConnected() == false {
			t.Errorf("Connected graph with %d nodes reported disconnected.", len(c.Nodes))
		}
	}
	g = NewGraph()
	g.Nodes[0] = NewNode(0)
	g.Nodes[1] = NewNode(1)
	if g.IsConnected() == true {
		t.Errorf("Two empty nodes are not connected.")
	}
}

func TestMstOnTree(t *testing.T) {
	g := mkgraph(traversalType)

	mst := g.Mst()

	must := []*Edge{
		{0, 1},
		{1, 3},
		{0, 2},
		{2, 4},
	}

	// let's log the mst
	t.Logf("Printing out Mst edges: \n")
	for from, dests := range mst {
		for to, _ := range dests {
			t.Logf("\t %d -> %d", from, to)
		}
	}

	for _, e := range must {
		if !mst.Has(e.From, e.To) {
			t.Errorf("Mst does not have edge %#v", *e)
		}
	}
}

func TestMstCyclic(t *testing.T) {
	g := mkgraph(connectedGraph)
	mst := g.Mst()

	g2 := NewGraph()

	index := make([]bool, len(connectedGraph))

	t.Logf("Printing out Mst edges: \n")
	for from, dests := range mst {
		for to, _ := range dests {
			t.Logf("\t %d -> %d", from, to)
			g2.AddEdge(from, to)
			index[from] = true
			index[to] = true
		}
	}

	if !g2.IsConnected() {
		t.Error("Built up graph is not connected.")
	}

	if !g2.IsUndirected() {
		t.Error("The graph is not undirected.")
	}

	for nid, check := range index {
		if !check {
			t.Errorf("Node %d not reached for mst.", nid)
		}
	}

}

// Check if an edge is in the list of edges.
func isIn(possible []Edge, e *Edge) bool {
	for _, p := range possible {
		if p == *e {
			return true
		}
	}
	return false
}

var connectedGraph = [][]int{
	{1, 2, 8},       // 0
	{0, 3, 2},       // 1
	{0, 1, 7, 8, 4}, // 2
	{1},             // 3
	{5, 6, 2},       // 4
	{4, 6},          // 5
	{5, 4},          // 6
	{8, 2},          // 7
	{7, 2, 0},       // 8
}

func TestRemoveRandom(t *testing.T) {
	g := mkgraph(connectedGraph)
	mst := g.Mst()
	rmv := g.RemoveRandomEdges(10, mst)

	t.Logf("Removed the following edges: ")
	for _, e := range rmv {
		t.Logf("\t %#v", *e)
		// check if the edge is actually removed
		if _, ok := g.Nodes[e.From].Neighbours[e.To]; ok {
			t.Errorf(" ^-- edge not removed.")
		}
	}

	if len(rmv) != 4 {
		t.Errorf("Removed incorrect number of edges (%d) from connected graph.", len(rmv))
	}

	if !g.IsConnected() {
		t.Error("The graph is not connected anymore.")
	}

	if !g.IsUndirected() {
		t.Error("The graph is not undirected.")
	}
}
