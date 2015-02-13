package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Define flags.
var (
	flagRemove  = flag.Int("remove", 0, "Number of edges to remove from each graph.")
	flagGraph   = flag.String("graph", "-", "Input file for an edge list csv.")
	flagOutput  = flag.String("o", "processed", "Output file base name. This program actually outputs many files, depending on the number of connected components in the graph.")
	flagVerbose = flag.Bool("verbose", false, "Whether to print lots of debug information on stdout.")
)

func main() {
	flag.Parse()

	var inputStream *bufio.Reader
	if *flagGraph == "-" {
		inputStream = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(*flagGraph)
		if err != nil {
			fmt.Printf("Error reading input file %s\n", err.Error())
		}
		inputStream = bufio.NewReader(file)
		defer file.Close()
	}

	csvReader := csv.NewReader(inputStream)

	graph := NewGraph()

	for {
		from, to, err := record(csvReader)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error ", err)
			return
		}
		if *flagVerbose {
			fmt.Printf("Reading pair\t %d - %d\n", from, to)
		}
		graph.AddEdge(from, to)
	}

	fmt.Fprintf(os.Stderr, "Finished reading graph.\n")

	connectedGraphs := graph.ConnectedGraphs()

	fmt.Fprintf(os.Stderr, "Found %d connected graphs. Now processing them one by one... \n", len(connectedGraphs))
	var wg sync.WaitGroup
	for i, g := range connectedGraphs {
		wg.Add(1)
		go func(g *Graph, filename string, i, total int) {
			defer wg.Done()
			oneGraph(g, filename, i, total)
		}(g, *flagOutput+strconv.Itoa(i+1), i+1, len(connectedGraphs))
	}
	wg.Wait()
}

func oneGraph(g *Graph, filename string, i, total int) {
	if *flagVerbose {
		fmt.Println("Starting graph #", i)
	}

	var (
		outputEdgesFilename   = filename + "_edges.csv"
		outputRemovedFilename = filename + "_removed.csv"
	)

	// opening output files...
	outputEdges, err := os.Create(outputEdgesFilename)
	if err != nil {
		fmt.Printf("Graph #%d: Unable to open 'edges' output file %s. ( %s )", i, outputEdgesFilename, err)
		return
	}
	outputRemoved, err := os.Create(outputRemovedFilename)
	if err != nil {
		fmt.Printf("Graph #%d: Unable to open 'removed' output file %s. ( %s )", i, outputRemovedFilename, err)
		return
	}

	spanningTree := g.Mst()
	fmt.Fprintf(os.Stderr, "Computed spanning tree for graph %d.\n", i)

	removed := g.RemoveRandomEdges(*flagRemove, spanningTree)
	fmt.Fprintf(os.Stderr, "Some random edges removed for graph #%d.\n", i)

	for _, edge := range removed {
		if *flagVerbose {
			fmt.Printf("%d, %d, \"removed\"\n", edge.From, edge.To)
		}
		_, err = fmt.Fprintf(outputRemoved, "%d, %d\n", edge.From, edge.To)
		if err != nil {
			fmt.Printf("Graph #%d: Error printing removed edge. (%s)", i, err)
		}
	}

	allEdges := g.EdgeList()
	for _, edge := range allEdges {
		if *flagVerbose {
			fmt.Printf("%d, %d, \"remaining\"\n", edge.From, edge.To)
		}
		_, err = fmt.Fprintf(outputEdges, "%d, %d\n", edge.From, edge.To)
		if err != nil {
			fmt.Printf("Graph #%d: Error printing remaining edge. (%s)", i, err)
		}
	}

	fmt.Fprintf(os.Stderr, "Finished processing graph number %d of %d.\n", i, total)
}

func record(r *csv.Reader) (int, int, error) {
	rec, err := r.Read()
	if err != nil {
		return 0, 0, err
	}
	if len(rec) != 2 {
		return 0, 0, fmt.Errorf("Wrong record %#v.", rec)
	}
	rec[0], rec[1] = strings.TrimSpace(rec[0]), strings.TrimSpace(rec[1])
	a, err := strconv.Atoi(rec[0])
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.Atoi(rec[1])
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}
