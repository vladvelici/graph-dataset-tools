package main

import "math/rand"

// Type to represent an edge set.
type Edge struct {
	From, To int
}

// Represents a node in a graph.
type Node struct {
	Id         int
	Neibourghs map[int]*Node
}

// Create a new node.
func NewNode(id, set int) *Node {
	return &Node{
		id,
		make(map[int]*Node),
	}
}

// Represents an undirected graph that forms connected graphs.
type Graph struct {
	Nodes map[int]*Node
}

// Create a new graph.
func NewGraph() *Graph {
	return &Graph{
		make(map[int]*Node),
	}
}

// Fetches the node by ID. Creates the node if it does not exist.
func (g *Graph) fetch(node int) *Node {
	f, ok := g.Nodes[node]
	if !ok || f == nil {
		f = NewNode(node, -1)
		g.Nodes[node] = f
	}
	return f
}

// List of edges to be printed out.
func (g *Graph) EdgeList() []*Edge {
	res := make([]*Edge, 0)
	for nid, node := range g.Nodes {
		for toId, _ := range node.Neibourghs {
			res = append(res, &Edge{nid, toId})
		}
	}
	return res
}

// Add a directed edge, taking care of connected components.
// Panics if either from or to are nil.
func (g *Graph) directedEdge(from, to *Node) {
	if _, ok := from.Neibourghs[to.Id]; ok {
		return
	}
	from.Neibourghs[to.Id] = to
}

// Add an undirected edge to the graph.
func (g *Graph) AddEdge(fromId, toId int) {
	from := g.fetch(fromId)
	to := g.fetch(toId)
	g.directedEdge(from, to)
	g.directedEdge(to, from)
}

// Add a node to the graph. Used internally for connected components.
func (g *Graph) addNode(node *Node) {
	g.Nodes[node.Id] = node
}

// Internal type used for traversals.
type tr_index map[int]interface{}

// Set node as visited.
func (t tr_index) visit(node *Node) {
	t[node.Id] = nil
}

// Check if a node was visited.
func (t tr_index) visited(node *Node) bool {
	_, ok := t[node.Id]
	return ok
}

// Return a list of connected graphs.
func (g *Graph) ConnectedGraphs() []*Graph {
	result := make([]*Graph, 0)

	index := make(tr_index)

	for _, node := range g.Nodes {
		if index.visited(node) {
			continue
		}
		graph := NewGraph()
		Bfs(node, index.visit, index.visited, graph.addNode)
		result = append(result, graph)
	}

	return result
}

// Minimum spanning tree. Acually an undirected edge hash set.
type Mst map[int]map[int]interface{}

// from always less than to; from<to
func (m Mst) Add(from, to int) {
	if from > to {
		from, to = to, from
	}
	f := m[from]
	if f == nil {
		m[from] = make(map[int]interface{})
	}
	m[from][to] = nil
}

// Check if the Mst has the edge.
func (m Mst) Has(from, to int) bool {
	if from > to {
		from, to = to, from
	}
	if f, ok := m[from]; ok && f != nil {
		_, ok2 := f[to]
		return ok2
	}
	return false
}

// Minimum spanning tree computation.
func (g *Graph) Mst() Mst {
	var root *Node
	for _, node := range g.Nodes {
		root = node
		break
	}

	if root == nil {
		return nil
	}

	index := make(tr_index)
	res := make(Mst)

	DfsEdge(root, index.visit, index.visited, func(from, to *Node) {
		res.Add(from.Id, to.Id)
	})

	return res
}

// Remove random edges, keeping track of them.
//
// It does not remove critical edges. Computes the minimum spanning tree and only removes
// random edges that are not part of it.
//
// Uses reservoir sampling.
func (g *Graph) RemoveRandomEdges(n int, restrictions Mst) []*Edge {
	result := make([]*Edge, 0, n)

	defer func() {
		// do the work of actually removing those edges
		for _, edge := range result {
			delete(g.Nodes[edge.From].Neibourghs, edge.To)
			delete(g.Nodes[edge.To].Neibourghs, edge.From)
		}
	}()

	ch := make(chan *Edge)
	go func(channel chan *Edge, ignore Mst) {
		seen := make(Mst)
		for _, node := range g.Nodes {
			for to, _ := range node.Neibourghs {
				if seen.Has(node.Id, to) || ignore.Has(node.Id, to) {
					continue
				}
				channel <- &Edge{node.Id, to}
				seen.Add(node.Id, to)
			}
		}
		close(channel)
	}(ch, restrictions)

	for i := 0; i < n; i++ {
		edge, ok := <-ch
		if !ok {
			return result
		}
		result = append(result, edge)
	}

	for node := range ch {
		n++
		rnd := rand.Intn(n)
		if rnd < len(result) {
			result[rnd] = node
		}
	}

	return result
}
