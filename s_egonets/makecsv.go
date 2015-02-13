package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var mapping = make(map[int]int)
var id = 1
var pairs = make([]*pair, 0, 1000)

func mkmap(nr int) {
	_, ok := mapping[nr]
	if !ok {
		mapping[nr] = id
		id++
	}
}

type pair struct {
	a, b int
}

func (p *pair) write(output io.Writer) error {
	if p == nil {
		return fmt.Errorf("empty pair?? What the hack??")
	}
	_, err := fmt.Fprintf(output, "%d, %d\n", mapping[p.a], mapping[p.b])
	return err
}

func addPair(a, b string) error {
	aint, err := strconv.Atoi(a)
	if err != nil {
		return err
	}
	bint, err := strconv.Atoi(b)
	if err != nil {
		return err
	}
	mkmap(aint)
	mkmap(bint)
	pairs = append(pairs, &pair{aint, bint})
	return nil
}

func readFile(input string) error {
	file, err := os.Open(input)
	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		return err
	}
	inputStream := bufio.NewReader(file)
	for {
		first, err := inputStream.ReadString(':')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		second, err := inputStream.ReadString('\n')
		if err != nil {
			return err
		}
		first = strings.Trim(first, ",: ")
		ss := strings.Split(second, " ")
		for _, w := range ss {
			w = strings.TrimSpace(w)
			if w != "" {
				if err = addPair(first, w); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("At least one input file and one output file.")
		return
	}

	outputFile, err := os.Create(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println(err)
		return
	}
	output := bufio.NewWriter(outputFile)

	for _, name := range os.Args[1 : len(os.Args)-1] {
		err = readFile(name)
		if err != nil {
			fmt.Printf("Error at %s: %s\n", name, err)
		}
	}

	for _, p := range pairs {
		if err := p.write(output); err != nil {
			fmt.Println(err)
		}
	}

	output.Flush()
	outputFile.Close()

	fmt.Printf("Wrote %d pairs. Last node ID allocated: %d (starting at 1).\n", len(pairs), id-1)
}
