/*
Processing inputs of form:
node1, node2, [whatever]

and outputting things of form:
node1_processed, node2_processed, [whatever]

Operations supported:

- create a mapping index
- apply a mapping index to a file
- revert a mapping index to a file
- overwrite the file at the end
*/
package main

import (
	"flag"
	"fmt"
	"github.com/vladvelici/graph-dataset-tools/util"
	"io"
	"os"
)

var helpMessage = ` Action explanations:

-action apply -index file.json [filenames]			apply given index (file.json) to [filenames]. Unseen nodes are added and written to the index.
-action revert -index file.json [filenames]			revert the given index and apply to [filenames]. Unseen nodes are untouched.
-action mkapply -index file.json [filenames]		apply autoincrement while creating an index and saving it to file.json.
-action index -index file.json [filenames]			create an index and save it to file.json, then quit.

Following is a flag usage:
`

func help() {
	fmt.Println(helpMessage)
	flag.PrintDefaults()
}

var (
	flagAction = flag.String("action", "", "The action to perform. Valid options: apply, revert, mkapply, index.")
	flagIndex  = flag.String("index", "-", "The index file. Writing or reading depends on action")
	flagPrefix = flag.String("prefix", "mappped_", "The prefix to append to output files, if not overwriting.")
	flagHelp   = flag.Bool("help", false, "Show this help message")
	flagH      = flag.Bool("h", false, "Show this help message")
)

func main() {
	flag.Usage = help
	flag.Parse()

	if *flagH || *flagHelp {
		help()
		return
	}

	var err error
	switch *flagAction {
	case "apply":
		files := flag.Args()
		if len(files) == 0 {
			fmt.Println("Need at least one input graph file.")
			return
		}
		err = apply(*flagIndex, files, *flagPrefix)
	case "revert":
		files := flag.Args()
		if len(files) == 0 {
			fmt.Println("Need at least one input graph file.")
			return
		}
		err = revert(*flagIndex, files, *flagPrefix)
	case "mkapply":
		files := flag.Args()
		if len(files) == 0 {
			fmt.Println("Need at least one input graph file.")
			return
		}
		err = mkapply(*flagIndex, files, *flagPrefix)
	case "index":
		files := flag.Args()
		if len(files) == 0 {
			fmt.Println("Need at least one input graph file.")
			return
		}
		err = index(*flagIndex, files, *flagPrefix)
	case "help":
		fallthrough
	case "h":
		fallthrough
	default:
		help()
		return
	}

	if err != nil {
		fmt.Println(err)
	}
}

// Apply index at indexPath to given files. Write updates to index. Output files are prefixed with prefix.
func apply(indexPath string, files []string, prefix string) error {
	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	index, err := ReadMapping(indexFile)
	if err != nil {
		return err
	}

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		output, err := os.Create(prefix + file)
		if err != nil {
			return err
		}

		inputCsv := util.NewReader(input)
		outputCsv := util.NewWriter(output)

		for {
			a, b, pass, err := inputCsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			tmp, ok := index.Node(a)
			if !ok {
				fmt.Println("%s: Found node %d which was not in index. Allocated to %d.", file, a, tmp)
			}
			a = tmp
			tmp, ok = index.Node(b)
			if !ok {
				fmt.Println("%s: Found node %d which was not in index. Allocated to %d.", file, b, tmp)
			}
			b = tmp

			err = outputCsv.Write(a, b, pass)
			if err != nil {
				return err
			}
		}

		ferr, err1, err2 := outputCsv.Flush(), input.Close(), output.Close()
		if ferr != nil {
			return ferr
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}

	err = indexFile.Close()
	if err != nil {
		return fmt.Errorf("Cannot close (reading) index file, so index file cannot be written. (%s)", err.Error())
	}
	indexFile, err = os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Cannot open index file for writing. (%s)", err.Error())
	}
	err = index.Write(indexFile)
	if err != nil {
		return fmt.Errorf("Cannot write index (index.Write). (%s)", err.Error())
	}
	indexFile.Close()
	return nil
}

// revert reverts the files using the index at indexPath. Unknown nodes are untouched. Output files are prefixed with prefix.
func revert(indexPath string, files []string, prefix string) error {
	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	index, err := ReadMapping(indexFile)
	if err != nil {
		return err
	}

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		output, err := os.Create(prefix + file)
		if err != nil {
			return err
		}

		inputCsv := util.NewReader(input)
		outputCsv := util.NewWriter(output)

		for {
			a, b, pass, err := inputCsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			tmp, ok := index.Allocation(a)
			if !ok {
				fmt.Println("%s: Found node %d which was not in index. Using %d in output.", file, a, tmp)
			}
			a = tmp
			tmp, ok = index.Allocation(b)
			if !ok {
				fmt.Println("%s: Found node %d which was not in index. Using %d in output.", file, b, tmp)
			}
			b = tmp

			err = outputCsv.Write(a, b, pass)
			if err != nil {
				return err
			}
		}

		ferr, err1, err2 := outputCsv.Flush(), input.Close(), output.Close()
		if ferr != nil {
			return ferr
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}

	err = indexFile.Close()
	if err != nil {
		return fmt.Errorf("Cannot close (reading) index file, so index file cannot be written. (%s)", err.Error())
	}
	indexFile, err = os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Cannot open index file for writing. (%s)", err.Error())
	}
	err = index.Write(indexFile)
	if err != nil {
		return fmt.Errorf("Cannot write index (index.Write). (%s)", err.Error())
	}
	indexFile.Close()
	return nil
}

// mkapply makes an autoincrement index over the given files. Writes the index to indexPath. Output files are prefixed with prefix.
func mkapply(indexPath string, files []string, prefix string) error {
	index := NewMapping()

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		output, err := os.Create(prefix + file)
		if err != nil {
			return err
		}

		inputCsv := util.NewReader(input)
		outputCsv := util.NewWriter(output)

		for {
			a, b, pass, err := inputCsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			a, _ = index.Node(a)
			b, _ = index.Node(b)

			err = outputCsv.Write(a, b, pass)
			if err != nil {
				return err
			}
		}

		ferr, err1, err2 := outputCsv.Flush(), input.Close(), output.Close()
		if ferr != nil {
			return ferr
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}

	indexFile, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Cannot open index file for writing. (%s)", err.Error())
	}
	err = index.Write(indexFile)
	if err != nil {
		return fmt.Errorf("Cannot write index (index.Write). (%s)", err.Error())
	}
	indexFile.Close()
	return nil
}

// index makes an autoincrement index over the given files. Writes the index to indexPath. Does not output any files.
func index(indexPath string, files []string, prefix string) error {
	index := NewMapping()

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			return err
		}

		inputCsv := util.NewReader(input)

		for {
			a, b, _, err := inputCsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			a, _ = index.Node(a)
			b, _ = index.Node(b)
		}

		if err = input.Close(); err != nil {
			return err
		}
	}

	indexFile, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Cannot open index file for writing. (%s)", err.Error())
	}
	err = index.Write(indexFile)
	if err != nil {
		return fmt.Errorf("Cannot write index (index.Write). (%s)", err.Error())
	}
	indexFile.Close()
	return nil
}
