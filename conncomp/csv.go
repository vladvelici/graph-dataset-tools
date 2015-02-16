package main

import (
	"io"
	"os"

	"github.com/vladvelici/graph-dataset-tools/util"
)

// Convenience function to read a graph from a file, by path.
func ReadGraph(path string) (*Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := util.NewReader(file)
	g := NewGraph()
	for {
		from, to, _, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return g, err
		}
		g.AddDirectedEdge(from, to)
	}
	return g, nil
}
