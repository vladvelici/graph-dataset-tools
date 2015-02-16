package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vladvelici/graph-dataset-tools/util"
)

// Define flags.
var (
	flagAction  = flag.String("action", "", "Possibile actions: details, components, remove")
	flagN       = flag.Int("n", 0, "Number of edges to remove from each graph.")
	flagOutput  = flag.String("o", "component_", "Output file prefix. It will be followed by the component number and an underscore.")
	flagVerbose = flag.Bool("verbose", false, "Whether to print lots of debug information on stdout.")
	flagHelp    = flag.Bool("help", false, "Show this help message")
	flagH       = flag.Bool("h", false, "Show this help message")
)

var helpMessage = `conncomp is a tool that deals with connected components tasks.

Possible uses (can optionally add -o and -verbose flags before filelist):

-action details filelist			Outputs details about the graphs in the files. Number of edges, number of nodes, directed/undirected, connected and no. of components.
-action components filelist			Splits the given graph in connected components files (prefix can be changed with -o flag).
-action remove -n P filelist		Removes at most P% random edges from the given graphs. Prefix of file controlled with -o flag.
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

// get a graph, split it into compoments, write back the graph..
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

// get a (perhaps directed) graph and make it directed.
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
			fmt.Printf("%s: (FLUSH) Output might be corrupted.")
		}
		err = wr.Close()
		if err != nil {
			fmt.Printf("%s: (CLOSE) Output might be corrupted.")
		}
	}
}

// get a graph, assume it's connected and undirected. Remove random edges from it.
func actionRemove() {
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
