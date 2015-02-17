package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/vladvelici/graph-dataset-tools/util"
)

// Define flags.
var (
	flagAction  = flag.String("action", "", "Possibile actions: details, components, remove")
	flagN       = flag.Float64("n", 0, "\\% of edges to remove from each graph.")
	flagOutput  = flag.String("o", "component_", "Output file prefix. It will be followed by the component number and an underscore.")
	flagVerbose = flag.Bool("verbose", false, "Whether to print lots of debug information on stdout.")
	flagAll     = flag.Bool("all", false, "Whether to print all edges (to->from and from->to) when removing, in the removed list.")
	flagForce   = flag.Bool("force", false, "Do not make safety checks. Might be faster.")
	flagHelp    = flag.Bool("help", false, "Show this help message")
	flagH       = flag.Bool("h", false, "Show this help message")
)

var helpMessage = `conncomp is a tool that deals with connected components tasks.

Possible uses (can optionally add -o and -verbose flags before filelist):

-action details filelist			Outputs details about the graphs in the files. Number of edges, number of nodes, directed/undirected, connected and no. of components.
-action components filelist			Splits the given graph in connected components files (prefix can be changed with -o flag).
-action remove -n P [-all] filelist	Removes at most P% random edges from the given graphs. Prefix of file controlled with -o flag.
-action force-undirected			Force the graph into an undirected graph. Writes to a different file.

-h or -help to display this message and quit.
`

func help() {
	fmt.Println(helpMessage)
	flag.PrintDefaults()
}

type Action func()

func main() {
	flag.Usage = help
	flag.Parse()

	if *flagHelp || *flagH {
		help()
		return
	}

	if *flagAction == "" {
		fmt.Println("Need to specify an action! See -help.")
		return
	}

	controller := map[string]Action{
		"details":          actionDetails,
		"components":       actionComponents,
		"remove":           actionRemove,
		"force-undirected": actionForceUndirected,
	}

	action, ok := controller[*flagAction]
	if !ok {
		fmt.Print("Invalid action! See -help for info. Available actions are: ")
		for txt, _ := range controller {
			fmt.Print(txt)
			fmt.Print(" ")
		}
		fmt.Println()
		return
	}

	action()
}

// Return whether the given graphs are directed on undirected.
func actionDetails() {
	files := flag.Args()
	fmt.Println("filename\t\tType\t\tConnected?\t#components\t#nodes\t#edges")
	fmt.Println("========\t\t====\t\t==========\t===========\t======\t======")
	for _, f := range files {
		graph, err := ReadGraph(f)
		if err != nil {
			fmt.Printf("%s \t\t Cannot read graph (%s).\n", f, err.Error())
			continue
		}

		fmt.Printf("%s\t\t", f)
		undir := graph.IsUndirected()
		if undir {
			fmt.Print("Undirected\t")
		} else {
			fmt.Print("Directed\t")
		}

		components := graph.ConnectedGraphs()
		if len(components) <= 1 {
			fmt.Print("connected\t")
		} else {
			fmt.Print("disconnected\t")
		}

		fmt.Printf("%d\t\t", len(components))
		fmt.Printf("%d\t", len(graph.Nodes))

		var noEdges int
		edges := graph.EdgeList()
		if undir {
			noEdges = len(edges) / 2
		} else {
			noEdges = len(edges)
		}

		fmt.Printf("%d\n", noEdges)
	}
}

// Splits graphs into their connected compoments, and writes those components as separate files.
func actionComponents() {
	files := flag.Args()
	for _, f := range files {
		graph, err := ReadGraph(f)
		if err != nil {
			fmt.Printf("%s: Cannot read graph. Skipping. (%s)\n", f, err.Error())
			continue
		}

		components := graph.ConnectedGraphs()

		for i, comp := range components {
			fname := *flagOutput + strconv.Itoa(i) + "_" + f
			wr, err := os.Create(fname)
			if err != nil {
				fmt.Printf("%s: Cannot write to %s, skipping connected component #%d. (%s)\n", f, fname, i, err.Error())
			}
			writer := util.NewWriter(wr)
			comp.EachEdge(func(from, to *Node) bool {
				err := writer.Write(from.Id, to.Id, nil)
				if err != nil {
					fmt.Printf("%s: Cannot write edge to %s, skipping remaining of connected component #%d. (%s)\n", f, fname, i, err.Error())
					return false
				}
				return true
			})
			err = writer.Flush()
			if err != nil {
				fmt.Printf("%s: (FLUSH) Output might be corrupted.")
			}
			err = wr.Close()
			if err != nil {
				fmt.Printf("%s: (CLOSE) Output might be corrupted.")
			}
		}
	}
}

// Get (perhaps directed) graphs and make them directed by forcing all edges to be bi-directional.
func actionForceUndirected() {
	files := flag.Args()
	for _, f := range files {
		graph, err := ReadGraph(f)
		if err != nil {
			fmt.Printf("%s: Cannot read graph. Skipping. (%s)\n", f, err.Error())
			continue
		}

		fname := *flagOutput + f
		wr, err := os.Create(fname)
		if err != nil {
			fmt.Printf("%s: Cannot write to %s, skipping file. (%s)\n", f, fname, err.Error())
			continue
		}
		writer := util.NewWriter(wr)
		edges := graph.EdgeList()
		for _, e := range edges {
			// it has from -> to if we got here.
			// adding to -> from.
			graph.AddDirectedEdge(e.To, e.From)
		}
		graph.EachEdge(func(from, to *Node) bool {
			err := writer.Write(from.Id, to.Id, nil)
			if err != nil {
				fmt.Printf("%s: Cannot write edge to %s, aborting this graph. (%s)\n", f, fname, err.Error())
				return false
			}
			return true
		})
		err = writer.Flush()
		if err != nil {
			fmt.Printf("%s: (FLUSH) Output might be corrupted.", f)
		}
		err = wr.Close()
		if err != nil {
			fmt.Printf("%s: (CLOSE) Output might be corrupted.", f)
		}
	}
}

// write a graph
func writeGraph(graph *Graph, fname string) error {
	wr, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("Cannot write to %s, skipping file. (%s)", fname, err.Error())
	}
	writer := util.NewWriter(wr)
	err = nil
	graph.EachEdge(func(from, to *Node) bool {
		err := writer.Write(from.Id, to.Id, nil)
		if err != nil {
			err = fmt.Errorf("Cannot write edge to %s, aborting this graph. (%s)", fname, err.Error())
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("(FLUSH) Output might be corrupted. (%s)", err.Error())
	}
	err = wr.Close()
	if err != nil {
		return fmt.Errorf("(CLOSE) Output might be corrupted. (%s)", err.Error())
	}
	return nil
}

// Remove random edges from a graph using a spanning tree to assure connectness.
func actionRemove() {
	files := flag.Args()
	for _, f := range files {
		graph, err := ReadGraph(f)
		if err != nil {
			fmt.Printf("%s: Cannot read graph. Skipping. (%s)\n", f, err.Error())
			continue
		}

		if !*flagForce && !graph.IsConnected() {
			fmt.Printf("%s: Graph not connected. Skipping. Use the --force to do it anyway.\n", f)
			continue
		}

		if !*flagForce && !graph.IsUndirected() {
			fmt.Printf("%s: Graph is directed. Skipping. Use the --force to do it anyway.\n", f)
			continue
		}

		edges := len(graph.EdgeList())
		remove := int(math.Floor(*flagN*float64(edges)/2 + 0.5))
		mst := graph.Mst()
		removed := graph.RemoveRandomEdges(remove, mst)

		// write out the processed graph
		fname := *flagOutput + f
		err = writeGraph(graph, fname)
		if err != nil {
			fmt.Printf("%s: %s\n", err.Error())
			continue
		}

		fname = *flagOutput + "removed_" + f
		wr, err := os.Create(fname)
		if err != nil {
			fmt.Printf("%s: Cannot write to %s, skipping file. (%s)", f, fname, err.Error())
			continue
		}

		writer := util.NewWriter(wr)
		for _, edge := range removed {
			err = writer.Write(edge.From, edge.To, nil)
			if err != nil {
				fmt.Printf("%s: Cannot write edge to %s. Skipping remaining of graph. (%s)", f, fname, err.Error())
				break
			}
			if *flagAll {
				err = writer.Write(edge.To, edge.From, nil)
				if err != nil {
					fmt.Printf("%s: Cannot write (reverse) edge to %s. Skipping remaining of graph. (%s)", f, fname, err.Error())
					break
				}
			}
		}

		err = writer.Flush()
		if err != nil {
			fmt.Printf("%s: (FLUSH) Output might be corrupted.", f)
		}
		err = wr.Close()
		if err != nil {
			fmt.Printf("%s: (CLOSE) Output might be corrupted.", f)
		}

	}
}
